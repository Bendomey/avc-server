package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBBlogPostToGQLUser transforms [blog post] db input to gql type
func DBBlogPostToGQLUser(i *models.BlogPost) interface{} {
	return map[string]interface{}{
		"id":        i.ID.String(),
		"title":     i.Title,
		"image":     i.Image,
		"status":    i.Status,
		"tag":       DBTagToGQLUser(&i.Tag),
		"details":   i.Details,
		"createdBy": DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
}

//GQLBlogPostIDToDbID transforms id string from gql input and returns uuid from string
func GQLBlogPostIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
