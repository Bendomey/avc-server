package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var servicesQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"services": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.ServiceType)),
			Description: "Get Services list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterServicesType,
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

				_Response, err := svcs.ServiceServices.ReadServices(p.Context, filterQuery, name)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBServiceToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"servicesLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get services Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterServicesType,
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

				_Response, err := svcs.ServiceServices.ReadServicesLength(p.Context, filterQuery, name)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"service": {
			Type:        schemas.ServiceType,
			Description: "Get single service",
			Args: graphql.FieldConfigArgument{
				"serviceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				serviceID := p.Args["serviceId"].(string)

				_Response, err := svcs.ServiceServices.ReadService(p.Context, serviceID)
				if err != nil {
					return nil, err
				}
				return transformations.DBServiceToGQLUser(_Response), nil
			},
		},
	}
}

var servicesMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createService": {
			Type:        graphql.NewNonNull(schemas.ServiceType),
			Description: "Create service",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeForService),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					name := p.Args["name"].(string)
					serviceType := p.Args["type"].(string)
					takePrice, priceOk := p.Args["price"].(float64)
					takeDescription, descriptionrOk := p.Args["description"].(string)

					var description *string
					var price *float64

					//validations
					if priceOk {
						price = &takePrice
					} else {
						price = nil
					}

					if descriptionrOk {
						description = &takeDescription
					} else {
						description = nil
					}

					_Response, err := svcs.ServiceServices.CreateService(p.Context, name, price, description, models.ServiceType(serviceType), adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBServiceToGQLUser(_Response), nil
				},
			),
		},
		"updateService": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Service",
			Args: graphql.FieldConfigArgument{
				"serviceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"type": &graphql.ArgumentConfig{
					Type: schemas.EnumTypeForService,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					serviceID := p.Args["serviceId"].(string)
					takeServiceType, serviceTypeOk := p.Args["type"].(string)
					takeName, nameOk := p.Args["name"].(string)
					takePrice, priceOk := p.Args["price"].(float64)
					takeDescription, descriptionrOk := p.Args["description"].(string)

					var name, description *string
					var price *float64
					var serviceType *string

					//validations
					if nameOk {
						name = &takeName
					} else {
						name = nil
					}

					//validations
					if priceOk {
						price = &takePrice
					} else {
						price = nil
					}

					if descriptionrOk {
						description = &takeDescription
					} else {
						description = nil
					}

					if serviceTypeOk {
						serviceType = &takeServiceType
					} else {
						serviceType = nil
					}

					_Response, err := svcs.ServiceServices.UpdateService(p.Context, serviceID, name, price, description, serviceType)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deleteService": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete Service",
			Args: graphql.FieldConfigArgument{
				"serviceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					serviceID := p.Args["serviceId"].(string)

					_Response, err := svcs.ServiceServices.DeleteService(p.Context, serviceID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeServiceResolver exposes the Services Reesolver
func ExposeServiceResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    servicesQuery(services),
		Mutation: servicesMutation(services),
	}
}
