package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var packagesQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"packages": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.PackageType)),
			Description: "Get Packages list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterPackagesType,
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
				var name, packagesType *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					takeType, typeOk := takeFilter["type"].(string)
					if nameOk {
						name = &takeName
					}
					if typeOk {
						packagesType = &takeType
					}
				}

				_Response, err := svcs.PackageServices.ReadPackages(p.Context, filterQuery, name, packagesType)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBPackageToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"packagesLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get packages Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterPackagesType,
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
				var name, packagesType *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					takeType, typeOk := takeFilter["type"].(string)
					if nameOk {
						name = &takeName
					}
					if typeOk {
						packagesType = &takeType
					}
				}

				_Response, err := svcs.PackageServices.ReadPackagesLength(p.Context, filterQuery, name, packagesType)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"package": {
			Type:        schemas.PackageType,
			Description: "Get single package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				packageID := p.Args["packageId"].(string)

				_Response, err := svcs.PackageServices.ReadPackage(p.Context, packageID)
				if err != nil {
					return nil, err
				}
				return transformations.DBPackageToGQLUser(_Response), nil
			},
		},
	}
}

var packagesMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createPackage": {
			Type:        graphql.NewNonNull(schemas.PackageType),
			Description: "Create package",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"amountPerMonth": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"amountPerYear": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"packageServices": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(schemas.CustomPackageServices)),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					name := p.Args["name"].(string)
					customPackages := p.Args["packageServices"].([]interface{})
					takeAmountPerMonth, amountPerMonthOk := p.Args["amountPerMonth"].(int)
					takeAmountPerYear, amountPerYearOk := p.Args["amountPerYear"].(int)
					takeDescription, descriptionOk := p.Args["description"].(string)

					var customPackageConvertList []services.CustomPackageService
					for _, customPackage := range customPackages {
						h := customPackage.(map[string]interface{})
						customPackageConvert := services.CustomPackageService{
							ServiceId: h["serviceId"].(string),
						}
						if h["quantity"] == nil {
							b := h["isActive"].(bool)
							customPackageConvert.IsActive = &b
						} else {
							q := h["quantity"].(int)
							customPackageConvert.Quantity = &q
						}
						customPackageConvertList = append(customPackageConvertList, customPackageConvert)
					}

					var amountPerMonth, amountPerYear *int
					var description *string

					//validations
					if amountPerMonthOk {
						amountPerMonth = &takeAmountPerMonth
					} else {
						amountPerMonth = nil
					}

					//validations
					if amountPerYearOk {
						amountPerYear = &takeAmountPerYear
					} else {
						amountPerYear = nil
					}

					//validations
					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					_Response, err := svcs.PackageServices.CreatePackage(p.Context, name, adminData.ID, amountPerMonth, description, amountPerYear, customPackageConvertList)
					if err != nil {
						return nil, err
					}
					return transformations.DBPackageToGQLUser(_Response), nil
				},
			),
		},
		"createCustomPackage": {
			Type:        graphql.NewNonNull(schemas.PackageType),
			Description: "Create Custom package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"packageServices": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.NewList(schemas.CustomPackageServices)),
				},
			},
			Resolve: utils.AuthenticateUser(
				func(p graphql.ResolveParams, userData *utils.UserFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)
					customPackages := p.Args["packageServices"].([]interface{})
					takeName, nameOk := p.Args["name"].(string)
					takeDescription, descriptionOk := p.Args["description"].(string)

					var customPackageConvertList []services.CustomPackageService
					for _, customPackage := range customPackages {
						h := customPackage.(map[string]interface{})
						customPackageConvert := services.CustomPackageService{
							ServiceId: h["serviceId"].(string),
						}
						if h["quantity"] == nil {
							b := h["isActive"].(bool)
							customPackageConvert.IsActive = &b
						} else {
							q := h["quantity"].(int)
							customPackageConvert.Quantity = &q
						}
						customPackageConvertList = append(customPackageConvertList, customPackageConvert)
					}
					var description, name *string

					//validations
					if nameOk {
						name = &takeName
					} else {
						name = nil
					}

					//validations
					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					_Response, err := svcs.PackageServices.CreateCustomPackage(p.Context, userData.ID, packageID, customPackageConvertList, name, description)
					if err != nil {
						return nil, err
					}
					return transformations.DBPackageToGQLUser(_Response), nil
				},
			),
		},
		"approveCustomPackage": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Approve Custom Package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"amountPerMonth": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"amountPerYear": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)
					amountPerMonth := p.Args["amountPerMonth"].(int)
					amountPerYear := p.Args["amountPerYear"].(int)

					_Response, err := svcs.PackageServices.ApprovePackage(p.Context, packageID, adminData.ID, amountPerMonth, amountPerYear)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"updatePackage": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"amountPerMonth": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"amountPerYear": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
				"description": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)
					takeName, nameOk := p.Args["name"].(string)
					takeAmountPerMonth, amountPerMonthOk := p.Args["amountPerMonth"].(int)
					takeAmountPerYear, amountPerYearOk := p.Args["amountPerYear"].(int)
					takeDescription, descriptionOk := p.Args["description"].(string)

					var amountPerMonth, amountPerYear *int
					var name, description *string

					//validations
					if nameOk {
						name = &takeName
					} else {
						name = nil
					}

					//validations
					if amountPerMonthOk {
						amountPerMonth = &takeAmountPerMonth
					} else {
						amountPerMonth = nil
					}

					//validations
					if amountPerYearOk {
						amountPerYear = &takeAmountPerYear
					} else {
						amountPerYear = nil
					}

					//validations
					if descriptionOk {
						description = &takeDescription
					} else {
						description = nil
					}

					_Response, err := svcs.PackageServices.UpdatePackage(p.Context, packageID, name, description, amountPerMonth, amountPerYear)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deletePackage": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete Package",
			Args: graphql.FieldConfigArgument{
				"packageId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)

					_Response, err := svcs.PackageServices.DeletePackage(p.Context, packageID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposePackageResolver exposes the packages Reesolver
func ExposePackageResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    packagesQuery(services),
		Mutation: packagesMutation(services),
	}
}
