package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var tagsQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"tags": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.LegalAreaType)),
			Description: "Get tags list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterLegalAreasType,
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
				var name *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					if nameOk {
						name = &takeName
					}
				}

				_Response, err := svcs.TagServices.ReadTags(p.Context, filterQuery, name)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBTagToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"tagsLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get tags Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterLegalAreasType,
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
				var name *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					if nameOk {
						name = &takeName
					}
				}

				_Response, err := svcs.TagServices.ReadTagsLength(p.Context, filterQuery, name)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"tag": {
			Type:        schemas.TagType,
			Description: "Get single tag",
			Args: graphql.FieldConfigArgument{
				"tagId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tagID := p.Args["tagId"].(string)

				_Response, err := svcs.TagServices.ReadTag(p.Context, tagID)
				if err != nil {
					return nil, err
				}
				return transformations.DBTagToGQLUser(_Response), nil
			},
		},
	}
}

var tagsMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{}
}

// ExposeTagResolver exposes the legal ares resolver
func ExposeTagResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    tagsQuery(services),
		Mutation: tagsMutation(services),
	}
}
