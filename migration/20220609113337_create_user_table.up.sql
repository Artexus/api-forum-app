CREATE OR REPLACE FUNCTION trigger_set_updated() RETURNS trigger
AS
$$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$
LANGUAGE 'plpgsql';

CREATE TABLE IF NOT EXISTS users
(
    id SERIAL PRIMARY KEY,
    username VARCHAR (64) NOT NULL,
    email VARCHAR (64) NOT NULL,
    password VARCHAR(64) NOT NULL, 
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,
    created_at TIMESTAMP
);

CREATE TRIGGER update_users_updated_at 
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_updated();

    