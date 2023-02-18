CREATE TABLE user_projects (
    project_id INT,
    user_id INT,
    FOREIGN KEY (project_id) REFERENCES projects(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);