# Create the table for the realtime data
CREATE TABLE videos (
    source_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
    chat_id VARCHAR(255) NOT NULL,
    scheduled_at TIMESTAMP,
    updated_at TIMESTAMP NOT NULL
);

# add unique index
CREATE UNIQUE INDEX videos_source_id_idx ON videos (source_id);
