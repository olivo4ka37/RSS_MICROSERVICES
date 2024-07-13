-- +goose Up
CREATE TABLE Sources (
                         id SERIAL PRIMARY KEY,
                         url TEXT NOT NULL
);

ALTER TABLE Sources ADD CONSTRAINT unique_url UNIQUE (url);

CREATE TABLE articles (
                          id SERIAL PRIMARY KEY,
                          title TEXT NOT NULL,
                          link TEXT NOT NULL,
                          description TEXT,
                          published TIMESTAMP,
                          source_id INTEGER NOT NULL,
                          FOREIGN KEY (source_id) REFERENCES Sources (id)
);
CREATE UNIQUE INDEX unique_link ON articles(link);

CREATE TABLE administrators (
                                id SERIAL PRIMARY KEY,
                                uuid UUID NOT NULL UNIQUE
);

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       uuid UUID NOT NULL UNIQUE,
                       last_login TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_subscriptions (
                                    user_id INTEGER NOT NULL,
                                    source_id INTEGER NOT NULL,
                                    PRIMARY KEY (user_id, source_id),
                                    FOREIGN KEY (user_id) REFERENCES users (id),
                                    FOREIGN KEY (source_id) REFERENCES Sources (id)
);

INSERT INTO administrators (uuid) VALUES ('490d5457-17ec-4a36-a412-5208f6630505');
INSERT INTO users (uuid) VALUES ('f4c1755f-5991-4750-91ba-a42f93fd1d4c');

-- +goose Down
DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS administrators;
DROP TABLE IF EXISTS users;
ALTER TABLE Sources DROP CONSTRAINT unique_url;
DROP TABLE IF EXISTS Sources;
