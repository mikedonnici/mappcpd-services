package graphql

import (
	"os"

	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/graphql-go/graphql"
)

// Query resolves member queries, is a 'viewer' field for the member (user) identified by the token
var Query = &graphql.Field{
	Description: "Member queries require a valid JSON Web Encoded for auth and data in child nodes will always " +
		"belong to the member identified by the token.",
	Type: memberType,
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

// memberType defines fields for a Member.
var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "member", // This is the object Type name
	Description: "member query object that provides access to data for the member identified by the token.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's unique id number",
		},
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "A fresh token",
		},
		"active": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Boolean flag indicating if the member is currently active in the system",
		},
		"title": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's membership title",
		},
		"firstName": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's first name",
		},
		"middleNames": &graphql.Field{
			Type:        graphql.String,
			Description: "One or more middle names",
		},
		"lastName": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's surname / family name",
		},
		"postNominal": &graphql.Field{
			Type:        graphql.String,
			Description: "Option string of preferred post nominals, eg 'Ph.D', 'OAM' etc",
		},
		"dateOfBirth": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's date of birth, as a string value",
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's primary email address",
		},
		"mobile": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's mobile phone number",
		},
		"locations": &graphql.Field{
			Type:        graphql.NewList(locationType),
			Description: "One or more contact locations",
		},
		"qualifications": &graphql.Field{
			Type:        graphql.NewList(qualificationType),
			Description: "The member's qualifications",
		},
		"positions": &graphql.Field{
			Type:        graphql.NewList(positionType),
			Description: "The member's positions or appointments to committees, councils etc",
		},

		// child nodes / sub queries
		"activity":    activityQuery,
		"activities":  activitiesQuery,
		"evaluation":  currentEvaluationQuery,
		"evaluations": evaluationsQuery,
	},
})

// locationType defines fields for a contact location
var locationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "location",
	Description: "A contact location belonging to a member",
	Fields: graphql.Fields{
		"order": &graphql.Field{
			Type: graphql.Int,
		},
		"type": &graphql.Field{
			Type: graphql.String,
		},
		"address": &graphql.Field{
			Type: graphql.String,
		},
		"city": &graphql.Field{
			Type: graphql.String,
		},
		"state": &graphql.Field{
			Type: graphql.String,
		},
		"postcode": &graphql.Field{
			Type: graphql.String,
		},
		"country": &graphql.Field{
			Type: graphql.String,
		},
		"phone": &graphql.Field{
			Type: graphql.String,
		},
		"fax": &graphql.Field{
			Type: graphql.String,
		},
		"email": &graphql.Field{
			Type: graphql.String,
		},
		"url": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// contactType defines fields for a member's contact information, containing one or more locations
var contactType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "contact",
	Description: "Member contact information which may include one or more contact locations.",
	Fields: graphql.Fields{
		"emailPrimary": &graphql.Field{
			Type: graphql.String,
		},
		"emailSecondary": &graphql.Field{
			Type: graphql.String,
		},
		"mobile": &graphql.Field{
			Type: graphql.String,
		},
		"locations": &graphql.Field{
			Type: graphql.NewList(locationType),
		},
	},
})

// qualificationType defines fields for a qualification obtained by a Member
var qualificationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "qualification",
	Description: "An academic qualification obtained by the Member",
	Fields: graphql.Fields{
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"year": &graphql.Field{
			Type: graphql.String,
		},
	},
})

// positionType defines fields for a position held by a Member
var positionType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "position",
	Description: "A position or affiliation with a council, committee or group",
	Fields: graphql.Fields{
		"orgCode": &graphql.Field{
			Type: graphql.String,
		},
		"orgName": &graphql.Field{
			Type: graphql.String,
		},
		"code": &graphql.Field{
			Type: graphql.String,
		},
		"name": &graphql.Field{
			Type: graphql.String,
		},
		"description": &graphql.Field{
			Type: graphql.String,
		},
		"startDate": &graphql.Field{
			Type: graphql.String,
		},
		"endDate": &graphql.Field{
			Type: graphql.String,
		},
	},
})
