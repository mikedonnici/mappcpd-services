package graphql

import (
	"fmt"
	"os"

	"github.com/cardiacsociety/web-services/internal/cpd"
	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"
)

// activitySave handles mutation (add / update) of a member activity
var activitySave = &graphql.Field{
	Description: "Add or update a member activity. If `activityId` is present in the argument object, and the record " +
		"belongs to the member identified by the token, then it will be updated. If `activityId` is not present, or does not belong " +
		"to the authenticated user, a new member activity record will be created.",
	Type: activityType, // this type will be returned this operation
	Args: graphql.FieldConfigArgument{
		"obj": &graphql.ArgumentConfig{
			Type:        activityInputType, // this is the type required as the arg
			Description: "An object containing the necessary fields to add or update a member activity",
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

		maObj, ok := p.Args["obj"].(map[string]interface{})
		if ok {

			ma := activityInputData{}
			err := ma.unpack(maObj)
			if err != nil {
				return nil, err
			}

			// set activity id from the activity type id
			ma.ActivityID, err = activityIDByTypeID(ma.TypeID)
			if err != nil {
				msg := fmt.Sprintf("Error fetching activity with activity type id = %v", ma.TypeID)
				return nil, errors.Wrap(err, msg)
			}

			// update record
			if ma.ID > 0 {
				return updateActivity(memberID, ma)
			}

			// add record, ensure not a duplicate
			dupID, err := activityDuplicateID(memberID, ma)
			if dupID > 0 {
				msg := fmt.Sprintf("The activity is an exact duplicate of id %v. To copy an activity at least one "+
					"field must be changed, eg date.", dupID)
				return nil, errors.New(msg)
			}
			if err != nil {
				msg := fmt.Sprintf("Error checking for duplicate activity - %s", err.Error())
				return nil, errors.New(msg)
			}

			return addActivity(memberID, ma)
		}

		return nil, nil
	},
}

// activityInputType defines fields for mutating a member activity
var activityInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "activitySaveInput",
	Description: "An input object type used as an argument for adding / updating a member activity",
	Fields: graphql.InputObjectConfigFieldMap{
		// optional member activity id - if supplied then it is an update
		"id": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "Optional id of the member activity record, if present will update existing, otherwise will add new.",
		},

		// typeId specifies the type of activity
		"typeId": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the activity type",
		},

		// date on which the activity was undertaken
		"date": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The date on which the activity was undertaken",
		},

		// quantity, generally in hours
		"quantity": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Float},
			Description: "The number of units of the activity being recorded, generally the number of hours",
		},

		// description supplied by the user
		"description": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The specifics of the activityQuery described by the member",
		},

		"evidence": &graphql.InputObjectFieldConfig{
			Type:         graphql.Boolean,
			Description:  "A flag to indicate that the user has evidence to support this activity record",
			DefaultValue: false,
		},
	},
})

// activityDelete handles mutation (add / update) of a member activity
var activityDelete = &graphql.Field{
	Description: "Delete an activity that belongs to the member identified by the token",
	Type:        graphql.String, // this type will be returned this operation
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type:        graphql.Int, // this is the type required as the arg
			Description: "The id of the record to be deleted",
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

		activityID, ok := p.Args["id"].(int)
		if ok {
			return "CPD deleted", cpd.Delete(DS, memberID, activityID)
		}
		return nil, nil
	},
}
