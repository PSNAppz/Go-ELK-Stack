CREATE TABLE IF NOT EXISTS user_project_logs (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    user_id INT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);