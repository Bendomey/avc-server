package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBFaqToGQLUser transforms [faq] db input to gql type
func DBFaqToGQLUser(i *models.Faq) interface{} {
	return map[string]interface{}{
		"id":        i.ID.String(),
		"question":  i.Question,
		"answer":    i.Answer,
		"createdBy": DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
}

//GQLFaqIDToDbID transforms id string from gql input and returns uuid from string
func GQLFaqIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
