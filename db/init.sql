-- clear tables
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS forum CASCADE;
-- DROP TABLE IF EXISTS threads CASCADE;
-- DROP TABLE IF EXISTS posts CASCADE;

-- install module with case-insensitive string
CREATE EXTENSION IF NOT EXISTS citext;

-- tables
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL       UNIQUE,
    nickname    CITEXT       NOT NULL PRIMARY KEY,
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
    id          SERIAL                      NOT NULL PRIMARY KEY,
    slug        CITEXT                      DEFAULT '',
    title       TEXT                        NOT NULL,
    author      CITEXT                      NOT NULL,
    forum       CITEXT                      NOT NULL,
    message     TEXT                        NOT NULL,
    votes       INT DEFAULT 0               NOT NULL,
    created     TIMESTAMP WITH TIME ZONE    DEFAULT now(),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug)
);

CREATE TABLE IF NOT EXISTS posts (
    id          BIGSERIAL                   NOT NULL PRIMARY KEY,
    parent      BIGINT                      NOT NULL DEFAULT 0,
    author      CITEXT                      NOT NULL,
    message     TEXT                        NOT NULL,
    isEdited    BOOLEAN                     NOT NULL DEFAULT FALSE,
    forum       CITEXT                      NOT NULL,
    thread      INTEGER                     NOT NULL,
    created     TIMESTAMP WITH TIME ZONE    NOT NULL DEFAULT NOW(),
    FOREIGN KEY (author) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
    FOREIGN KEY (thread) REFERENCES threads (id)
);

CREATE TABLE IF NOT EXISTS votes (
	nickname 	CITEXT	NOT NULL,
  	thread 		INT		NOT NULL,
  	voice     	INT		NOT NULL,
	FOREIGN KEY (nickname) REFERENCES users(nickname),
	FOREIGN KEY (thread) REFERENCES threads(id)
);


-- create triggers
CREATE OR REPLACE FUNCTION vote_insert()
  RETURNS TRIGGER AS $vote_insert$
    BEGIN
        UPDATE threads
        SET votes = votes + NEW.voice
        WHERE id = NEW.thread;
        RETURN NULL;
    END;
$vote_insert$  LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS vote_insert ON votes;
CREATE TRIGGER vote_insert AFTER INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE vote_insert();

CREATE OR REPLACE FUNCTION vote_update() RETURNS TRIGGER AS $vote_update$
BEGIN
	IF OLD.voice = NEW.voice
		THEN RETURN NULL;
	END IF;
  	UPDATE threads
	SET
		votes = votes + CASE WHEN NEW.voice = -1 THEN -2 ELSE 2 END
  	WHERE id = NEW.thread;
  	RETURN NULL;
END;
$vote_update$ LANGUAGE  plpgsql;

DROP TRIGGER IF EXISTS vote_update ON votes;
CREATE TRIGGER vote_update AFTER UPDATE ON votes FOR EACH ROW EXECUTE PROCEDURE vote_update();