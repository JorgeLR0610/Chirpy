-- +goose Up
ALTER TABLE refresh_tokens 
DROP CONSTRAINT refresh_tokens_user_id_key;

-- +goose Down
ALTER TABLE refresh_tokens 
ADD CONSTRAINT user_id_unique 
UNIQUE (user_id);