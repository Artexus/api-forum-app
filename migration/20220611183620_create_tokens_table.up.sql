--CREATE TYPE IF NOT EXISTS token_type AS ENUM ('access_token', 'refresh_token');
CREATE TABLE IF NOT EXISTS tokens
(
	id SERIAL NOT NULL,
	user_id INT REFERENCES users(id),
	type token_type NOT NULL,
	token TEXT NOT NULL, 
	expired_at INT NOT NULL, 
	created_at TIMESTAMP, 
	updated_at TIMESTAMP, 
	deleted_at TIMESTAMP
);