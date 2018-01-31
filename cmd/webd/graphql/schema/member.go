package schema

import (
	"errors"
	"fmt"
	"time"

	"github.com/graphql-go/graphql"

	"github.com/mappcpd/web-services/internal/members"
	"github.com/mappcpd/web-services/internal/platform/jwt"
	"github.com/mappcpd/web-services/internal/utility"
)

// Member struct - a simpler representation than members.Member
type Member struct {
	ID             int                      `json:"id"`
	Active         bool                     `json:"active"`
	Title          string                   `json:"title"`
	FirstName      string                   `json:"firstName"`
	MiddleNames    string                   `json:"middleNames"`
	LastName       string                   `json:"lastName"`
	PostNominal    string                   `json:"postNominal"`
	DateOfBirth    string                   `json:"dateOfBirth"`
	Email          string                   `json:"email"`
	Mobile         string                   `json:"mobile"`
	Locations      []members.MemberLocation `json:"locations"`
	Qualifications []members.Qualification  `json:"qualifications"`
	Positions      []members.Position       `json:"positions"`
}

// MemberActivity is a simpler representation of the member activity than the nested one in the current REST api.
type MemberActivity struct {
	ID          int       `json:"id"`
	Date        string    `json:"date"`
	DateTime    time.Time `json:"dateTime"`
	Credit      float64   `json:"credit"`
	CategoryID  int       `json:"categoryId"`
	Category    string    `json:"category"`
	TypeID      int       `json:"typeId"`
	Type        string    `json:"type"`
	Description string    `json:"description"`
}

// MemberEvaluation representations the member evaluation data
type MemberEvaluation struct {
	//ID          int       `json:"id"`
	Name           string  `json:"name"`
	StartDate      string  `json:"startDate"`
	EndDate        string  `json:"endDate"`
	CreditRequired float64 `json:"creditRequired"`
	CreditObtained float64 `json:"creditObtained"`
	Closed         bool    `json:"closed"`
}

// getMember fetches the basic member record
func getMember(id int) (Member, error) {
	var m Member
	mp, err := getMemberProfile(id)
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

// getMemberProfile fetches a single member record by id
func getMemberProfile(memberID int) (members.Member, error) {
	// MemberByID returns a pointer to a members.Member so dereference in return
	m, err := members.MemberByID(memberID)
	return *m, err
}

// getMemberActivities fetches activities for a member. By default it returns the entire set,
// ordered by activity date desc. Some filters have been added here for the caller's convenience.
func getMemberActivities(memberID int, filter map[string]interface{}) ([]MemberActivity, error) {
	var xa []MemberActivity

	// This returns a nested struct which is simplified below.
	xma, err := members.MemberActivitiesByMemberID(memberID)

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
		a := MemberActivity{
			ID:          v.ID,
			Date:        v.Date,
			DateTime:    v.DateISO,
			Credit:      v.Credit,
			CategoryID:  v.Category.ID,
			Category:    v.Category.Name,
			TypeID:      v.Activity.ID,
			Type:        v.Activity.Name,
			Description: v.Description,
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

// unpack an object into a value of type MemberActivity
func (ma *MemberActivity) unpack(obj map[string]interface{}) error {
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
	if val, ok := obj["credit"].(float64); ok {
		ma.Credit = val
	}
	if val, ok := obj["categoryId"].(int); ok {
		ma.CategoryID = int(val)
	}
	if val, ok := obj["typeId"].(int); ok {
		ma.TypeID = int(val)
	}
	if val, ok := obj["description"].(string); ok {
		ma.Description = val
	}

	return nil
}

// getMemberActivity fetches a single activities by id.
// It verifies that the activity is owned by the member by memberID.
func getMemberActivity(memberID, activityID int) (MemberActivity, error) {

	var a MemberActivity

	// This returns a nested struct which we can simplify
	ma, err := members.MemberActivityByID(activityID)
	if err != nil {
		return a, err
	}

	// Verify owner match
	if ma.MemberID != memberID {
		msg := fmt.Sprintf("MemberActivity with id %v does not belong to member with id %v", activityID, memberID)
		return a, errors.New(msg)
	}

	a.ID = ma.ID
	a.Date = ma.Date
	a.DateTime = ma.DateISO
	a.Credit = ma.Credit
	a.CategoryID = ma.Category.ID
	a.Category = ma.Category.Name
	a.TypeID = ma.Activity.ID
	a.Type = ma.Activity.Name
	a.Description = ma.Description

	return a, nil
}

// addMemberActivity adds a member activity
func addMemberActivity(memberID int, memberActivity MemberActivity) (MemberActivity, error) {

	// Create the required type for the insert
	// todo: add evidence and categoryId
	ma := members.MemberActivityRow{
		MemberID:    memberID,
		ActivityID:  memberActivity.TypeID,
		Date:        memberActivity.Date,
		Quantity:    memberActivity.Credit,
		Description: memberActivity.Description,
	}

	// A return value for the new record
	var mar MemberActivity

	// This just returns the new record id, so re-fetch the member activity record
	// so that all the fields are populated for the response.
	newID, err := members.AddMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return getMemberActivity(memberID, newID)

}

// updateMemberActivity adds a member activity
func updateMemberActivity(memberID int, memberActivity MemberActivity) (MemberActivity, error) {

	// Create the required type for the insert
	ma := members.MemberActivityRow{
		MemberID:    memberID,
		ID:          memberActivity.ID,     // id of the activity instance
		ActivityID:  memberActivity.TypeID, // id of the activity type
		Date:        memberActivity.Date,
		Quantity:    memberActivity.Credit,
		Description: memberActivity.Description,
	}

	// A return value for the new record
	var mar MemberActivity

	// This just returns an error so re-fetch the member activity record
	// so that all the fields are populated for the response.
	err := members.UpdateMemberActivity(ma)
	if err != nil {
		return mar, err
	}

	return getMemberActivity(memberID, ma.ID)

}

// GetMemberEvaluations fetches evaluation data for a member.
func GetMemberEvaluations(memberID int) ([]MemberEvaluation, error) {

	var xme []MemberEvaluation

	// This returns a nested struct which is simplified below.
	xma, err := members.EvaluationsByMemberID(memberID)

	for _, v := range xma {
		e := MemberEvaluation{
			Name:           v.Name,
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

// getCurrentEvaluation fetches the current evaluation period data for a member.
func getCurrentEvaluation(memberID int) (MemberEvaluation, error) {

	var me MemberEvaluation

	// This returns a nested struct which is simplified below.
	ce, err := members.CurrentEvaluation(memberID)
	if err != nil {
		return me, err
	}

	me.Name = ce.Name
	me.StartDate = ce.StartDate
	me.EndDate = ce.EndDate
	me.CreditRequired = float64(ce.CreditRequired)
	me.CreditObtained = float64(ce.CreditObtained)
	me.Closed = ce.Closed

	return me, nil
}

// MemberUser is exported field attached to the root query. It is a top-level 'viewer' query field that ensures data
// is restricted to the member (user) identified by the token.
var MemberUser = &graphql.Field{
	Description: "The memberUser field acts as a 'viewer' and requires a valid JSON Web Token ( see https://jwt.io). " +
		"Data in child fields will always belong to the member identified by the token.",
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
			// Validate the token, and extract the member id
			at, err := jwt.Check(token)
			if err != nil {
				return nil, err
			}
			//fmt.Println(at.Claims)
			id := at.Claims.ID
			// At this point we have a valid token from which we've extracted an id.
			// As a final step we can verify that the id is a valid user in the system,
			// for example, that it is active. Although this is a bit redundant for each request?
			return getMember(id)
		}
		return nil, nil
	},
}

// memberType represents the memberUser node, and provides a path to the child nodes.
var memberType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "member",
	Description: "Member query object that provides access to data for the member identified by the token.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.String,
			Description: "The member's unique id number",
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
			Type:        graphql.NewList(memberLocationType),
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

		// sub queries
		"activity":    memberActivity,
		"activities":  memberActivities,
		"evaluation":  memberCurrentEvaluation,
		"evaluations": memberEvaluations,
	},
})

// Contact represents a contact 'card' - that is, a single contact record that pertains to a Member.
var memberContactType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "contact",
	Description: "A contact record belonging to a member",
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
			Type: graphql.NewList(memberLocationType),
		},
	},
})

// location represents one or more contact locations pertaining to a member
var memberLocationType = graphql.NewObject(graphql.ObjectConfig{
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

// position represents a position held by a member
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

// qualification represents a qualification obtained by the member
var qualificationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "qualification",
	Description: "An academic qualification obtained by the member",
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

// memberActivity represents a Member memberActivity record (not memberActivity type record)
var memberActivityType = graphql.NewObject(graphql.ObjectConfig{
	Name: "memberActivity",
	Description: "An activity record belonging to a member. This is an instance of an activity recorded " +
		"by a member, having been completed on a particular date, with additional information such as duration and description.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the member activity record",
		},
		"date": &graphql.Field{
			Type:        graphql.String,
			Description: "The date the activity was undertaken, as string format 'YYYY-MM-DD'.",
		},
		"dateTime": &graphql.Field{
			Type: graphql.DateTime,
			Description: "The date the activity was undertaken. Note only a date string is required, eg '2017-12-07' and " +
				"any time information is discarded. This field returns the date in RFC3339 format with the time set " +
				"to 00:00:00 UTC to facilitate date ordering and other date-related operations.",
		},
		"credit": &graphql.Field{
			Type:        graphql.Float,
			Description: "Value or credit for the memberActivity",
		},
		"categoryId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The memberActivity category id",
		},
		"category": &graphql.Field{
			Type:        graphql.String,
			Description: "The top-level category of the memberActivity",
		},
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "The type of memberActivity",
		},
		"typeId": &graphql.Field{
			Type:        graphql.Int,
			Description: "The memberActivity type id",
		},
		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "The specifics of the memberActivity described by the member",
		},
	},
})

// memberActivityInput is an input object type used as an argument for adding / updating a memberActivity
var memberActivityInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name:        "memberActivityInput",
	Description: "An input object type used as an argument for adding / updating a memberActivity",
	Fields: graphql.InputObjectConfigFieldMap{

		// optional member activity id - if supplied then it is an update
		"id": &graphql.InputObjectFieldConfig{
			Type:        graphql.Int,
			Description: "Optional id of the member memberActivity record - if supplied then will update existing.",
		},

		"date": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The date of the memberActivity",
		},

		"credit": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Float},
			Description: "Value or credit for the memberActivity",
		},

		"typeId": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "The memberActivity type id",
		},

		"description": &graphql.InputObjectFieldConfig{
			Type:        &graphql.NonNull{OfType: graphql.String},
			Description: "The specifics of the memberActivity described by the member",
		},
	},
})

// memberEvaluationType represents data about the points credited vs points required for an evaluation period.
var memberEvaluationType = graphql.NewObject(graphql.ObjectConfig{
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

// memberActivity query field fetches a single memberActivity that belongs to a member
var memberActivity = &graphql.Field{
	Description: "Fetches a single member memberActivity by memberActivity id.",
	Type:        memberActivityType,
	Args: graphql.FieldConfigArgument{
		"activityId": &graphql.ArgumentConfig{
			Type:        &graphql.NonNull{OfType: graphql.Int},
			Description: "ID of the member memberActivity",
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

		activityID, ok := p.Args["activityId"].(int)
		if ok {
			return getMemberActivity(memberID, int(activityID))
		}

		return nil, nil
	},
}

// memberActivities field fetches multiple memberActivities belonging to a member
var memberActivities = &graphql.Field{
	Description: "Fetches a list of member memberActivities",
	Type:        graphql.NewList(memberActivityType),
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

		return getMemberActivities(memberID, f)
	},
}

// memberCurrentEvaluation field fetches the current evaluation period data
var memberCurrentEvaluation = &graphql.Field{
	Description: "Fetches activity data for the current evaluation period",
	Type:        memberEvaluationType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return getCurrentEvaluation(memberID)
	},
}

// memberEvaluations field fetches a list of activity evaluations for the member
var memberEvaluations = &graphql.Field{
	Description: "Fetches a history of member activity evaluation periods",
	Type:        graphql.NewList(memberEvaluationType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Check(token.(string))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return GetMemberEvaluations(memberID)
	},
}

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
			//return data.getMember(id)
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
		"belongs to the member identified by the token, then it will be updated. If `activityId` is not present, or does not belong " +
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
			ma := MemberActivity{}
			err := ma.unpack(maObj)
			if err != nil {
				return nil, err
			}

			if ma.ID > 0 {
				return updateMemberActivity(memberID, ma)
			}

			return addMemberActivity(memberID, ma)
		}

		return nil, nil
	},
}
