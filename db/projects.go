package db

import (
	"database/sql"

	"github.com/PSNAppz/Fold-ELK/models"
)

func (db Database) CreateProject(project *models.Project, user *models.User, hashtags *[]models.Hashtag) error {
	var projectId int
	query := `INSERT INTO projects(name, slug, description) VALUES ($1, $2, $3) RETURNING id`
	err := db.Conn.QueryRow(query, project.Name, project.Slug, project.Description).Scan(&projectId)
	if err != nil {
		return err
	}

	// Insert project-user association
	if user != nil {
		userProjectQuery := `INSERT INTO user_projects(project_id, user_id) VALUES ($1, $2)`
		_, err = db.Conn.Exec(userProjectQuery, projectId, user.ID)
		if err != nil {
			return err
		}
	}

	// Insert project-hashtag association
	if hashtags != nil {
		for _, hashtag := range *hashtags {
			projectHashtagQuery := `INSERT INTO project_hashtags(project_id, hashtag_id) VALUES ($1, $2)`
			_, err = db.Conn.Exec(projectHashtagQuery, projectId, hashtag.ID)
			if err != nil {
				return err
			}
		}
	}

	// log the operation for logstash to pick up and send to elasticsearch
	// Here we are doing this at app level.
	logQuery := `INSERT INTO project_logs(project_id, operation) VALUES ($1, $2)`
	project.ID = projectId
	_, err = db.Conn.Exec(logQuery, project.ID, insertOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) UpdateProject(projectId int, project models.Project) error {
	query := "UPDATE projects SET name=$1, slug=$2, description=$3 WHERE id=$4"
	_, err := db.Conn.Exec(query, project.Name, project.Slug, project.Description, projectId)
	if err != nil {
		return err
	}

	project.ID = projectId
	logQuery := "INSERT INTO project_logs(project_id, operation) VALUES ($1, $2, $3)"
	_, err = db.Conn.Exec(logQuery, project.ID, updateOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) DeleteProject(projectId int) error {
	query := "DELETE FROM projects WHERE id=$1"
	_, err := db.Conn.Exec(query, projectId)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoRecord
		}
		return err
	}

	logQuery := "INSERT INTO project_logs(project_id, operation) VALUES ($1, $2)"
	_, err = db.Conn.Exec(logQuery, projectId, deleteOp)
	if err != nil {
		db.Logger.Err(err).Msg("could not log operation for logstash")
	}
	return nil
}

func (db Database) GetProjectById(projectId int) (models.Project, error) {
	project := models.Project{}
	query := "SELECT id, name, slug, description FROM projects WHERE id = $1"
	row := db.Conn.QueryRow(query, projectId)
	switch err := row.Scan(&project.ID, &project.Name, &project.Slug, &project.Description); err {
	case sql.ErrNoRows:
		return project, ErrNoRecord
	default:
		return project, err
	}
}

// Get all projects by a user
func (db Database) GetProjectsByUserId(userId int) ([]models.Project, error) {
	var list []models.Project
	query := "SELECT id, name, slug, description FROM projects WHERE id IN (SELECT project_id FROM user_projects WHERE user_id = $1)"
	rows, err := db.Conn.Query(query, userId)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Slug, &project.Description)
		if err != nil {
			return list, err
		}
		list = append(list, project)
	}
	return list, nil
}

func (db Database) GetProjects() ([]models.Project, error) {
	var list []models.Project
	query := "SELECT id, name, slug, description FROM projects ORDER BY id DESC"
	rows, err := db.Conn.Query(query)
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var project models.Project
		err := rows.Scan(&project.ID, &project.Name, &project.Slug, &project.Description)
		if err != nil {
			return list, err
		}
		list = append(list, project)
	}
	return list, nil
}
