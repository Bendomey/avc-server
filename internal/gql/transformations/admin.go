package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBAdminToGQLAdmin transforms [user] db input to gql type
func DBAdminToGQLAdmin(i *models.Admin) interface{} {
	return map[string]interface{}{
		"id":              i.ID.String(),
		"fullname":        i.FullName,
		"email":           i.Email,
		"phone":           i.Phone,
		"role":            i.Role,
		"phoneVerifiedAt": i.PhoneVerifiedAt,
		"suspendedAt":     i.SuspendedAt,
		"suspendedReason": i.SuspendedReason,
		"createdAt":       i.CreatedAt,
		"updatedAt":       i.UpdatedAt,
	}
}

//GQLInputIDToDbID transforms id string from gql input and returns uuid from string
func GQLInputIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
