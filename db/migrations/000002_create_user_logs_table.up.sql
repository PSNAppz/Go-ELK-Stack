CREATE TABLE IF NOT EXISTS users_logs (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    operation VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);