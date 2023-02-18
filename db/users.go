package db

import (
	"database/sql"

	"github.com/PSNAppz/Fold-ELK/models"
)

func (db Database) CreateUser(user *models.User) error {
	var id int
	query := `INSERT INTO users(name) VALUES ($1) RETURNING id`
	err := db.Conn.QueryRow(query, user.Name).Scan(&id)
	if err != nil {
		return err
	}

	// log the operation for logstash to pick up and send to elasticsearch
	// Here we are doing this at app level.
	logQuery := `INSERT INTO user_logs(user_id, operation) VALUES ($1, $2)`
	user.ID = id
	_, err = db.Conn.Exec(logQuery, user.ID, insertOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) UpdateUser(userId int, user models.User) error {
	query := "UPDATE users SET name=$1 WHERE id=$2"
	_, err := db.Conn.Exec(query, user.Name, userId)
	if err != nil {
		return err
	}

	user.ID = userId
	logQuery := "INSERT INTO user_logs(user_id, operation) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(logQuery, user.ID, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeleteUser(userId int) error {
	query := "DELETE FROM users WHERE id=$1"
	_, err := db.Conn.Exec(query, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO user_logs(user_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, userId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) GetUserById(userId int) (models.User, error) {
	user := models.User{}
	query := "SELECT id, name FROM users WHERE id = $1"
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.ID, &user.Name); err {
	case sql.ErrNoRows:
		return user, ErrNoRecord
	default:
		return user, err
	}
}

func (db Database) GetUsers() ([]models.User, error) {
	var list []models.User
	query := "SELECT id, name FROM users ORDER BY id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name)
		if err != nil {
			return list, err
		}
		list = append(list, user)
	}
	return list, nil
}
