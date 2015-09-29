/* SQL migration AddUsers */

CREATE TABLE users (
    id SERIAL NOT NULL,
    created_at timestamp,
    updated_at timestamp,
    status int,
    role int,
    email text,
    name text,
    summary text,
    encrypted_password text,
    points int
);
ALTER TABLE users OWNER TO hackernews_server;