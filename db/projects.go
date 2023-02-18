package db

import (
	"database/sql"

	"github.com/PSNAppz/Fold-ELK/models"
)

func (db Database) CreateProject(project *models.Project) error {
	var id int
	query := `INSERT INTO projects(name, slug, description) VALUES ($1, $2, $3) RETURNING id`
	err := db.Conn.QueryRow(query, project.Name, project.Slug, project.Description).Scan(&id)
	if err != nil {
		return err
	}

	// log the operation for logstash to pick up and send to elasticsearch
	// Here we are doing this at app level.
	logQuery := `INSERT INTO project_logs(project_id, operation) VALUES ($1, $2)`
	project.ID = id
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
	query := "SELECT id, name, slug, descriotion FROM projects WHERE id = $1"
	row := db.Conn.QueryRow(query, projectId)
	switch err := row.Scan(&project.ID, &project.Name, &project.Slug, &project.Description); err {
	case sql.ErrNoRows:
		return project, ErrNoRecord
	default:
		return project, err
	}
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
