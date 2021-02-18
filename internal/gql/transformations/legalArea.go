package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBLegalAreaToGQLUser transforms [user] db input to gql type
func DBLegalAreaToGQLUser(i *models.LegalArea) interface{} {
	return map[string]interface{}{
		"id":          i.ID.String(),
		"name":        i.Name,
		"description": i.Description,
		"image":       i.Image,
		"createdBy":   DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt":   i.CreatedAt,
		"updatedAt":   i.UpdatedAt,
	}
}

//GQLLegalAreaIDToDbID transforms id string from gql input and returns uuid from string
func GQLLegalAreaIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
