package transformations

import (
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBPaymentToGQLUser(i *models.Payment) interface{} {
	return map[string]interface{}{
		"id":               i.ID.String(),
		"code":             i.Code.String(),
		"amount":           i.Amount,
		"servicing":        DBServicingToGQLUser(i.Servicing),
		"subscription":     DBSubscriptionToGQLUser(i.Subscription),
		"status":           fmt.Sprintf("%s", i.Status),
		"createdBy":        DBUserToGQLUser(&i.CreatedBy),
		"authorizationUrl": i.AuthorizationUrl,
		"accessCode":       i.AccessCode,
		"createdAt":        i.CreatedAt,
		"updatedAt":        i.UpdatedAt,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLPaymentIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
