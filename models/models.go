package models

import "time"

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type Hashtag struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

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

type UserProject struct {
	ProjectID int `db:"project_id"`
	UserID    int `db:"user_id"`
}
