package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var blogPostsQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"posts": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.BlogPostType)),
			Description: "Get blog posts list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterBlogPostsType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				argument := p.Args
				filterQuery, filterErr := utils.GenerateQuery(argument)
				if filterErr != nil {
					return nil, filterErr
				}

				//fields
				takeFilter, filterOk := argument["filter"].(map[string]interface{})
				var title, details, status, tag *string

				if filterOk {
					takeTitle, titleOk := takeFilter["title"].(string)
					if titleOk {
						title = &takeTitle
					}

					takeDetails, detailsOk := takeFilter["details"].(string)
					if detailsOk {
						details = &takeDetails
					}

					takeStatus, statusOk := takeFilter["status"].(string)
					if statusOk {
						status = &takeStatus
					}

					takeTag, tagOk := takeFilter["tag"].(string)
					if tagOk {
						tag = &takeTag
					}
				}

				_Response, err := svcs.BlogPostServices.ReadPosts(p.Context, filterQuery, status, title, details, tag)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBBlogPostToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"postsLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get blog posts Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterBlogPostsType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				argument := p.Args
				filterQuery, filterErr := utils.GenerateQuery(argument)
				if filterErr != nil {
					return nil, filterErr
				}
				//fields
				takeFilter, filterOk := argument["filter"].(map[string]interface{})
				var title, details, status, tag *string

				if filterOk {
					takeTitle, titleOk := takeFilter["title"].(string)
					if titleOk {
						title = &takeTitle
					}

					takeDetails, detailsOk := takeFilter["details"].(string)
					if detailsOk {
						details = &takeDetails
					}

					takeStatus, statusOk := takeFilter["status"].(string)
					if statusOk {
						status = &takeStatus
					}

					takeTag, tagOk := takeFilter["tag"].(string)
					if tagOk {
						tag = &takeTag
					}
				}

				_Response, err := svcs.BlogPostServices.ReadPostsLength(p.Context, filterQuery, status, title, details, tag)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"post": {
			Type:        schemas.BlogPostType,
			Description: "Get single blog post",
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				postID := p.Args["postId"].(string)

				_Response, err := svcs.BlogPostServices.ReadPost(p.Context, postID)
				if err != nil {
					return nil, err
				}
				return transformations.DBBlogPostToGQLUser(_Response), nil
			},
		},
	}
}

var blogPostsMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createPost": {
			Type:        graphql.NewNonNull(schemas.BlogPostType),
			Description: "Create blog post",
			Args: graphql.FieldConfigArgument{
				"title": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"details": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"tag": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"status": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeTestBlogPostStatus),
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					title := p.Args["title"].(string)
					details := p.Args["details"].(string)
					tag := p.Args["tag"].(string)
					status := p.Args["status"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var image *string
					//validations
					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.BlogPostServices.CreatePost(p.Context, title, details, tag, status, image, adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBBlogPostToGQLUser(_Response), nil
				},
			),
		},
		"updatePost": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update blog post",
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"title": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"details": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"tag": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
				"status": &graphql.ArgumentConfig{
					Type: schemas.EnumTypeTestBlogPostStatus,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					postID := p.Args["postId"].(string)
					takeTitle, titleOk := p.Args["title"].(string)
					takeDetails, detailsOk := p.Args["details"].(string)
					takeTag, tagOk := p.Args["tag"].(string)
					takeStatus, statusOk := p.Args["status"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var title, details, tag, status, image *string

					//validations
					if titleOk {
						title = &takeTitle
					} else {
						title = nil
					}

					if detailsOk {
						details = &takeDetails
					} else {
						details = nil
					}

					if tagOk {
						tag = &takeTag
					} else {
						tag = nil
					}

					if statusOk {
						status = &takeStatus
					} else {
						status = nil
					}

					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.BlogPostServices.UpdatePost(p.Context, postID, title, details, tag, status, image)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deletePost": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete blog post",
			Args: graphql.FieldConfigArgument{
				"postId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					postID := p.Args["postId"].(string)

					_Response, err := svcs.BlogPostServices.DeletePost(p.Context, postID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeBlogPostResolver exposes the blog posts Reesolver
func ExposeBlogPostResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    blogPostsQuery(services),
		Mutation: blogPostsMutation(services),
	}
}
