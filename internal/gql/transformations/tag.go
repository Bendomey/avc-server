package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBTagToGQLUser transforms [user] db input to gql type
func DBTagToGQLUser(i *models.Tag) interface{} {
	return map[string]interface{}{
		"id":        i.ID.String(),
		"name":      i.Name,
		"createdBy": DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
}

//GQLTagIDToDbID transforms id string from gql input and returns uuid from string
func GQLTagIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
