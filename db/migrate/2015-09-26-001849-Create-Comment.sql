DROP TABLE IF EXISTS comments;
CREATE TABLE comments (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
text text, 
parent_id integer,
points integer,
user_id integer,
story_id integer
);
ALTER TABLE comments OWNER TO hackernews_server;
