package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
	"github.com/mappcpd/web-services/internal/platform/jwt"
)

// MemberUserInput is top-level fields for mutations performed by member users.
var MemberUserInput = &graphql.Field{
	Description: "Top-level input field for member data.",
	Type:        memberInputType,
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "Valid JSON web token",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		token, ok := p.Args["token"].(string)
		if ok {
			// Validate the token, and extract the member id
			at, err := jwt.Check(token)
			if err != nil {
				return nil, err
			}
			id := at.Claims.ID

			// For the memberInput type we only want the member id, and, to be honest, don't really even need that
			//return data.GetMember(id)
			return map[string]interface{}{"id": id}, nil
		}
		return nil, nil
	},
}

// memberInputType is entry point for member mutations, requiring a valid token.
var memberInputType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberInput",
	Description: "Top-level input for member fields",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "Unique id of the member performing the operation, extracted from the token.",
		},

		"setActivity": memberActivityInput,
	},
})

// memberActivity will either add a new member memberActivity, or edit an existing one, when the member memberActivity id is provided.
var memberActivityInput = &graphql.Field{
	Description: "Add or update a member activity. If `activityId` is present in the argument object, and the record " +
		"belongs to the authenticated member, then it will be updated. If activityId is not present, or does not belong " +
		"to the authenticated user, a new member activity record will be created.",
	Type: memberActivityType, // this is what will return from this operation, ie the type we are mutating
	Args: graphql.FieldConfigArgument{
		"obj": &graphql.ArgumentConfig{
			Type:        memberActivityInputType, // this is the type required as the arg
			Description: "An object containing the necessary fields to add or update a member activity",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Always extract the member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		maObj, ok := p.Args["obj"].(map[string]interface{})
		if ok {
			ma := data.MemberActivity{}
			err := ma.Unpack(maObj)
			if err != nil {
				return nil, err
			}

			if ma.ID > 0 {
				return data.UpdateMemberActivity(memberID, ma)
			}

			return data.AddMemberActivity(memberID, ma)
		}

		return nil, nil
	},
}
