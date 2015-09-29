/* Setup tables for cms */
CREATE TABLE fragmenta_metadata (
    id SERIAL NOT NULL,
    updated_at timestamp,
    fragmenta_version text,
    migration_version text,
    status int
);

ALTER TABLE fragmenta_metadata OWNER TO "hackernews_server";
