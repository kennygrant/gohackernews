/* Setup tables for gohackernews */
CREATE TABLE fragmenta_metadata (
id SERIAL NOT NULL,
updated_at timestamp,
fragmenta_version text,
migration_version text,
status integer
);

CREATE TABLE comments (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
parent_id integer,
dotted_ids text, 
points integer,
rank integer,
level integer,
text text, 
user_name text,
story_name text,
user_id integer,
story_id integer
);

CREATE TABLE stories (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
name text,
url text,
rank integer,
summary text,
user_id integer,
user_name text,
points integer,
comment_count integer
);

CREATE TABLE votes (
created_at timestamp,
comment_id integer,
story_id integer,
user_id integer,
user_ip text,
points integer
);

CREATE TABLE flags (
created_at timestamp,
comment_id integer,
story_id integer,
user_id integer,
user_ip text,
points integer
);

CREATE TABLE users (
id SERIAL NOT NULL,
created_at timestamp,
updated_at timestamp,
status integer,
role integer,
email text,
name text,
summary text,
encrypted_password text,
points integer
);

ALTER TABLE fragmenta_metadata OWNER TO gohackernews_server;
ALTER TABLE comments OWNER TO gohackernews_server;
ALTER TABLE flags OWNER TO gohackernews_server;
ALTER TABLE users OWNER TO gohackernews_server;
ALTER TABLE votes OWNER TO gohackernews_server;
ALTER TABLE stories OWNER TO gohackernews_server;
