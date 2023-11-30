package models

import "time"

type Event struct {
	ID    string `json:"id" bson:"_id"`
	Type  string `json:"type" bson:"type"`
	Actor struct {
		ID           int    `json:"id" bson:"id"`
		Login        string `json:"login" bson:"login"`
		DisplayLogin string `json:"display_login" bson:"display_login"`
		GravatarID   string `json:"gravatar_id" bson:"gravatar_id"`
		Url          string `json:"url" bson:"url"`
		AvatarUrl    string `json:"avatar_url" bson:"avatar_url"`
	} `json:"actor" bson:"actor"`
	Repo struct {
		ID    int    `json:"id" bson:"id"`
		Name  string `json:"name" bson:"name"`
		Url   string `json:"url" bson:"url"`
		Stars int    `json:"stars" bson:"stars"`
	} `json:"repo" bson:"repo"`
	Payload   map[string]any `json:"payload" bson:"payload"`
	Public    bool           `json:"public" bson:"public"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
}

type Actor struct {
	ID           int    `json:"id" bson:"id"`
	Login        string `json:"login" bson:"login"`
	DisplayLogin string `json:"display_login" bson:"display_login"`
	GravatarID   string `json:"gravatar_id" bson:"gravatar_id"`
	Url          string `json:"url" bson:"url"`
	AvatarUrl    string `json:"avatar_url" bson:"avatar_url"`
}

type ApiRepo struct {
	StargazersCount int `json:"stargazers_count" bson:"stargazers_count"`
}

type EventRepo struct {
	ID    int    `json:"id" bson:"id"`
	Name  string `json:"name" bson:"name"`
	Url   string `json:"url" bson:"url"`
	Stars int    `json:"stars" bson:"stars"`
}
