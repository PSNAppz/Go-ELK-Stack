package models

import "time"

type Project struct {
	ID          int       `db:"id"`
	Name        string    `db:"name"`
	Slug        string    `db:"slug"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
}

type ProjectHashtag struct {
	HashtagID int `db:"hashtag_id"`
	ProjectID int `db:"project_id"`
}

type CreateProjectRequest struct {
	Name        string    `json:"name" binding:"required"`
	Slug        string    `json:"slug" binding:"required"`
	Description string    `json:"description" binding:"required"`
	UserID      int       `json:"user_id" binding:"required"`
	Hashtags    []Hashtag `json:"hashtags"`
}
