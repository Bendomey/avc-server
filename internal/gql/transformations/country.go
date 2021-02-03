package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBCountryToGQLUser transforms [user] db input to gql type
func DBCountryToGQLUser(i *models.Country) interface{} {
	return map[string]interface{}{
		"id":          i.ID.String(),
		"name":        i.Name,
		"description": i.Description,
		"currency":    i.Currency,
		"image":       i.Image,
		"createdBy":   DBUserToGQLUser(&i.CreatedBy),
		"createdAt":   i.CreatedAt,
		"updatedAt":   i.UpdatedAt,
	}
}

//GQLCountryIDToDbID transforms id string from gql input and returns uuid from string
func GQLCountryIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
