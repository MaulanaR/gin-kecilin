package utils

import (
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type HelperUsecaseHandler struct {
	Filters           map[string][]string
	Sort              string
	AllowedSortFields map[string]bool
}

func (uc *HelperUsecaseHandler) SetSort() bson.D {
	if uc.Sort != "" {
		parts := strings.Split(uc.Sort, ",")
		sortDoc := bson.D{}
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			dir := int32(1)
			field := p
			if strings.HasPrefix(p, "-") {
				dir = -1
				field = strings.TrimPrefix(p, "-")
			}
			if uc.AllowedSortFields[field] {
				sortDoc = append(sortDoc, bson.E{Key: field, Value: dir})
			}
		}
		if len(sortDoc) > 0 {
			return sortDoc
		}
	}

	// default
	return bson.D{}
}

func (uc *HelperUsecaseHandler) SetFilter() bson.M {
	filter := bson.M{}

	for key, values := range uc.Filters {
		if len(values) == 0 {
			continue
		}
		val := values[0]

		if strings.Contains(key, "[$like]") {
			field := strings.Replace(key, "[$like]", "", 1)
			filter[field] = bson.M{"$regex": val, "$options": "i"}
		} else if strings.Contains(key, "[$eq]") {
			field := strings.Replace(key, "[$eq]", "", 1)
			filter[field] = val
		} else if strings.Contains(key, "[$in]") {
			field := strings.Replace(key, "[$in]", "", 1)
			filter[field] = bson.M{"$in": values}
		} else {
			filter[key] = val
		}
	}

	return filter
}

type Pagination struct {
	Page       int  `json:"page"`
	Limit      int  `json:"limit"`
	TotalCount int  `json:"total_count"`
	TotalPages int  `json:"total_pages"`
	HasNext    bool `json:"has_next"`
	HasPrev    bool `json:"has_prev"`
}

type Response struct {
	Status     string      `json:"status"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

func (r *Response) BuildResponse() map[string]interface{} {
	return map[string]interface{}{
		"status":  r.Status,
		"message": r.Message,
		"results": map[string]interface{}{
			"list":       r.Data,
			"pagination": r.Pagination,
		},
		"timestamp": time.Now(),
	}
}

func (r *Response) BuildSingleResponse() map[string]interface{} {
	return map[string]interface{}{
		"status":    r.Status,
		"message":   r.Message,
		"results":   r.Data,
		"timestamp": time.Now(),
	}
}
