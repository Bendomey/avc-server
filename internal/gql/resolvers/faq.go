package resolvers

import (
	"github.com/Bendomey/avc-server/internal/gql/schemas"
	"github.com/Bendomey/avc-server/internal/gql/transformations"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/graphql-go/graphql"
)

var faqsQuery = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"faqs": {
			Type:        graphql.NewNonNull(graphql.NewList(schemas.FaqType)),
			Description: "Get faqs list",
			Args: graphql.FieldConfigArgument{
				"pagination": &graphql.ArgumentConfig{
					Type: schemas.PaginationType,
				},
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterFaqsType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				argument := p.Args
				filterQuery, filterErr := utils.GenerateQuery(argument)
				if filterErr != nil {
					return nil, filterErr
				}

				//fields
				takeFilter, filterOk := argument["filter"].(map[string]interface{})
				var question, answer *string

				if filterOk {
					takeQuestion, questionOk := takeFilter["question"].(string)
					takeAnswer, answerOk := takeFilter["question"].(string)
					if questionOk {
						question = &takeQuestion
					}

					if answerOk {
						answer = &takeAnswer
					}
				}

				_Response, err := svcs.FaqServices.ReadFAQs(p.Context, filterQuery, question, answer)
				if err != nil {
					return nil, err
				}
				var newResponse []interface{}
				//loop to get the types
				for _, dbRec := range _Response {
					rec := transformations.DBFaqToGQLUser(dbRec)
					newResponse = append(newResponse, rec)
				}
				return newResponse, nil
			},
		},
		"faqsLength": {
			Type:        graphql.NewNonNull(graphql.Int),
			Description: "Get faqs Counts",
			Args: graphql.FieldConfigArgument{
				"filter": &graphql.ArgumentConfig{
					Type: schemas.FilterFaqsType,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				argument := p.Args
				filterQuery, filterErr := utils.GenerateQuery(argument)
				if filterErr != nil {
					return nil, filterErr
				}
				//fields
				takeFilter, filterOk := argument["filter"].(map[string]interface{})
				var question, answer *string

				if filterOk {
					takeQuestion, questionOk := takeFilter["question"].(string)
					if questionOk {
						question = &takeQuestion
					} else {
						question = nil
					}

					takeAnswer, answerOk := takeFilter["answer"].(string)
					if answerOk {
						answer = &takeAnswer
					} else {
						answer = nil
					}
				}

				_Response, err := svcs.FaqServices.ReadFAQsLength(p.Context, filterQuery, question, answer)
				if err != nil {
					return nil, err
				}
				return _Response, nil
			},
		},
		"faq": {
			Type:        schemas.FaqType,
			Description: "Get single faq",
			Args: graphql.FieldConfigArgument{
				"faqId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				faqID := p.Args["faqId"].(string)

				_Response, err := svcs.FaqServices.ReadFAQ(p.Context, faqID)
				if err != nil {
					return nil, err
				}
				return transformations.DBFaqToGQLUser(_Response), nil
			},
		},
	}
}

var faqsMutation = func(svcs services.Services) map[string]*graphql.Field {

	return map[string]*graphql.Field{
		"createFaq": {
			Type:        graphql.NewNonNull(schemas.FaqType),
			Description: "Create faq",
			Args: graphql.FieldConfigArgument{
				"question": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"answer": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					question := p.Args["question"].(string)
					answer := p.Args["answer"].(string)

					_Response, err := svcs.FaqServices.CreateFAQ(p.Context, question, answer, adminData.ID)
					if err != nil {
						return nil, err
					}
					return transformations.DBFaqToGQLUser(_Response), nil
				},
			),
		},
		"updateFaq": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Update faq",
			Args: graphql.FieldConfigArgument{
				"faqId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
				"question": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"answer": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					faqID := p.Args["faqId"].(string)
					takeQuestion, questionOk := p.Args["question"].(string)
					takeAnswer, answerOk := p.Args["answer"].(string)

					var question, answer *string

					//validations
					if questionOk {
						question = &takeQuestion
					} else {
						question = nil
					}

					if answerOk {
						answer = &takeAnswer
					} else {
						answer = nil
					}

					_Response, err := svcs.FaqServices.UpdateFAQ(p.Context, faqID, question, answer)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
		"deleteFaq": {
			Type:        graphql.NewNonNull(graphql.Boolean),
			Description: "Delete faq",
			Args: graphql.FieldConfigArgument{
				"faqId": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.ID),
				},
			},
			Resolve: utils.AuthenticateAdmin(
				func(p graphql.ResolveParams, adminData *utils.AdminFromToken) (interface{}, error) {
					faqID := p.Args["faqId"].(string)

					_Response, err := svcs.FaqServices.DeleteFAQ(p.Context, faqID)
					if err != nil {
						return nil, err
					}
					return _Response, nil
				},
			),
		},
	}
}

// ExposeFaqResolver exposes the faqs Reesolver
func ExposeFaqResolver(services services.Services) ResolverLoader {
	return ResolverLoader{
		Query:    faqsQuery(services),
		Mutation: faqsMutation(services),
	}
}
