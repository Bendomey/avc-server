package schemas

import "github.com/graphql-go/graphql"

// LawyerType defines typings for lawyers
var LawyerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Lawyer",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"name": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"user": &graphql.Field{
				Type: UserType,
			},
			"digitalAddress": &graphql.Field{
				Type: graphql.String,
			},
			"addressCountry": &graphql.Field{
				Type: graphql.String,
			},
			"addressCity": &graphql.Field{
				Type: graphql.String,
			},
			"addressStreetNumber": &graphql.Field{
				Type: graphql.String,
			},
			"addressNumber": &graphql.Field{
				Type: graphql.String,
			},
			"tin": &graphql.Field{
				Type: graphql.String,
			},
			"licenseNumber": &graphql.Field{
				Type: graphql.String,
			},
			"firstYearOfBarAdmission": &graphql.Field{
				Type: graphql.String,
			},
			"nationalIDFront": &graphql.Field{
				Type: graphql.String,
			},
			"nationalIDBack": &graphql.Field{
				Type: graphql.String,
			},
			"barMembershipCard": &graphql.Field{
				Type: graphql.String,
			},
			"lawCertificate": &graphql.Field{
				Type: graphql.String,
			},
			"cv": &graphql.Field{
				Type: graphql.String,
			},
			"coverLetter": &graphql.Field{
				Type: graphql.String,
			},
			"suspendedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"suspendedBy": &graphql.Field{
				Type: AdminType,
			},
			"approvedAt": &graphql.Field{
				Type: graphql.DateTime,
			},
			"approvedBy": &graphql.Field{
				Type: AdminType,
			},
			"createdAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
			"updatedAt": &graphql.Field{
				Type: graphql.NewNonNull(graphql.DateTime),
			},
		},
	},
)
