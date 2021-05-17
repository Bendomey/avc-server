package transformations

import (
	"github.com/Bendomey/avc-server/internal/orm/models"
	"github.com/gofrs/uuid"
)

// DBServiceToGQLUser transforms [user] db input to gql type
func DBServicingFieldToGQLUser(i *models.ServicingField) interface{} {
	return map[string]interface{}{
		"id":        i.ID.String(),
		"business":  DBBusinessToGQL(i.Business),
		"trademark": DBTrademarkToGQL(i.Trademark),
		"document":  DBDocumentToGQL(i.Document),
		"createdAt": i.CreatedAt,
		"updatedAt": i.UpdatedAt,
	}
}

func DBBusinessToGQL(i models.Business) interface{} {
	return map[string]interface{}{
		"country":        DBCountryToGQLUser(&i.Country),
		"entityType":     i.EntityType,
		"name":           i.Name,
		"owners":         i.Owners,
		"directors":      i.Directors,
		"address":        i.Address,
		"numberOfShares": i.NumberOfShares,
		"initialCapital": i.InitialCapital,
		"Industry":       i.Industry,
	}
}

func DBTrademarkToGQL(i models.Trademark) interface{} {
	return map[string]interface{}{
		"country":                   DBCountryToGQLUser(&i.Country),
		"ownershipType":             i.OwnershipType,
		"owners":                    i.Owners,
		"address":                   i.Address,
		"classificationOfTrademark": i.ClassificationOfTrademark,
		"uploads":                   i.Uploads,
	}

}

func DBDocumentToGQL(i models.Document) interface{} {
	return map[string]interface{}{
		"type":             i.Type,
		"natureOfDoc":      i.NatureOfDoc,
		"deadline":         i.Deadline,
		"existingDocument": i.ExistingDocuments,
		"newDocuments":     i.NewDocuments,
	}
}

//GQLServiceIDToDbID transforms id string from gql input and returns uuid from string
func GQLServicingFieldIDToDbID(i string) (*uuid.UUID, error) {
	updID, err := uuid.FromString(i)
	if err != nil {
		return nil, err
	}
	return &updID, nil
}
