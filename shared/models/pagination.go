package models

type Page[T any] struct {
	Page  int64 `json:"page" bson:"page"`
	Size  int64 `json:"size" bson:"size"`
	Total int64 `json:"total" bson:"total"`
	Count int   `json:"count" bson:"count"`
	Items []T   `json:"items" bson:"items"`
}
