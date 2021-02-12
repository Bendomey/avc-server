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
		"setupAt":         i.SetupAt,
		"createdAt":       i.CreatedAt,
		"updatedAt":       i.UpdatedAt,
	}
}

// DBUserToGQLLawyer transforms [user] db input to gql type
func DBUserToGQLLawyer(i *models.Lawyer) interface{} {
	if i == nil {
		return nil
	}
	return map[string]interface{}{
		"id":                      i.ID.String(),
		"digitalAddress":          i.DigitalAddress,
		"addressCountry":          i.AddressCountry,
		"user":                    DBUserToGQLUser(&i.User),
		"addressCity":             i.AddressCity,
		"addressStreetName":       i.AddressStreetName,
		"addressNumber":           i.AddressNumber,
		"firstYearOfBarAdmission": i.FirstYearOfBarAdmission,
		"licenseNumber":           i.LicenseNumber,
		"tin":                     i.TIN,
		"nationalIDFront":         i.NationalIDFront,
		"nationalIDBack":          i.NationalIDBack,
		"bARMembershipCard":       i.BARMembershipCard,
		"LawCertificate":          i.LawCertificate,
		"CV":                      i.CV,
		"coverLetter":             i.CoverLetter,
		"createdAt":               i.CreatedAt,
		"updatedAt":               i.UpdatedAt,
	}
}

// DBUserToGQLCustomer transforms [user] db input to gql type
func DBUserToGQLCustomer(i *models.Customer) interface{} {
	if i == nil {
		return nil
	}

	return map[string]interface{}{
		"id":                           i.ID.String(),
		"type":                         i.Type,
		"digitalAddress":               i.DigitalAddress,
		"addressCountry":               i.AddressCountry,
		"addressCity":                  i.AddressCity,
		"addressStreetName":            i.AddressStreetName,
		"addressNumber":                i.AddressNumber,
		"companyName":                  i.CompanyName,
		"companyEntityType":            i.CompanyEntityType,
		"tin":                          i.TIN,
		"companyEntityTypeOther":       i.CompanyEntityTypeOther,
		"companyCountryOfRegistration": i.CompanyCountryOfRegistration,
		"companyDateOfRegistration":    i.CompanyDateOfRegistration,
		"companyRegistrationNumber":    i.CompanyRegistrationNumber,
		"user":                         DBUserToGQLUser(&i.User),
		"createdAt":                    i.CreatedAt,
		"updatedAt":                    i.UpdatedAt,
	}
}
