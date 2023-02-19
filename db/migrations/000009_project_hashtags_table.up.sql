CREATE TABLE IF NOT EXISTS project_hashtags (
    id SERIAL PRIMARY KEY,
    project_id INT NOT NULL,
    hashtag_id INT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);