ALTER TABLE user_projects
ADD CONSTRAINT fk_user_projects_project_id
FOREIGN KEY (project_id) REFERENCES projects(id)
ON DELETE CASCADE;