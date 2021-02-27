package schemas

import "github.com/graphql-go/graphql"

var EnumTypeTestBlogPostStatus = graphql.NewEnum(graphql.EnumConfig{
	Name: "BlogPostStatus",
	Values: graphql.EnumValueConfigMap{
		"Active": &graphql.EnumValueConfig{
			Value: "Active",
		},
		"Draft": &graphql.EnumValueConfig{
			Value: "Draft",
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
				Type: EnumTypeTestBlogPostStatus,
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

//FilterBlogPostsType  for filtering blog post type
var FilterBlogPostsType = graphql.NewInputObject(
	graphql.InputObjectConfig{
		Name: "GetTagsFilter",
		Fields: graphql.InputObjectConfigFieldMap{
			"title": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"details": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"status": &graphql.InputObjectFieldConfig{
				Type: EnumTypeTestBlogPostStatus,
			},
			"tag": &graphql.InputObjectFieldConfig{
				Type: graphql.ID,
			},
			"search": &graphql.InputObjectFieldConfig{
				Type: SearchType,
			},
			"order": &graphql.InputObjectFieldConfig{
				Type: enumTypeForOrder,
			},
			"orderBy": &graphql.InputObjectFieldConfig{
				Type: graphql.String,
			},
			"dateRange": &graphql.InputObjectFieldConfig{
				Type: DateRangeType,
			},
		},
	},
)
