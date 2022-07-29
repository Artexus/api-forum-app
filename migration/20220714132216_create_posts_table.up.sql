CREATE TABLE
IF NOT EXISTS posts
(
	id SERIAL PRIMARY KEY,
	user_id INT REFERENCES users
(id),
	file_id INT references files
(id),
	title VARCHAR,
(255),
	description TEXT,
	created_at TIMESTAMP NOT NULL DEFAULT NOW
(),, 
	updated_at TIMESTAMP NOT NULL DEFAULT NOW
(), 
	deleted_at TIMESTAMP
);