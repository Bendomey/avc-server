package transformations

import (
	"fmt"

	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBPackageServiceToGQLUser(i *models.PackageService) interface{} {
	return map[string]interface{}{
		"id":        i.ID.String(),
		"service":   DBServiceToGQLUser(&i.Service),
		"package":   DBPackageToGQLUser(&i.Package),
		"type":      fmt.Sprintf("%s", i.Type),
		"quantity":  i.Quantity,
		"isActive":  i.IsActive,
		"createdBy": DBAdminToGQLAdmin(&i.CreatedBy),
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLPackageServiceIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
