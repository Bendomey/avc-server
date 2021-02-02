package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/Bendomey/goutilities/pkg/validatetoken"
	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

//AdminFromToken unmarshals cliams from jwt to get admin id
type AdminFromToken struct {
	ID string `json:"id"`
}

func extractAdminToken(unattendedToken string) (string, error) {
	//remove bearer
	strArr := strings.Split(unattendedToken, " ")
	if len(strArr) != 2 {
		return "", errors.New("AuthorizationFailed")
	}
	return strArr[1], nil
}

func validateAdmin(unattendedToken string) (*AdminFromToken, error) {
	//extract token
	token, extractTokenErr := extractAdminToken(unattendedToken)
	if extractTokenErr != nil {
		return nil, extractTokenErr
	}

	//extract token metadata
	rawToken, validateError := validatetoken.ValidateJWTToken(token, os.Getenv("ADMIN_SECRET"))
	if validateError != nil {
		return nil, errors.New("AuthorizationFailed")
	}

	claims, ok := rawToken.Claims.(jwt.MapClaims)
	var adminFromTokenImplementation AdminFromToken
	if ok && rawToken.Valid {
		adminFromTokenImplementation.ID = claims["id"].(string)
	}

	//check if its exists in db
	// _, err := userService.GetUser(ctx, userFromTokenImplementation.ID)
	// if err != nil {
	// 	return nil, err
	// }
	return &adminFromTokenImplementation, nil
}

type manipulateAnythingFromAdmin func(params graphql.ResolveParams, adminData *AdminFromToken) (interface{}, error)

// AuthenticateAdmin checks if the user trying to access that resource is truly an Admin
func AuthenticateAdmin(fn manipulateAnythingFromAdmin) func(params graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		token, tokenExtractionErr := GetContextInjected(params.Context)
		if tokenExtractionErr != nil {
			return nil, tokenExtractionErr
		}
		validated, validateError := validateAdmin(token)
		if validateError != nil {
			return nil, validateError
		}
		return fn(params, validated)
	}
}
