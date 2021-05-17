package transformations

import (
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBServicingToGQLUser(i *models.Servicing) interface{} {
	return map[string]interface{}{
		"id":      i.ID.String(),
		"service": DBServiceToGQLUser(&i.Service),
		"cost":    i.Cost,
		// "payment":       DBPaymentToGQLUser(&i.Payment),
		"status":        fmt.Sprintf("%s", i.Status),
		"subscription":  DBSubscriptionToGQLUser(&i.Subscription),
		"lawyer":        DBUserToGQLUser(&i.Lawyer),
		"serviceFields": DBServicingFieldToGQLUser(&i.ServiceFields),
		"createdBy":     DBUserToGQLUser(&i.CreatedBy),
		"createdAt":     i.CreatedAt,
		"updatedAt":     i.UpdatedAt,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLServicingIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
