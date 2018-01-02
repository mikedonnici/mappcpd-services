package schema

import (
	"github.com/graphql-go/graphql"
)

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
