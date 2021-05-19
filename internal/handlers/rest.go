package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Bendomey/avc-server/internal/services"
)

type GetMeResponse struct {
	success      bool
	data         *interface{}
	errorMessage *string
}

var errorResponseIfTokenNotFound string

func init() {
	errorResponseIfTokenNotFound = "Please add your token to the header"

}

// GetMe is simple keep-alive/ping handler
func GetMe(services services.Services) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		setupResponse(&w, r)
		fmt.Print(r.Header.Get("token"))
		userId := r.Header.Get("token")

		if userId == "" {
			//if header is not passed
			json.NewEncoder(w).Encode(GetMeResponse{
				success:      false,
				errorMessage: &errorResponseIfTokenNotFound,
			})
		} else {

			// if header is passed
			res, _ := services.UserServices.GetMe(context.TODO(), userId)
			json.NewEncoder(w).Encode(res)
		}

	}
}
