package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBPackageToGQLUser transforms [user] db input to gql type
func DBPackageToGQLUser(i *models.Package) interface{} {
	return map[string]interface{}{
		"id":             i.ID.String(),
		"name":           i.Name,
		"amountPerMonth": i.AmountPerMonth,
		"amountPerYear":  i.AmountPerYear,
		"description":    i.Description,
		"status":         i.Status,
		"createdBy":      DBAdminToGQLAdmin(&i.CreatedBy),
		"requestedBy":    DBUserToGQLUser(&i.RequestedBy),
		"createdAt":      i.CreatedAt,
		"updatedAt":      i.UpdatedAt,
	}
}

//GQLPackageIDToDbID transforms id string from gql input and returns uuid from string
func GQLPackageIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
