package schema

import (
	"fmt"
	"time"

	"github.com/graphql-go/graphql"
	"github.com/pkg/errors"

	"github.com/mappcpd/web-services/internal/member"
	"github.com/mappcpd/web-services/internal/member/activity"
	"github.com/mappcpd/web-services/internal/platform/jwt"
	"github.com/mappcpd/web-services/internal/utility"
	"github.com/mappcpd/web-services/internal/attachments"
)

// localMember is a local representation of localMember.Member
type localMember struct {
	ID             int                     `json:"id"`
	Token          string                  `json:"token"`
	Active         bool                    `json:"active"`
	Title          string                  `json:"title"`
	FirstName      string                  `json:"firstName"`
	MiddleNames    string                  `json:"middleNames"`
	LastName       string                  `json:"lastName"`
	PostNominal    string                  `json:"postNominal"`
	DateOfBirth    string                  `json:"dateOfBirth"`
	Email          string                  `json:"email"`
	Mobile         string                  `json:"mobile"`
	Locations      []member.MemberLocation `json:"locations"`
	Qualifications []member.Qualification  `json:"qualifications"`
	Positions      []member.Position       `json:"positions"`
}

// memberActivity is a leaner representation of members.memberActivity
type memberActivity struct {
	// ID is the unique id of the localMember activity record
	ID int `json:"id"`

	// Date on which the activity was undertaken, and an equivalent Time value
	Date     string    `json:"date"`
	DateTime time.Time `json:"dateTime"`

	// Quantity is generally the number of hours, but may be other units
	Quantity float64 `json:"quantity"`

	// Credit is generally the number of hours multiplied by credit-per-unit for a particular activity
	CreditPerUnit float64 `json:"creditPerUnit"`

	// Credit is generally the number of hours multiplied by credit-per-unit for a particular activity
	Credit float64 `json:"credit"`

	// Description is the user-input that further describes the activity itself
	Description string `json:"description"`

	// The following are descriptive of the type of activity undertaken. yes, this is a disaster - will fix later
	// The data relationship is: Category -> Activity -> Type, which is straightforward enough. However the Type
	// was added in much later for a compliance reason and creates some confusion as, in many parts of this code,
	// the word 'type' is used to describe activity (type) from the ce_activity table.

	// ActivityID is the id of the activity (type), ie, a record from the ce_activity table. Until the new 'type'
	// came along this was often described as activityType in var names etc. This was to avoid confusion with an
	// actual localMember activity, but has now caused more confusion.
	ActivityID int `json:"activityId"`
	// Activity is the string name of the ce_activity record
	Activity string `json:"activity"`

	// CategoryID is the parent category, ie a record from ce_activity_category
	CategoryID int `json:"categoryId"`
	// Category is the name
	Category string `json:"category"`

	// TypeID now refers to the activity sub-type, ie a record from the ce_activity_type table
	TypeID int `json:"typeId"`
	// type is the string name of the activity sub-type
	Type string `json:"type"`

	// Attachments
	//Attachments []Attachment

	// todo: remove this UploadURL is a signed URL that allows for uploading file attachments
	UploadURL string `json:"uploadUrl"`
}

// memberActivityInput represents an object for mutating a localMember activity
type memberActivityInput struct {
	// ID - if present triggers and update, else record will be added
	ID int `json:"id"`

	// Date on which the activity was undertaken as a string "YYYY-MM-DD"
	Date string `json:"date"`

	// Quantity of the units relevant to the activity, generally hours
	Quantity float64 `json:"quantity"`

	// Description is the user-input that further describes the activity itself
	Description string `json:"description"`

	// ActivityID is the id of the activity (type) - can look up from type
	ActivityID int `json:"activityId"`

	// TypeID now refers to the activity sub-type, ie a record from the ce_activity_type table
	TypeID int `json:"typeId"`
}

// memberActivityAttachment represents an file associated with a localMember activity
type memberActivityAttachment struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

// memberEvaluation representations the localMember evaluation data
type memberEvaluation struct {
	//ID          int       `json:"id"`
	ReportName     string  `json:"name"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	CreditRequired float64 `json:"creditRequired"`
	CreditObtained float64 `json:"creditObtained"`
	Closed         bool    `json:"closed"`
}

// memberQueryField resolves localMember queries, is a 'viewer' field for the localMember (user) identified by the token
var memberQueryField = &graphql.Field{
	Description: "Member queries require a valid JSON Web Token for auth and data in child nodes will always " +
		"belong to the localMember identified by the token.",
	Type: memberQueryObject,
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "Valid JSON web token",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		token, ok := p.Args["token"].(string)
		if ok {

			// Validate the token, and extract the localMember id
			at, err := jwt.Check(token)
			if err != nil {
				return nil, err
			}
			id := at.Claims.ID

			// At this point we have a valid token from which we've extracted an id.
			// As a final step we can verify that the id is a valid user in the system,
			// for example, that it is active. Although this is a bit redundant for each request?

			// create the localMember value
			m, err := memberData(id)
			if err != nil {
				return nil, err
			}

			// set the fresh token
			m.Token, err = member.FreshToken(token)
			if err != nil {
				return m, err
			}

			return m, nil
		}

		return nil, nil
	},
}

// memberQueryObject defines fields for a localMember.
var memberQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "localMember",
	Description: "localMember query object that provides access to data for the localMember identified by the token.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's unique id number",
		},
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "A fresh token",
		},
		"active": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Boolean flag indicating if the localMember is currently active in the system",
		},
		"title": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's membership title",
		},
		"firstName": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's first name",
		},
		"middleNames": &graphql.Field{
			Type:        graphql.String,
			Description: "One or more middle names",
		},
		"lastName": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's surname / family name",
		},
		"postNominal": &graphql.Field{
			Type:        graphql.String,
			Description: "Option string of preferred post nominals, eg 'Ph.D', 'OAM' etc",
		},
		"dateOfBirth": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's date of birth, as a string value",
		},
		"email": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's primary email address",
		},
		"mobile": &graphql.Field{
			Type:        graphql.String,
			Description: "The localMember's mobile phone number",
		},
		"locations": &graphql.Field{
			Type:        graphql.NewList(locationQueryObject),
			Description: "One or more contact locations",
		},
		"qualifications": &graphql.Field{
			Type:        graphql.NewList(qualificationQueryObject),
			Description: "The localMember's qualifications",
		},
		"positions": &graphql.Field{
			Type:        graphql.NewList(positionQueryObject),
			Description: "The localMember's positions or appointments to committees, councils etc",
		},

		// child nodes / sub queries
		"activity":    memberActivityQueryField,
		"activities":  memberActivitiesQueryField,
		"evaluation":  memberCurrentEvaluationQueryField,
		"evaluations": memberEvaluationsQueryField,
	},
})

// locationQueryObject defines fields for a contact location
var locationQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "location",
	Description: "A contact location belonging to a localMember",
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

// contactQueryObject defines fields for a localMember's contact information, containing one or more locations
var contactQueryObject = graphql.NewObject(graphql.ObjectConfig{
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
			Type: graphql.NewList(locationQueryObject),
		},
	},
})

// qualificationQueryObject defines fields for a qualification obtained by a localMember
var qualificationQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "qualification",
	Description: "An academic qualification obtained by the localMember",
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

// positionQueryObject defines fields for a position held by a localMember
var positionQueryObject = graphql.NewObject(graphql.ObjectConfig{
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

// memberActivityQueryObject defines fields for a localMember activity
var memberActivityQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberActivity",
	Description: "An instance of an activity recorded by a localMember - ie an entry in the CPD diary.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "ID of the localMember activity record.",
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
			Description: "Descriptive details about the activity, supplied by the localMember.",
		},

		"attachments": memberActivityAttachmentsQueryField,
	},
})

// memberActivityAttachmentQueryObject defines fields for a localMember activity attachment
var memberActivityAttachmentQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberActivityAttachment",
	Description: "An attachment associated with the localMember activity",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the localMember activity attachment record",
		},
		// todo this should be a signed url
		"url": &graphql.Field{
			Type:        graphql.String,
			Description: "The url for accessing the file",
		},
	},
})

// memberEvaluationQueryObject defines fields for a localMember evaluation
var memberEvaluationQueryObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberEvaluation",
	Description: "An evaluation of activity credited and required, for a given period of time - eg a calendar year.",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "The name fo the evaluation period - eg Annual CPD Requirement",
		},
		"startDate": &graphql.Field{
			Type:        graphql.String,
			Description: "The start date of the evaluation period.",
		},
		"endDate": &graphql.Field{
			Type:        graphql.String,
			Description: "The end date of the evaluation period.",
		},
		"creditRequired": &graphql.Field{
			Type:        graphql.Float,
			Description: "Value or credit required to satisfy the evaluation period requirements.",
		},
		"creditObtained": &graphql.Field{
			Type:        graphql.Float,
			Description: "Actual activity credit gained for the period.",
		},
		"closed": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "Indicated if the evaluation period is closed.",
		},
	},
})

// memberActivityQueryField resolves a query for a single localMember activity
var memberActivityQueryField = &graphql.Field{
	Description: "Fetches a single localMember activity by id.",
	Type:        memberActivityQueryObject,
	Args: graphql.FieldConfigArgument{
		"activityId": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the localMember memberActivityQueryField",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Always extract the localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		activityID, ok := p.Args["activityId"].(int)
		if ok {
			return memberActivityData(memberID, int(activityID))
		}

		return nil, nil
	},
}

// memberActivityAttachmentsQueryField resolves a query for localMember activity attachments
var memberActivityAttachmentsQueryField = &graphql.Field{
	Description: "Fetches a list of attachments for a localMember activity",
	Type:        graphql.NewList(memberActivityAttachmentQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract localMember id from the token, available thus:
		//token := p.Info.VariableValues["token"]
		//at, err := jwt.Check(token.(string))
		//if err != nil {
		//	return nil, err
		//}
		//memberID := at.Claims.ID

		// Get the localMember activity id from the parent
		maID := p.Source.(memberActivity).ID
		//types, err := activityTypesData(id)
		//if err != nil {
		//	return nil, nil
		//}

		return memberActivityAttachmentsData(maID)
	},
}

// memberActivitiesQueryField resolves a query for localMember activities
var memberActivitiesQueryField = &graphql.Field{
	Description: "Fetches a list of localMember activities",
	Type:        graphql.NewList(memberActivityQueryObject),
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

		// Extract localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
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
			f["from"], err = utility.DateStringToTime(from)
			if err != nil {
				return nil, err
			}
		}
		to, ok := p.Args["to"].(string)
		if ok {
			f["to"], err = utility.DateStringToTime(to)
			if err != nil {
				return nil, err
			}
		}

		return memberActivitiesData(memberID, f)
	},
}

// memberCurrentEvaluationQueryField resolves queries for the current evaluation period
var memberCurrentEvaluationQueryField = &graphql.Field{
	Description: "Fetches activity data for the current evaluation period",
	Type:        memberEvaluationQueryObject,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return memberCurrentEvaluationData(memberID)
	},
}

// memberEvaluationsQueryField resolves queries for multiple localMember evaluation periods
var memberEvaluationsQueryField = &graphql.Field{
	Description: "Fetches a history of localMember activity evaluation periods",
	Type:        graphql.NewList(memberEvaluationQueryObject),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return memberEvaluationsData(memberID)
	},
}

// memberMutationField handles mutations for localMember data
var memberMutationField = &graphql.Field{
	Description: "Top-level input field for localMember data.",
	Type:        memberMutationObject,
	Args: graphql.FieldConfigArgument{
		"token": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "Valid JSON web token",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		token, ok := p.Args["token"].(string)
		if ok {
			// Validate the token, and extract the localMember id
			at, err := jwt.Check(token)
			if err != nil {
				return nil, err
			}
			id := at.Claims.ID

			// create the localMember value
			m, err := memberData(id)
			if err != nil {
				return nil, err
			}

			// set the fresh token
			m.Token, err = member.FreshToken(token)
			if err != nil {
				return m, err
			}

			return m, nil
		}

		return nil, nil
	},
}

// memberMutationObject defines fields for mutating localMember data
var memberMutationObject = graphql.NewObject(graphql.ObjectConfig{
	Name:        "memberInput",
	Description: "Top-level input for localMember fields",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "Unique id of the localMember performing the operation, extracted from the token.",
		},
		"token": &graphql.Field{
			Type:        graphql.String,
			Description: "A fresh token",
		},

		"saveActivity":   memberActivitySaveField,
		"deleteActivity": memberActivityDeleteField,
	},
})

// memberActivitySaveField handles mutation (add / update) of a localMember activity
var memberActivitySaveField = &graphql.Field{
	Description: "Add or update a localMember activity. If `activityId` is present in the argument object, and the record " +
		"belongs to the localMember identified by the token, then it will be updated. If `activityId` is not present, or does not belong " +
		"to the authenticated user, a new localMember activity record will be created.",
	Type: memberActivityQueryObject, // this type will be returned this operation
	Args: graphql.FieldConfigArgument{
		"obj": &graphql.ArgumentConfig{
			Type:        memberActivitySaveObject, // this is the type required as the arg
			Description: "An object containing the necessary fields to add or update a localMember activity",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Always extract the localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		maObj, ok := p.Args["obj"].(map[string]interface{})
		if ok {

			ma := memberActivityInput{}
			err := ma.unpack(maObj)
			if err != nil {
				return nil, err
			}

			// set activity id from the activity type id
			ma.ActivityID, err = activityIDByActivityTypeID(ma.TypeID)
			if err != nil {
				msg := fmt.Sprintf("Error fetching activity with activity type id = %v", ma.TypeID)
				return nil, errors.Wrap(err, msg)
			}

			// update record
			if ma.ID > 0 {
				return updateMemberActivity(memberID, ma)
			}

			// add record, ensure not a duplicate
			dupID := memberActivityDuplicate(memberID, ma)
			if dupID > 0 {
				msg := fmt.Sprintf("Cannot add new activity as it is a duplicate of localMember activity id %v - "+
					"include { id: %v } in the object argument to update instead", dupID, dupID)
				return nil, errors.New(msg)
			}

			return addMemberActivity(memberID, ma)
		}

		return nil, nil
	},
}

// memberActivitySaveObject defines fields for mutating a localMember activity
var memberActivitySaveObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "memberActivityInput",
	Description: "An input object type used as an argument for adding / updating a localMember activity",
	Fields: graphql.InputObjectConfigFieldMap{
		// optional localMember activity id - if supplied then it is an update
		"id": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "Optional id of the localMember activity record, if present will update existing, otherwise will add new.",
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
			Description: "The specifics of the memberActivityQueryField described by the localMember",
		},
	},
})

// memberActivityDeleteField handles mutation (add / update) of a localMember activity
var memberActivityDeleteField = &graphql.Field{
	Description: "Delete an activity that belongs to the localMember identified by the token",
	Type:        graphql.String, // this type will be returned this operation
	Args: graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type:        graphql.Int, // this is the type required as the arg
			Description: "The id of the record to be deleted",
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Always extract the localMember id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		activityID, ok := p.Args["id"].(int)
		if ok {
			return "Record deleted", deleteMemberActivity(memberID, activityID)
		}
		return nil, nil
	},
}

// memberActivityDeleteObject defines fields for deleting a localMember activity
var memberActivityDeleteObject = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "memberActivityInput",
	Description: "An input object type used as an argument for adding / updating a localMember activity",
	Fields: graphql.InputObjectConfigFieldMap{
		// optional localMember activity id - if supplied then it is an update
		"id": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the localMember activity record to be deleted.",
		},
	},
})

// unpack an object into a value of type MemberActivity
func (ma *memberActivity) unpack(obj map[string]interface{}) error {
	if val, ok := obj["id"].(int); ok {
		ma.ID = val
	}
	if val, ok := obj["date"].(string); ok {
		ma.Date = val
		d, err := utility.DateStringToTime(val)
		if err != nil {
			return err
		}
		ma.DateTime = d
	}
	if val, ok := obj["quantity"].(float64); ok {
		ma.Quantity = val
	}
	if val, ok := obj["creditPerUnit"].(float64); ok {
		ma.CreditPerUnit = val
	}
	if val, ok := obj["credit"].(float64); ok {
		ma.Credit = val
	}
	if val, ok := obj["categoryId"].(int); ok {
		ma.CategoryID = int(val)
	}
	if val, ok := obj["activityId"].(int); ok {
		ma.ActivityID = int(val)
	}
	if val, ok := obj["typeId"].(int); ok {
		ma.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		ma.Description = val
	}

	return nil
}

// unpack an object into a value of type MemberActivityInput
func (mai *memberActivityInput) unpack(obj map[string]interface{}) error {
	if val, ok := obj["id"].(int); ok {
		mai.ID = val
	}
	if val, ok := obj["date"].(string); ok {
		mai.Date = val
	}
	if val, ok := obj["quantity"].(float64); ok {
		mai.Quantity = val
	}
	if val, ok := obj["typeId"].(int); ok {
		mai.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		mai.Description = val
	}

	return nil
}

// memberData fetches the basic localMember record
func memberData(id int) (localMember, error) {
	var m localMember
	mp, err := memberProfileData(id)
	if err != nil {
		return m, err
	}
	m.ID = mp.ID
	m.Active = mp.Active
	m.Title = mp.Title
	m.FirstName = mp.FirstName
	m.MiddleNames = mp.MiddleNames
	m.LastName = mp.LastName
	m.DateOfBirth = mp.DateOfBirth
	m.Email = mp.Contact.EmailPrimary
	m.Mobile = mp.Contact.Mobile
	m.PostNominal = mp.PostNominal
	m.Locations = mp.Contact.Locations
	m.Qualifications = mp.Qualifications
	m.Positions = mp.Positions

	return m, nil
}

// memberProfileData fetches a single localMember record by id
func memberProfileData(memberID int) (member.Member, error) {
	// MemberByID returns a pointer to a members.localMember so dereference in return
	m, err := member.MemberByID(memberID)
	return *m, err
}

// memberActivitiesData fetches activities for a localMember.
func memberActivitiesData(memberID int, filter map[string]interface{}) ([]memberActivity, error) {
	var xa []memberActivity

	// This returns a nested struct which is simplified below.
	xma, err := activity.MemberActivitiesByMemberID(memberID)

	// Set up date filters
	from, okFrom := filter["from"].(time.Time)
	to, okTo := filter["to"].(time.Time)
	if okFrom && okTo {
		if from.After(to) {
			return xa, errors.New("from date cannot be after to date")
		}
	}

	for _, v := range xma {

		// Apply date filters, skip to next iteration if the data is outside the range
		if okFrom {
			if v.DateISO.Before(from) {
				continue
			}
		}
		if okTo {
			if v.DateISO.After(to) {
				continue
			}
		}

		// Passed through date filters, add the record to our simplified struct
		a := memberActivity{
			ID:            v.ID,
			Date:          v.Date,
			DateTime:      v.DateISO,
			Quantity:      v.CreditData.Quantity,
			CreditPerUnit: v.CreditData.UnitCredit,
			Credit:        v.Credit,
			CategoryID:    v.Category.ID,
			Category:      v.Category.Name,
			ActivityID:    v.Activity.ID,
			Activity:      v.Activity.Name,
			TypeID:        int(v.Type.ID.Int64), // null-able field
			Type:          v.Type.Name,
			Description:   v.Description,
		}
		xa = append(xa, a)
	}

	// Although less efficient, apply 'last' n filter last - otherwise it cannot be used in conjunction with
	// the date filters.
	last, ok := filter["last"].(int)
	if ok {
		// Activities are returned in reverse order so returning the 'last' n items, ie the most *recent*, means
		// slicing from the index 0. If n is greater than the total, just return the total.
		if last < len(xma) {
			xa = xa[:last]
		}
	}

	return xa, err
}

// memberActivityData fetches a single localMember activity by ID after verifying ownership by memberID
func memberActivityData(memberID, memberActivityID int) (memberActivity, error) {

	var a memberActivity

	// This returns a nested struct which we can simplify
	ma, err := activity.MemberActivityByID(memberActivityID)
	if err != nil {
		return a, err
	}

	// Verify owner match
	if ma.MemberID != memberID {
		msg := fmt.Sprintf("Member activity (id %v) does not belong to localMember (id %v)", memberActivityID, memberID)
		return a, errors.New(msg)
	}

	a.ID = ma.ID
	a.Date = ma.Date
	a.DateTime = ma.DateISO
	a.Quantity = ma.CreditData.Quantity
	a.CreditPerUnit = ma.CreditData.UnitCredit
	a.Credit = ma.Credit
	a.CategoryID = ma.Category.ID
	a.Category = ma.Category.Name
	a.ActivityID = ma.Activity.ID
	a.Activity = ma.Activity.Name
	a.TypeID = int(ma.Type.ID.Int64)
	a.Type = ma.Type.Name
	a.Description = ma.Description

	return a, nil
}

// addMemberActivity adds a localMember activity
func addMemberActivity(memberID int, activityInput memberActivityInput) (memberActivity, error) {

	// Create the required type for the insert
	// todo: add evidence and attachment
	ma := activity.MemberActivityInput{
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	// A return value for the new record
	var mar memberActivity

	// This just returns the new record id, so re-fetch the localMember activity record
	// so that all the fields are populated for the response.
	newID, err := activity.AddMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return memberActivityData(memberID, newID)
}

// updateMemberActivity updates an existing localMember activity record
func updateMemberActivity(memberID int, activityInput memberActivityInput) (memberActivity, error) {

	// Create the required value
	ma := activity.MemberActivityInput{
		ID:          activityInput.ID,
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	// A return value for the new record
	var mar memberActivity

	// This just returns an error so re-fetch the localMember activity record
	// so that all the fields are populated for the response.
	err := activity.UpdateMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return memberActivityData(memberID, ma.ID)
}

func deleteMemberActivity(memberID, activityID int) error {
	return activity.DeleteMemberActivity(memberID, activityID)
}

// memberActivityDuplicate returns the id of a matching localMember activity, or 0 if not found
func memberActivityDuplicate(memberID int, activityInput memberActivityInput) int {

	// Create the required value
	ma := activity.MemberActivityInput{
		ID:          activityInput.ID,
		MemberID:    memberID,
		ActivityID:  activityInput.ActivityID,
		TypeID:      activityInput.TypeID,
		Date:        activityInput.Date,
		Quantity:    activityInput.Quantity,
		Description: activityInput.Description,
	}

	return activity.DuplicateMemberActivity(ma)
}

// memberEvaluationsData fetches evaluation data for a localMember.
func memberEvaluationsData(memberID int) ([]memberEvaluation, error) {

	var xme []memberEvaluation

	// This returns a nested struct which is simplified below.
	xma, err := activity.MemberActivityReports(memberID)

	for _, v := range xma {
		e := memberEvaluation{
			ReportName:     v.ReportName,
			StartDate:      v.StartDate,
			EndDate:        v.EndDate,
			CreditRequired: float64(v.CreditRequired),
			CreditObtained: float64(v.CreditObtained),
			Closed:         v.Closed,
		}
		xme = append(xme, e)
	}

	return xme, err
}

// memberCurrentEvaluationData fetches the current evaluation period data for a localMember
func memberCurrentEvaluationData(memberID int) (memberEvaluation, error) {

	var me memberEvaluation

	// This returns a nested struct which is simplified below.
	ce, err := activity.CurrentMemberActivityReport(memberID)
	if err != nil {
		return me, err
	}

	me.ReportName = ce.ReportName
	me.StartDate = ce.StartDate
	me.EndDate = ce.EndDate
	me.CreditRequired = float64(ce.CreditRequired)
	me.CreditObtained = float64(ce.CreditObtained)
	me.Closed = ce.Closed

	return me, nil
}

// memberActivityAttachmentsData fetches the attachments for a localMember activity
func memberActivityAttachmentsData(memberActivityID int) ([]attachments.Attachment, error) {

	return attachments.MemberActivityAttachments(memberActivityID)
}

// memberActivityAttachmentRequest requests a signed URL for uploading to S3
//func memberActivityAttachmentRequest(memberID int) string {
//
//	var url string
//
//	// Get the file set data
//	fs, err := fileset.ActivityAttachment()
//	if err != nil {
//		msg := "Could not determine the storage information for activity attachments - " + err.Error()
//		return msg
//	}
//	fmt.Println(fs)
//
//	// Use the file set information to create an upload value
//
//	return url
//}
