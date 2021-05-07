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
				var name *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					if nameOk {
						name = &takeName
					}
				}

				_Response, err := svcs.PackageServices.ReadPackages(p.Context, filterQuery, name)
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
				var name *string

				if filterOk {
					takeName, nameOk := takeFilter["name"].(string)
					if nameOk {
						name = &takeName
					}
				}

				_Response, err := svcs.PackageServices.ReadPackagesLength(p.Context, filterQuery, name)
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
				"tagId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				tagID := p.Args["packageId"].(string)

				_Response, err := svcs.PackageServices.ReadPackage(p.Context, tagID)
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
			Description: "Create tag",
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
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					name := p.Args["name"].(string)
					takeAmountPerMonth, amountPerMonthOk := p.Args["amountPerMonth"].(int)
					takeAmountPerYear, amountPerYearOk := p.Args["amountPerYear"].(int)

					var amountPerMonth, amountPerYear *int

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

					_Response, err := svcs.PackageServices.CreatePackage(p.Context, name, adminData.ID, amountPerMonth, amountPerYear)
					if err != nil {
						return nil, err
					}
					return transformations.DBPackageToGQLUser(_Response), nil
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
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					packageID := p.Args["packageId"].(string)
					takeName, nameOk := p.Args["name"].(string)
					takeAmountPerMonth, amountPerMonthOk := p.Args["amountPerMonth"].(int)
					takeAmountPerYear, amountPerYearOk := p.Args["amountPerYear"].(int)

					var amountPerMonth, amountPerYear *int
					var name *string

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

					_Response, err := svcs.PackageServices.UpdatePackage(p.Context, packageID, name, amountPerMonth, amountPerYear)
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
