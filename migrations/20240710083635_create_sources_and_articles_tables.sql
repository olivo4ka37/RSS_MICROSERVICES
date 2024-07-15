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
-- INSERT INTO users (uuid) VALUES ('f4c1755f-5991-4750-91ba-a42f93fd1d4c');
INSERT INTO users (uuid, last_login) VALUES ('f4c1755f-5991-4750-91ba-a42f93fd1d4c', '2023-01-01 00:00:00');
INSERT INTO Sources (url)
VALUES
    ('https://habr.com/ru/rss/hub/go/all/?fl=ru'),
    ('https://habr.com/ru/rss/best/daily/?fl=ru'),
    ('https://golangcode.com/index.xml'),
    ('https://forum.golangbridge.org/latest.rss'),
    ('https://appliedgo.net/index.xml'),
    ('https://blog.jetbrains.com/go/feed/'),
    ('https://dave.cheney.net/category/golang/feed'),
    ('https://changelog.com/gotime/feed'),
    ('https://golang.ch/feed/'),
    ('https://gosamples.dev/index.xml'),
    ('https://www.coolfields.co.uk/feed/'),
    ('https://webplatform.news/feed.xml');
INSERT INTO user_subscriptions (user_id, source_id) VALUES ('1','1'), ('1','2'), ('1','3') , ('1','4'),('1','5'), ('1','6'), ('1','7') , ('1','8'),('1','9'), ('1','10'), ('1','11') , ('1','12');
SELECT id, title, link, description, published, source_id
FROM articles
WHERE source_id = '1' AND published > '2023-01-01 00:00:00'
ORDER BY published DESC
    LIMIT '3' OFFSET '0';


-- +goose Down
DROP TABLE IF EXISTS user_subscriptions;
DROP TABLE IF EXISTS articles;
DROP TABLE IF EXISTS administrators;
DROP TABLE IF EXISTS users;
ALTER TABLE Sources DROP CONSTRAINT unique_url;
DROP TABLE IF EXISTS Sources;
