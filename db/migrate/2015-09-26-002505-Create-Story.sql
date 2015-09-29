DROP TABLE IF EXISTS stories;
CREATE TABLE stories (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
name text,
url text,
summary text,
user_id integer,
points integer
);
ALTER TABLE stories OWNER TO hackernews_server;
