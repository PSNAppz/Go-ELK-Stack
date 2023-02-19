CREATE TABLE IF NOT EXISTS hashtags_logs (
  id SERIAL PRIMARY KEY,
  hashtag_id INT NOT NULL,
  operation VARCHAR(20) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW()
);