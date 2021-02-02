package utils

import "github.com/graphql-go/graphql"

// GetReolvers helps you get reolvers from individual models
func GetReolvers(sample []map[string]*graphql.Field) map[string]*graphql.Field {
	var phg1 = map[string]*graphql.Field{}

	for _, val := range sample {
		for key, value := range val {
			phg1[key] = value
		}
	}

	return phg1
}
