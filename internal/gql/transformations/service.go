package transformations

import (
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBServiceToGQLUser(i *models.Service) interface{} {
	return map[string]interface{}{
		"id":          i.ID.String(),
		"name":        i.Name,
		"price":       i.Price,
		"description": i.Description,
		"type":        fmt.Sprintf("%s", i.Type),
		"createdBy":   DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt":   i.CreatedAt,
		"updatedAt":   i.UpdatedAt,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLServiceIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
