ALTER TABLE project_hashtags
ADD CONSTRAINT fk_hashtag_id
FOREIGN KEY (hashtag_id)
REFERENCES hashtags(id)
ON DELETE CASCADE;
