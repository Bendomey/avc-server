package schemas

import "github.com/graphql-go/graphql"

//PaginationType  for pagination type
var PaginationType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Pagination",
		Fields: graphql.InputObjectConfigFieldMap{
			"limit": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
			"skip": &graphql.InputObjectFieldConfig{
				Type: graphql.Int,
			},
		},
	},
)

//enumTypeForOrder for ordering
var enumTypeForOrder = graphql.NewEnum(graphql.EnumConfig{
	Name: "Order",
	Values: graphql.EnumValueConfigMap{
		"asc": &graphql.EnumValueConfig{
			Value: "asc",
		},
		"desc": &graphql.EnumValueConfig{
			Value: "desc",
		},
	},
})

//DateRangeType  for date range type
var DateRangeType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "DateRange",
		Fields: graphql.InputObjectConfigFieldMap{
			"start": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"end": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)

//SearchType  for searching
var SearchType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "Search",
		Fields: graphql.InputObjectConfigFieldMap{
			"criteria": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"searchFields": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String))),
			},
		},
	},
)
