CREATE TABLE projects (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    slug VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);