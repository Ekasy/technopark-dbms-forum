-- clear tables
-- DROP TABLE IF EXISTS users CASCADE;
-- DROP TABLE IF EXISTS forum CASCADE;
-- DROP TABLE IF EXISTS threads CASCADE;
-- DROP TABLE IF EXISTS posts CASCADE;
-- DROP TABLE IF EXISTS votes CASCADE;

-- install module with case-insensitive string
CREATE EXTENSION IF NOT EXISTS citext;


-----------------------
----- БЛОК ТАБЛИЦ -----
-----------------------
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
    path        BIGINT ARRAY,
    parentRoot  BIGINT,
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


-------------------------
----- БЛОК ИНДЕКСОВ -----
-------------------------
CREATE UNIQUE INDEX IF NOT EXISTS pindex_threads_slug ON threads(slug) WHERE TRIM(slug) <> '';



--------------------------
----- БЛОК ТРИГГЕРОВ -----
--------------------------

-- вставка голоса -> обновление треда
CREATE OR REPLACE FUNCTION vote_insert() RETURNS TRIGGER AS $vote_insert$
BEGIN
    UPDATE threads
    SET votes = votes + NEW.voice
    WHERE id = NEW.thread;
    RETURN NULL;
END;
$vote_insert$  LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS vote_insert ON votes;
CREATE TRIGGER vote_insert AFTER INSERT ON votes FOR EACH ROW EXECUTE PROCEDURE vote_insert();


-- обновление голоса -> обновление треда
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


-- вставка поста -> обновление родительского поста ()
CREATE OR REPLACE FUNCTION update_path() RETURNS TRIGGER AS $update_path$
BEGIN
    IF NEW.parent = 0
    THEN
    
        UPDATE posts
        SET path = ARRAY [NEW.id]
        WHERE id = NEW.id;
    
    ELSE
        
        UPDATE posts SET
            path = array_append(
                (SELECT path FROM posts WHERE id = NEW.parent), 
                NEW.id
            )
        WHERE id = NEW.id;

    END IF;

    RETURN NULL;
END;
$update_path$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_path ON posts;
CREATE TRIGGER update_path AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE update_path();


-- создание поста -> инкремент числа постов в форуме
CREATE OR REPLACE FUNCTION increment_posts_count() RETURNS TRIGGER AS $increment_posts_count$
BEGIN
    UPDATE forum SET 
        posts = (posts + 1)
    WHERE slug = NEW.forum;
    
    RETURN NULL;
END;
$increment_posts_count$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS increment_posts_count ON posts;
CREATE TRIGGER increment_posts_count AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE increment_posts_count();


-- создание трэда -> инкремент числа трэдов в форуме
CREATE OR REPLACE FUNCTION increment_threads_count() RETURNS TRIGGER AS $increment_threads_count$
BEGIN
    UPDATE forum SET 
        threads = (threads + 1)
    WHERE slug = NEW.forum;
    
    RETURN NULL;
END;
$increment_threads_count$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS increment_threads_count ON threads;
CREATE TRIGGER increment_threads_count AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE increment_threads_count();
