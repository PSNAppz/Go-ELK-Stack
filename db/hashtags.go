package db

import (
	"database/sql"

	"github.com/PSNAppz/Fold-ELK/models"
)

func (db Database) CreateHashtag(hashtag *models.Hashtag) error {
	var id int
	query := `INSERT INTO hashtags(name) VALUES ($1) RETURNING id`
	err := db.Conn.QueryRow(query, hashtag.Name).Scan(&id)
	if err != nil {
		return err
	}

	// log the operation for logstash to pick up and send to elasticsearch
	// Here we are doing this at app level.
	logQuery := `INSERT INTO hashtags_logs(hashtag_id, operation) VALUES ($1, $2)`
	hashtag.ID = id
	_, err = db.Conn.Exec(logQuery, hashtag.ID, insertOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) UpdateHashtag(hashtagId int, hashtag models.Hashtag) error {
	query := "UPDATE hashtags SET name=$1 WHERE id=$2"
	_, err := db.Conn.Exec(query, hashtag.Name, hashtagId)
	if err != nil {
		return err
	}

	hashtag.ID = hashtagId
	logQuery := "INSERT INTO hashtags_logs(hashtag_id, operation) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(logQuery, hashtag.ID, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeleteHashtag(hashtagId int) error {
	query := "DELETE FROM hashtags WHERE id=$1"
	_, err := db.Conn.Exec(query, hashtagId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO hashtags_logs(hashtag_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, hashtagId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) GetHashtagById(hashtagId int) (models.Hashtag, error) {
	hashtag := models.Hashtag{}
	query := "SELECT id, name FROM hashtags WHERE id = $1"
	row := db.Conn.QueryRow(query, hashtagId)
	switch err := row.Scan(&hashtag.ID, &hashtag.Name); err {
	case sql.ErrNoRows:
		return hashtag, ErrNoRecord
	default:
		return hashtag, err
	}
}

func (db Database) GetHashtags() ([]models.Hashtag, error) {
	var list []models.Hashtag
	query := "SELECT id, name FROM hashtags ORDER BY id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var hashtag models.Hashtag
		err := rows.Scan(&hashtag.ID, &hashtag.Name)
		if err != nil {
			return list, err
		}
		list = append(list, hashtag)
	}
	return list, nil
}
