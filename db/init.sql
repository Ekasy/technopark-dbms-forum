CREATE EXTENSION IF NOT EXISTS citext;


-----------------------
----- БЛОК ТАБЛИЦ -----
-----------------------
-- очистка таблиц
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS forum CASCADE;
DROP TABLE IF EXISTS threads CASCADE;
DROP TABLE IF EXISTS posts CASCADE;
DROP TABLE IF EXISTS votes CASCADE;
DROP TABLE IF EXISTS forum_users CASCADE;


CREATE TABLE IF NOT EXISTS users (
    nickname    CITEXT COLLATE "C"  NOT NULL PRIMARY KEY,
    fullname    TEXT                NOT NULL,
    email       CITEXT              NOT NULL UNIQUE,
    about       TEXT                NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS forum (
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
    parent      BIGINT                      NOT NULL,
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
	FOREIGN KEY (thread) REFERENCES threads(id),
    PRIMARY KEY (nickname, thread)
);

CREATE TABLE IF NOT EXISTS forum_users (
    nickname    CITEXT COLLATE "C"  NOT NULL,
    fullname    TEXT                NOT NULL,
    email       CITEXT              NOT NULL,
    about       TEXT                NOT NULL DEFAULT '',
    forum       CITEXT              NOT NULL,
    FOREIGN KEY (nickname) REFERENCES users (nickname),
    FOREIGN KEY (forum) REFERENCES forum (slug),
	PRIMARY KEY (nickname, forum)
);


-------------------------
----- БЛОК ИНДЕКСОВ -----
-------------------------
-- очистка индексов
DROP INDEX IF EXISTS index_users_nickname;
DROP INDEX IF EXISTS index_users_email;
DROP INDEX IF EXISTS index_forum_slug;
DROP INDEX IF EXISTS index_threads_id;
DROP INDEX IF EXISTS index_threads_slug;
DROP INDEX IF EXISTS composite_index_threads_id_slug;
DROP INDEX IF EXISTS composite_index_threads_forum_created;
DROP INDEX IF EXISTS composite_index_threads_created;
DROP INDEX IF EXISTS index_posts_id;
DROP INDEX IF EXISTS index_posts_thread_parent;
DROP INDEX IF EXISTS index_posts_thread_id;
DROP INDEX IF EXISTS index_posts_thread_path;
DROP INDEX IF EXISTS index_posts_thread_created_id;
DROP INDEX IF EXISTS index_posts_path_1_path;
DROP INDEX IF EXISTS index_forum_users_nickname;
DROP INDEX IF EXISTS index_forum_users_forum;

-- индексы для users

-- индексы для forum

-- индексы для threads
CREATE UNIQUE INDEX IF NOT EXISTS index_threads_slug                    ON threads(slug) WHERE TRIM(slug) <> '';
CREATE UNIQUE INDEX IF NOT EXISTS composite_index_threads_id_slug       ON threads(id, slug);
CREATE        INDEX IF NOT EXISTS composite_index_threads_forum_created ON threads(forum, created);
CLUSTER threads USING composite_index_threads_forum_created;

-- индексы для posts
CREATE UNIQUE INDEX IF NOT EXISTS index_posts_id_parent			ON posts(thread, id) WHERE parent != 0;

-- индексы для forum_users
CREATE INDEX IF NOT EXISTS index_forum_users_nickname ON forum_users(nickname, forum);


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


-- создание поста -> вставка юзера в forum_users
CREATE OR REPLACE FUNCTION post_paste_forum_user() RETURNS TRIGGER AS $post_paste_forum_user$
BEGIN
    INSERT INTO forum_users
    SELECT nickname, fullname, email, about, NEW.forum as forum 
    FROM users
    WHERE nickname = NEW.author
	ON CONFLICT DO NOTHING;
    
    RETURN NULL;
END;
$post_paste_forum_user$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS post_paste_forum_user ON posts;
CREATE TRIGGER post_paste_forum_user AFTER INSERT ON posts FOR EACH ROW EXECUTE PROCEDURE post_paste_forum_user();


-- создание трэда -> вставка юзера в forum_users
CREATE OR REPLACE FUNCTION thread_paste_forum_user() RETURNS TRIGGER AS $thread_paste_forum_user$
BEGIN
    INSERT INTO forum_users
    SELECT nickname, fullname, email, about, NEW.forum as forum 
    FROM users
    WHERE nickname = NEW.author
	ON CONFLICT DO NOTHING;
    
    RETURN NULL;
END;
$thread_paste_forum_user$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS thread_paste_forum_user ON threads;
CREATE TRIGGER thread_paste_forum_user AFTER INSERT ON threads FOR EACH ROW EXECUTE PROCEDURE thread_paste_forum_user();
