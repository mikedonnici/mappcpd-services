package graphql

import (
	"os"

	"github.com/cardiacsociety/web-services/internal/attachments"
	"github.com/cardiacsociety/web-services/internal/date"
	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/graphql-go/graphql"
)

// activityQuery resolves a query for a single member activity
var activityQuery = &graphql.Field{
	Description: "Fetches a single member activity by id.",
	Type:        activityType,
	Args: graphql.FieldConfigArgument{
		"activityId": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the member activityQuery",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Always extract the member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Decode(token.(string), os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		activityID, ok := p.Args["activityId"].(int)
		if ok {
			return mapActivityData(memberID, int(activityID))
		}

		return nil, nil
	},
}

// activitiesQuery resolves a query for member activities
var activitiesQuery = &graphql.Field{
	Description: "Fetches a list of member activities",
	Type:        graphql.NewList(activityType),
	Args: graphql.FieldConfigArgument{
		"last": &graphql.ArgumentConfig{
			Type:        graphql.Int,
			Description: "Fetch only the last (most recent) n records.",
		},
		"from": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Fetch activities from this date - format 'YYYY-MM-DD'",
		},
		"to": &graphql.ArgumentConfig{
			Type:        graphql.String,
			Description: "Fetch activities up to and including this date - format 'YYYY-MM-DD'",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Decode(token.(string), os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		// Filter arguments
		f := make(map[string]interface{})
		last, ok := p.Args["last"].(int)
		if ok {
			f["last"] = last
		}
		from, ok := p.Args["from"].(string)
		if ok {
			f["from"], err = date.StringToTime(from)
			if err != nil {
				return nil, err
			}
		}
		to, ok := p.Args["to"].(string)
		if ok {
			f["to"], err = date.StringToTime(to)
			if err != nil {
				return nil, err
			}
		}

		return mapActivitiesData(memberID, f)
	},
}

// activityType defines fields for a Member activity
var activityType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "activityData",
	Description: "An instance of an activity recorded by a Member - ie an entry in the CPD diary.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "ID of the Member activity record.",
		},
		"date": &graphql.Field{
			Type:        graphql.String,
			Description: "The date the activity was undertaken, format 'YYYY-MM-DD'.",
		},
		"dateTime": &graphql.Field{
			Type:        graphql.DateTime,
			Description: "The date the activity was undertaken in RFC3339 format time set to 00:00:00 UTC.",
		},
		"quantity": &graphql.Field{
			Type:        graphql.Float,
			Description: "Quantity, generally number of hours, for an activity",
		},
		"creditPerUnit": &graphql.Field{
			Type:        graphql.Float,
			Description: "Credit per unit of the activity, ie multiply this x quantity = total credit",
		},
		"credit": &graphql.Field{
			Type:        graphql.Float,
			Description: "The total credit for the activity, ie quanity x creditPerUnit.",
		},
		"activity": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the activity.",
		},
		"activityId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the activity.",
		},
		"categoryId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the category to which the activity belongs.",
		},
		"category": &graphql.Field{
			Type:        graphql.String,
			Description: "The name of the category to which the activity belongs",
		},
		"type": &graphql.Field{
			Type: graphql.String,
			Description: "Type represents a specific form, or example of an activity, Where 'category' is the broadest" +
				"descriptive attribute, 'type' is the most specific.",
		},
		"typeId": &graphql.Field{
			Type:        graphql.Int,
			Description: "Activity type ID.",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "Descriptive details about the activity, supplied by the Member.",
		},
		"evidence": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "A flag that indicates if the user has supporting evidence for the activity",
		},

		"attachments": activityAttachmentsQuery,
	},
})

// activityAttachmentsQuery resolves a query for member activity attachments
var activityAttachmentsQuery = &graphql.Field{
	Description: "Fetches a list of attachments for a member activity",
	Type:        graphql.NewList(activityAttachmentType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// member activity id from the parent node
		maID := p.Source.(activityData).ID

		return attachments.MemberActivityAttachments(DS, maID)
	},
}

// activityAttachmentType defines fields for a Member activity attachment
var activityAttachmentType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "activityAttachmentData",
	Description: "An attachment associated with the member activity",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the member activity attachment record",
		},
		// todo this should be a signed url
		"url": &graphql.Field{
			Type:        graphql.String,
			Description: "The url for accessing the file",
		},
	},
})
