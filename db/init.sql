-- clear tables
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS forum CASCADE;

-- install module with case-insensitive string
CREATE EXTENSION IF NOT EXISTS citext;

-- tables
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL       UNIQUE,
    nickname    CITEXT       COLLATE ucs_basic NOT NULL PRIMARY KEY,
    fullname    VARCHAR(255) NOT NULL,
    email       CITEXT       NOT NULL UNIQUE,
    about       TEXT         NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS forum (
    id          SERIAL       UNIQUE,
    slug        CITEXT       NOT NULL PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    "user"      CITEXT       NOT NULL REFERENCES users (nickname),
    posts       INTEGER      NOT NULL DEFAULT 0,
    threads     INTEGER      NOT NULL DEFAULT 0
);
