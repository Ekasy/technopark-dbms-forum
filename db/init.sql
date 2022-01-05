-- clear tables
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS forum CASCADE;
-- DROP TABLE IF EXISTS threads CASCADE;

-- install module with case-insensitive string
CREATE EXTENSION IF NOT EXISTS citext;

-- tables
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL       UNIQUE,
    nickname    CITEXT       COLLATE ucs_basic NOT NULL PRIMARY KEY,
    fullname    TEXT         NOT NULL,
    email       CITEXT       NOT NULL UNIQUE,
    about       TEXT         NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS forum (
    id          SERIAL       UNIQUE,
    slug        CITEXT       NOT NULL PRIMARY KEY,
    title       TEXT         NOT NULL,
    author      CITEXT       NOT NULL,
    posts       INTEGER      NOT NULL DEFAULT 0,
    threads     INTEGER      NOT NULL DEFAULT 0,
	FOREIGN KEY (author) REFERENCES users (nickname)
);

CREATE TABLE IF NOT EXISTS threads (
    id          SERIAL          NOT NULL PRIMARY KEY,
    slug        CITEXT          DEFAULT '',
    title       TEXT            NOT NULL,
    author      CITEXT          NOT NULL,
    forum       CITEXT          NOT NULL,
    message     TEXT            NOT NULL,
    votes       INT DEFAULT 0   NOT NULL,
    created     TIMESTAMP WITH TIME ZONE DEFAULT now(),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug)
);
