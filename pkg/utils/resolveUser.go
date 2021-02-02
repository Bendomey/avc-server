package utils

import (
	"errors"
	"os"
	"strings"

	"github.com/Bendomey/goutilities/pkg/validatetoken"
	"github.com/dgrijalva/jwt-go"
	"github.com/graphql-go/graphql"
)

//UserFromToken unmarshals cliams from jwt to get admin id
type UserFromToken struct {
	ID string `json:"id"`
}

func extractUserToken(unattendedToken string) (string, error) {
	//remove bearer
	strArr := strings.Split(unattendedToken, " ")
	if len(strArr) != 2 {
		return "", errors.New("AuthorizationFailed")
	}
	return strArr[1], nil
}

func validateUser(unattendedToken string) (*UserFromToken, error) {
	//extract token
	token, extractTokenErr := extractUserToken(unattendedToken)
	if extractTokenErr != nil {
		return nil, extractTokenErr
	}

	//extract token metadata
	rawToken, validateError := validatetoken.ValidateJWTToken(token, os.Getenv("ADMIN_SECRET"))
	if validateError != nil {
		return nil, errors.New("AuthorizationFailed")
	}

	claims, ok := rawToken.Claims.(jwt.MapClaims)
	var userFromTokenImplementation UserFromToken
	if ok && rawToken.Valid {
		userFromTokenImplementation.ID = claims["id"].(string)
	}

	//check if its exists in db
	// _, err := userService.GetUser(ctx, userFromTokenImplementation.ID)
	// if err != nil {
	// 	return nil, err
	// }
	return &userFromTokenImplementation, nil
}

type manipulateAnythingFromUser func(params graphql.ResolveParams, adminData *UserFromToken) (interface{}, error)

// AuthenticateUser checks if the user trying to access that resource is truly a user (customer, lawyer)
func AuthenticateUser(fn manipulateAnythingFromUser) func(params graphql.ResolveParams) (interface{}, error) {
	return func(params graphql.ResolveParams) (interface{}, error) {
		token, tokenExtractionErr := GetContextInjected(params.Context)
		if tokenExtractionErr != nil {
			return nil, tokenExtractionErr
		}
		validated, validateError := validateUser(token)
		if validateError != nil {
			return nil, validateError
		}
		return fn(params, validated)
	}
}
