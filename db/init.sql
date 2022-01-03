-- install module with case-insensitive string
CREATE EXTENSION IF NOT EXISTS citext;

-- tables
CREATE TABLE IF NOT EXISTS users (
    nickname CITEXT COLLATE ucs_basic NOT NULL PRIMARY KEY,
    fullname VARCHAR(255) NOT NULL,
    email CITEXT NOT NULL UNIQUE,
    about TEXT NOT NULL DEFAULT ''
);
