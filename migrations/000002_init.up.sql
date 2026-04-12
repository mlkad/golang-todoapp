CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 1,
  full_name VARCHAR(100) NOT NULL CHECK (char_length(full_name) BETWEEN 3 AND 100),
  phone_number VARCHAR(15) CHECK (
    phone_number ~ '^\+[0-9]+$'
    AND
    char_length(phone_number) BETWEEN 10 AND 15
  )
);

CREATE TABLE todoapp.tasks {
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 1,
  title VARCHAR(100) NOT NULL,
  description VARCHAR(1000),
  completed BOOLEAN NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  completed_at TIMESTAMPTZ,

  author_user_id INTEGER NOT NULL REFERENCES todoapp.users(id)
}