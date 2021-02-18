package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var legalAreasQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"legalAreas": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.LegalAreaType)),
			Description: "Get legal areas list",
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
				var name, description *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					takeDescription, descriptionOk := takeFilter["description"].(string)
					if nameOk {
						name = &takeName
					}
					if descriptionOk {
						description = &takeDescription
					}
				}

				_Response, err := svcs.LegalAreaServices.ReadLegalAreas(p.Context, filterQuery, name, description)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBLegalAreaToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"legalAreasLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get legal areas Counts",
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
				var name, description *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					takeDescription, descriptionOk := takeFilter["description"].(string)
					if nameOk {
						name = &takeName
					}
					if descriptionOk {
						description = &takeDescription
					}
				}

				_Response, err := svcs.LegalAreaServices.ReadLegalAreasLength(p.Context, filterQuery, name, description)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"legalArea": {
			Type:        schemas.LegalAreaType,
			Description: "Get single legal area",
			Args: graphql.FieldConfigArgument{
				"legalAreaId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				legalAreaID := p.Args["legalAreaId"].(string)

				_Response, err := svcs.LegalAreaServices.ReadLegalArea(p.Context, legalAreaID)
				if err != nil {
					return nil, err
				}
				return transformations.DBLegalAreaToGQLUser(_Response), nil
			},
		},
	}
}

var legalAreaMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createLegalArea": {
			Type:        graphql.NewNonNull(schemas.LegalAreaType),
			Description: "Create legal area",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					name := p.Args["name"].(string)
					takeDescription, descriptionOk := p.Args["description"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var description, image *string

					//validations
					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.LegalAreaServices.CreateLegalArea(p.Context, name, description, image, adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBLegalAreaToGQLUser(_Response), nil
				},
			),
		},
		"updateLegalArea": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update legal area",
			Args: graphql.FieldConfigArgument{
				"legalAreaId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					legalAreaID := p.Args["legalAreaId"].(string)
					takeName, nameOk := p.Args["name"].(string)
					takeDescription, descriptionOk := p.Args["description"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var description, image, name *string

					//validations
					if nameOk {
						name = &takeName
					} else {
						name = nil
					}

					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.LegalAreaServices.UpdateLegalArea(p.Context, legalAreaID, name, description, image)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deleteLegalArea": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update legal area",
			Args: graphql.FieldConfigArgument{
				"legalAreaId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					legalAreaID := p.Args["legalAreaId"].(string)

					_Response, err := svcs.LegalAreaServices.DeleteLegalArea(p.Context, legalAreaID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeLegalAreaResolver exposes the legal ares resolver
func ExposeLegalAreaResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    legalAreasQuery(services),
		Mutation: legalAreaMutation(services),
	}
}
