package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var countriesQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"countries": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.CountryType)),
			Description: "Get country list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterCountryType,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
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

					_Response, err := svcs.CountryServices.ReadCountries(p.Context, filterQuery, name, description)
					if err != nil {
						return nil, err
					}
					var newResponse []interface{}
					//loop to get the types
					for _, dbRec := range _Response {
						rec := transformations.DBCountryToGQLUser(dbRec)
						newResponse = append(newResponse, rec)
					}
					return newResponse, nil
				},
			),
		},
		"countriesLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get country Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterCountryType,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
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

					_Response, err := svcs.CountryServices.ReadCountriesLength(p.Context, filterQuery, name, description)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"country": {
			Type:        schemas.CountryType,
			Description: "Get single country",
			Args: graphql.FieldConfigArgument{
				"countryId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				adminID := p.Args["countryId"].(string)

				_Response, err := svcs.CountryServices.ReadCountry(p.Context, adminID)
				if err != nil {
					return nil, err
				}
				return transformations.DBCountryToGQLUser(_Response), nil
			},
		},
	}
}

var countryMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createCountry": {
			Type:        graphql.NewNonNull(schemas.CountryType),
			Description: "Create country",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"currency": &graphql.ArgumentConfig{
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
					takeCurrency, currencyOk := p.Args["currency"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var description, currency, image *string

					//validations
					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					if currencyOk {
						currency = &takeCurrency
					} else {
						currency = nil
					}

					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.CountryServices.CreateCountry(p.Context, name, description, currency, image, adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBCountryToGQLUser(_Response), nil
				},
			),
		},
		"updateCountry": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update country",
			Args: graphql.FieldConfigArgument{
				"countryId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"currency": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"image": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					countryID := p.Args["countryId"].(string)
					takeName, nameOk := p.Args["name"].(string)
					takeDescription, descriptionOk := p.Args["description"].(string)
					takeCurrency, currencyOk := p.Args["currency"].(string)
					takeImage, imageOk := p.Args["image"].(string)

					var description, currency, image, name *string

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

					if currencyOk {
						currency = &takeCurrency
					} else {
						currency = nil
					}

					if imageOk {
						image = &takeImage
					} else {
						image = nil
					}

					_Response, err := svcs.CountryServices.UpdateCountry(p.Context, countryID, name, description, currency, image)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deleteCountry": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update country",
			Args: graphql.FieldConfigArgument{
				"countryId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					countryID := p.Args["countryId"].(string)

					_Response, err := svcs.CountryServices.DeleteCountry(p.Context, countryID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeCountryResolver exposes the customer resolver
func ExposeCountryResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    countriesQuery(services),
		Mutation: countryMutation(services),
	}
}
