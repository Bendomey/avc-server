package schemas

import "github.com/graphql-go/graphql"

// ServicingFieldType defines typings for Servicing Fields - categories
var ServicingFieldType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ServicingField",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			// "description": &graphql.Field{
			// 	Type: BusinessType,
			// },
			// "trademark": &graphql.Field{
			// 	Type: DocumentType,
			// },
			// "document": &graphql.Field{
			// 	Type: DocumentType,
			// },
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)

var BusinessType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Business",
		Fields: graphql.Fields{
			"country": &graphql.Field{
				Type: CountryType,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"owners": &graphql.Field{
				Type: graphql.String,
			},
			"directors": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"numberOfShares": &graphql.Field{
				Type: graphql.String,
			},
			"entityType": &graphql.Field{
				Type: graphql.String,
			},
			"initialCapital": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var TrademarkType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Trademark",
		Fields: graphql.Fields{
			"country": &graphql.Field{
				Type: CountryType,
			},
			"ownershipType": &graphql.Field{
				Type: graphql.String,
			},
			"owners": &graphql.Field{
				Type: graphql.String,
			},
			"address": &graphql.Field{
				Type: graphql.String,
			},
			"classificationOfTrademark": &graphql.Field{
				Type: graphql.String,
			},
			"uploads": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var DocumentType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Document",
		Fields: graphql.Fields{
			"type": &graphql.Field{
				Type: graphql.String,
			},
			"natureOfDoc": &graphql.Field{
				Type: graphql.String,
			},
			"deadline": &graphql.Field{
				Type: graphql.DateTime,
			},
			"existingDocuments": &graphql.Field{
				Type: graphql.String,
			},
			"newDocuments": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)
