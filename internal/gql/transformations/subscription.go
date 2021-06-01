package transformations

import (
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBSubscriptionToGQLUser(i *models.Subscription) interface{} {
	return map[string]interface{}{
		"id":           i.ID.String(),
		"package":      DBPackageToGQLUser(&i.Package),
		"status":       fmt.Sprintf("%s", i.Status),
		"createdBy":    DBUserToGQLUser(&i.CreatedBy),
		"subscribedAt": i.SubscribeAt,
		"expiresAt":    i.ExpiresAt,
		"createdAt":    i.CreatedAt,
		"updatedAt":    i.UpdatedAt,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLSubscriptionIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
