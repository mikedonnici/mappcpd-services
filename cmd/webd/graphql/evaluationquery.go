package graphql

import (
	"os"

	"github.com/cardiacsociety/web-services/internal/platform/jwt"
	"github.com/graphql-go/graphql"
)

// currentEvaluationQuery resolves queries for the current evaluation period
var currentEvaluationQuery = &graphql.Field{
	Description: "Fetches activity data for the current evaluation period",
	Type:        evaluationType,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Decode(token.(string), os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return currentEvaluation(memberID)
	},
}

// evaluationsQuery resolves queries for multiple member evaluation periods
var evaluationsQuery = &graphql.Field{
	Description: "Fetches a history of member activity evaluation periods",
	Type:        graphql.NewList(evaluationType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {

		// Extract member id from the token, available thus:
		token := p.Info.VariableValues["token"]
		at, err := jwt.Decode(token.(string), os.Getenv("MAPPCPD_JWT_SIGNING_KEY"))
		if err != nil {
			return nil, err
		}
		memberID := at.Claims.ID

		return evaluations(memberID)
	},
}

// evaluationType defines fields for a member evaluation
var evaluationType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "evaluationData",
	Description: "An evaluation of activity credited and required, for a given period of time - eg a calendar year.",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type:        graphql.Int,
			Description: "The id of the member evaluation",
		},
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
