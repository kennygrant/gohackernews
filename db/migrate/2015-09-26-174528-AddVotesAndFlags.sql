/* SQL migration AddVotesAndFlags */
CREATE TABLE votes (
  created_at timestamp,
  comment_id integer,
  story_id integer,
  user_id integer,
  user_ip text,
  points integer
);
ALTER TABLE votes OWNER TO hackernews_server;
CREATE TABLE flags (
  created_at timestamp,
  comment_id integer,
  story_id integer,
  user_id integer,
  user_ip text,
  points integer
);
ALTER TABLE flags OWNER TO hackernews_server;