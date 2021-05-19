package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
)

type GetMeResponse struct {
	Success      bool                    `json:"success"`
	Data         *map[string]interface{} `json:"data"`
	ErrorMessage *string                 `json:"errorMessage"`
}

var errorResponseIfTokenNotFound string

func init() {
	errorResponseIfTokenNotFound = "Please add your token to the header"

}

// GetMe is simple keep-alive/ping handler
func GetMe(services services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		userId := r.Header.Get("Authorization")

		if userId == "" {
			//if header is not passed
			json.NewEncoder(w).Encode(GetMeResponse{
				Success:      false,
				ErrorMessage: &errorResponseIfTokenNotFound,
			})
		} else {

			// if header is passed
			_Response, err := services.UserServices.GetMe(context.TODO(), userId)
			if err != nil {
				errorMessage := err.Error()
				json.NewEncoder(w).Encode(GetMeResponse{
					Success:      false,
					ErrorMessage: &errorMessage,
				})

			} else {
				sendThis := map[string]interface{}{
					"user":     transformations.DBUserToGQLUser(&_Response.User),
					"lawyer":   transformations.DBUserToGQLLawyer(_Response.Lawyer),
					"customer": transformations.DBUserToGQLCustomer(_Response.Customer),
				}
				json.NewEncoder(w).Encode(GetMeResponse{
					Success: true,
					Data:    &sendThis,
				})
			}

		}

	}
}
