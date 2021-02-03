package utils

import (
	"time"
)

//FilterQuery type to help generate filter for queries
type FilterQuery struct {
	Limit     int            `json:"limit"`
	Skip      int            `json:"skip"`
	Order     string         `json:"order"`
	OrderBy   string         `json:"orderBy"`
	Search    *Search        `json:"search"`
	DateRange *DateRangeType `json:"dateRange"`
}

//DateRangeType typings
type DateRangeType struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

//Search typings
type Search struct {
	Criteria     string        `json:"criteria"`
	SearchFields []interface{} `json:"searchFields"`
}

// GenerateQuery takes a loook at what is coming from client and then generates a sieve
func GenerateQuery(argument map[string]interface{}) (*FilterQuery, error) {
	filterResult := FilterQuery{
		Limit:     -1,
		Skip:      -1,
		Order:     "asc",
		OrderBy:   "created_at",
		Search:    nil,
		DateRange: nil,
	}

	//pagination
	takePagination, paginationOk := argument["pagination"].(map[string]interface{})

	if paginationOk {
		limit, limitOk := takePagination["limit"].(int)
		if limitOk {
			filterResult.Limit = limit
		}

		skip, skipOk := takePagination["skip"].(int)
		if skipOk {
			filterResult.Skip = skip
		}
	}

	takeFilter, filterOk := argument["filter"].(map[string]interface{})

	if filterOk {
		//order
		order, orderOk := takeFilter["order"].(string)
		if orderOk {
			filterResult.Order = order
		}

		//orderBy
		orderBy, orderByOk := takeFilter["orderBy"].(string)
		if orderByOk {
			filterResult.OrderBy = orderBy
		}

		//dateRange
		dateRange, dateRangeOk := takeFilter["dateRange"].(map[string]interface{})

		if dateRangeOk {
			start := dateRange["start"].(time.Time)
			end := dateRange["end"].(time.Time)

			filterResult.DateRange = &DateRangeType{
				StartTime: start,
				EndTime:   end,
			}

		}

		//search
		search, searchOk := takeFilter["search"].(map[string]interface{})

		if searchOk {
			criteria := search["criteria"].(string)
			searchFields := search["searchFields"].([]interface{})
			filterResult.Search = &Search{
				Criteria:     criteria,
				SearchFields: searchFields,
			}

		}

	}

	return &filterResult, nil
}
