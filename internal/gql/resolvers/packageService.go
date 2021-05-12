package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var packageServicesQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"packageServices": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.PackageServiceType)),
			Description: "Get Package Services list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterPackageServicesType,
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
				var serviceID, packageID *string

				if filterOk {
					takeService, serviceOk := takeFilter["service"].(string)

					if serviceOk {
						serviceID = &takeService
					}

					takePackage, packageOk := takeFilter["package"].(string)

					if packageOk {
						packageID = &takePackage
					}
				}

				_Response, err := svcs.PackageServiceServices.ReadPackageServices(p.Context, filterQuery, serviceID, packageID)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}

				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBPackageServiceToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"packageServicesLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get package services Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterPackageServicesType,
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
				var serviceID, packageID *string

				if filterOk {
					takeService, serviceOk := takeFilter["service"].(string)

					if serviceOk {
						serviceID = &takeService
					}

					takePackage, packageOk := takeFilter["package"].(string)

					if packageOk {
						packageID = &takePackage
					}
				}
				_Response, err := svcs.PackageServiceServices.ReadPackageServicesLength(p.Context, filterQuery, serviceID, packageID)

				if err != nil {
					return nil, err
				}

				return _Response, nil
			},
		},
		"packageService": {
			Type:        schemas.PackageServiceType,
			Description: "Get single package service",
			Args: graphql.FieldConfigArgument{
				"packageServiceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				packageServiceID := p.Args["packageServiceId"].(string)

				_Response, err := svcs.PackageServiceServices.ReadPackageService(p.Context, packageServiceID)
				if err != nil {
					return nil, err
				}
				return transformations.DBPackageServiceToGQLUser(_Response), nil
			},
		},
	}
}

var packageServicesMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createPackageService": {
			Type:        graphql.NewNonNull(schemas.PackageServiceType),
			Description: "Create package service",
			Args: graphql.FieldConfigArgument{
				"serviceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"type": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeForPackageService),
				},
				"quantity": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"isActive": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					serviceID := p.Args["serviceId"].(string)
					packageID := p.Args["packageId"].(string)
					packageServiceType := p.Args["type"].(string)
					takeQuantity, quantityOk := p.Args["quantity"].(int)
					takeIsActive, isActiveOk := p.Args["isActive"].(bool)

					var quantity *int
					var isActive *bool

					//validations
					if quantityOk {
						quantity = &takeQuantity
					} else {
						quantity = nil
					}

					if isActiveOk {
						isActive = &takeIsActive
					} else {
						isActive = nil
					}

					_Response, err := svcs.PackageServiceServices.CreatePackageService(p.Context, serviceID, packageID, models.PackageServiceType(packageServiceType), quantity, isActive, adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBPackageServiceToGQLUser(_Response), nil
				},
			),
		},
		"updatePackageService": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Package Service",
			Args: graphql.FieldConfigArgument{
				"packageServiceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"serviceId": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.ID,
				},
				"type": &graphql.ArgumentConfig{
					Type: schemas.EnumTypeForPackageService,
				},
				"quantity": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"isActive": &graphql.ArgumentConfig{
					Type: graphql.Boolean,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageServiceID := p.Args["packageServiceId"].(string)
					takePackageServiceType, packageServiceTypeOk := p.Args["type"].(string)
					takeServiceID, serviceIDOk := p.Args["serviceId"].(string)
					takePackageID, packageIDOk := p.Args["packageId"].(string)
					takeQuantity, quantityOk := p.Args["quantity"].(int)
					takeIsActive, isActiveOk := p.Args["isActive"].(bool)

					var packageID, serviceID, packageServiceType *string
					var quantity *int
					var isActive *bool

					//validations
					if serviceIDOk {
						serviceID = &takeServiceID
					} else {
						serviceID = nil
					}

					//validations
					if packageIDOk {
						packageID = &takePackageID
					} else {
						packageID = nil
					}

					if quantityOk {
						quantity = &takeQuantity
					} else {
						quantity = nil
					}

					if isActiveOk {
						isActive = &takeIsActive
					} else {
						isActive = nil
					}

					if packageServiceTypeOk {
						packageServiceType = &takePackageServiceType
					} else {
						packageServiceType = nil
					}

					_Response, err := svcs.PackageServiceServices.UpdatePackageService(p.Context, packageServiceID, serviceID, packageID, packageServiceType, quantity, isActive)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deletePackageService": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete Package Service",
			Args: graphql.FieldConfigArgument{
				"packageServiceId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageServiceID := p.Args["packageServiceId"].(string)

					_Response, err := svcs.PackageServiceServices.DeletePackageService(p.Context, packageServiceID)
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
func ExposePackageServiceResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    packageServicesQuery(services),
		Mutation: packageServicesMutation(services),
	}
}
