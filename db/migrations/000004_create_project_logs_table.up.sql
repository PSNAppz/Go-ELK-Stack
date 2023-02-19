CREATE TABLE IF NOT EXISTS project_logs (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
