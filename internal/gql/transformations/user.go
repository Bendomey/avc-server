package transformations

import "github.com/Bendomey/avc-server/internal/orm/models"

// DBUserToGQLUser transforms [user] db input to gql type
func DBUserToGQLUser(i *models.User) interface{} {
	return map[string]interface{}{
		"id":              i.ID.String(),
		"type":            i.Type,
		"lastName":        i.LastName,
		"firstName":       i.FirstName,
		"otherNames":      i.OtherNames,
		"email":           i.Email,
		"phone":           i.Phone,
		"emailVerifiedAt": i.EmailVerifiedAt,
		"phoneVerifiedAt": i.PhoneVerifiedAt,
		"createdAt":       i.CreatedAt,
		"updatedAt":       i.UpdatedAt,
	}
}
