ALTER TABLE refresh_tokens
ADD CONSTRAINT refresh_tokens_user_id_key UNIQUE (user_id);