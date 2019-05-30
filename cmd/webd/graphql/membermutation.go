package graphql

import (
	"os"

	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/graphql-go/graphql"
)

// Mutation handles mutations for member data
var Mutation = &graphql.Field{
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

			at, err := jwt.Decode(token, os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
			if err != nil {
				return nil, err
			}
			id := at.Claims.ID

			m, err := mapMemberData(id)
			if err != nil {
				return nil, err
			}

			m.Token, err = freshToken(token)
			if err != nil {
				return m, err
			}

			return m, nil
		}

		return nil, nil
	},
}

// memberInputType defines fields for mutating member data
var memberInputType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberInput",
	Description: "Top-level input for member fields",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "Unique id of the member performing the operation, extracted from the token.",
		},
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "A fresh token",
		},
		"saveActivity":   activitySave,
		"deleteActivity": activityDelete,
	},
})
