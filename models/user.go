package models

import "time"

type User struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
}

type UserProject struct {
	ProjectID int `db:"project_id"`
	UserID    int `db:"user_id"`
}
