package schemas

import "github.com/graphql-go/graphql"

var enumTypeTestBlogPostStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "BlogPostStatus",
	Values: graphql.EnumValueConfigMap{
		"Active": &graphql.EnumValueConfig{
			Value: 0,
		},
		"Draft": &graphql.EnumValueConfig{
			Value: 1,
		},
	},
})

// BlogPostType defines typings for blog post
var BlogPostType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "BlogPost",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.NewNonNull(graphql.ID),
			},
			"title": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"image": &graphql.Field{
				Type: graphql.String,
			},
			"status": &graphql.Field{
				Type: enumTypeTestBlogPostStatus,
			},
			"tag": &graphql.Field{
				Type: TagType,
			},
			"details": &graphql.Field{
				Type: graphql.NewNonNull(graphql.String),
			},
			"createdBy": &graphql.Field{
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
