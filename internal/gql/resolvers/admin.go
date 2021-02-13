package resolvers

import (
	"errors"

	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var adminQuery = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"admins": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.AdminType)),
			Description: "Get Administrators in the system",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterAdminType,
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
					var name *string

					if filterOk {
						takeName, nameOk := takeFilter["fullname"].(string)
						if nameOk {
							name = &takeName
						}

					}

					_Response, err := svcs.AdminServices.ReadAdmins(p.Context, filterQuery, name)
					if err != nil {
						return nil, err
					}
					var newResponse []interface{}
					//loop to get the types
					for _, dbRec := range _Response {
						rec := transformations.DBAdminToGQLAdmin(&dbRec)
						newResponse = append(newResponse, rec)
					}
					return newResponse, nil
				},
			),
		},
		"adminsLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get Length of Administrators in the system",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterAdminType,
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
					var name *string

					if filterOk {
						takeName, nameOk := takeFilter["fullname"].(string)
						if nameOk {
							name = &takeName
						}

					}

					_Response, err := svcs.AdminServices.ReadAdminsLength(p.Context, filterQuery, name)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"admin": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.AdminType)),
			Description: "Get an Administrator with id",
			Args: graphql.FieldConfigArgument{
				"adminId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					adminID := p.Args["adminId"].(string)

					_Response, err := svcs.AdminServices.ReadAdmin(p.Context, adminID)
					if err != nil {
						return nil, err
					}
					return transformations.DBAdminToGQLAdmin(_Response), nil
				},
			),
		},
	}
}

var adminMutation = func(svcs services.Services) map[string]*graphql.Field {
	return map[string]*graphql.Field{
		"createAdmin": {
			Type:        graphql.NewNonNull(schemas.AdminType),
			Description: "Create Admin in the system",
			Args: graphql.FieldConfigArgument{
				"fullname": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"role": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(schemas.EnumTypeAdminRole),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					fullname := p.Args["fullname"].(string)
					email := p.Args["email"].(string)
					role := p.Args["role"].(string)
					if adminData.Role != "Admin" {
						return nil, errors.New("PermissionDenied")
					}
					_Response, err := svcs.AdminServices.CreateAdmin(p.Context, fullname, email, role, &adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBAdminToGQLAdmin(_Response), nil
				},
			),
		},
		"loginAdmin": {
			Type:        graphql.NewNonNull(schemas.LoginAdminType),
			Description: "Login Admin",
			Args: graphql.FieldConfigArgument{
				"email": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				email := p.Args["email"].(string)
				password := p.Args["password"].(string)
				_Response, err := svcs.AdminServices.LoginAdmin(p.Context, email, password)
				if err != nil {
					return nil, err
				}
				return map[string]interface{}{
					"admin": transformations.DBAdminToGQLAdmin(&_Response.Admin),
					"token": _Response.Token,
				}, nil
			},
		},
		"updateAdminPassword": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Admin Password",
			Args: graphql.FieldConfigArgument{
				"oldPassword": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"password": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					oldPassword := p.Args["oldPassword"].(string)
					newPassword := p.Args["password"].(string)
					_Response, err := svcs.AdminServices.UpdateAdminPassword(p.Context, adminData.ID, oldPassword, newPassword)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"updateAdminDetails": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Admin Details",
			Args: graphql.FieldConfigArgument{
				"fullname": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"email": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"role": &graphql.ArgumentConfig{
					Type: schemas.EnumTypeAdminRole,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					takeFullname, fullNameOk := p.Args["fullname"].(string)
					takeEmail, emailOk := p.Args["email"].(string)
					takeRole, roleOk := p.Args["email"].(string)
					var fullname, email, role *string
					if fullNameOk {
						fullname = &takeFullname
					} else {
						fullname = nil
					}

					if emailOk {
						email = &takeEmail
					} else {
						email = nil
					}

					if roleOk {
						role = &takeRole
					} else {
						role = nil
					}
					_Response, err := svcs.AdminServices.UpdateAdmin(p.Context, adminData.ID, fullname, email, role)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"updateAdminPhone": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update Admin Phone",
			Args: graphql.FieldConfigArgument{
				"phone": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					phone := p.Args["phone"].(string)
					_Response, err := svcs.AdminServices.UpdateAdminPhone(p.Context, adminData.ID, phone)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deleteAdmin": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete Admin",
			Args: graphql.FieldConfigArgument{
				"adminId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					adminID := p.Args["adminId"].(string)
					if adminData.Role != "Admin" {
						return nil, errors.New("PermissionDenied")
					}
					_Response, err := svcs.AdminServices.DeleteAdmin(p.Context, adminID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"suspendAdmin": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Suspend Admin",
			Args: graphql.FieldConfigArgument{
				"adminId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"reason": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					adminID := p.Args["adminId"].(string)
					reason := p.Args["reason"].(string)
					if adminData.Role != "Admin" {
						return nil, errors.New("PermissionDenied")
					}
					_Response, err := svcs.AdminServices.SuspendAdmin(p.Context, adminID, adminData.ID, reason)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"restoreAdmin": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Restore Admin",
			Args: graphql.FieldConfigArgument{
				"adminId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					adminID := p.Args["adminId"].(string)
					if adminData.Role != "Admin" {
						return nil, errors.New("PermissionDenied")
					}
					_Response, err := svcs.AdminServices.RestoreAdmin(p.Context, adminID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeAdminResolver exposes the admin resolver
func ExposeAdminResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    adminQuery(services),
		Mutation: adminMutation(services),
	}
}
