CREATE TABLE
IF NOT EXISTS files
(
	id SERIAL PRIMARY KEY,
	path varchar
(255),
	created_at TIMESTAMP, 
	updated_at TIMESTAMP, 
	deleted_at TIMESTAMP
);