CREATE TABLE IF NOT EXISTS participants (
    id             SERIAL PRIMARY KEY,
    name           VARCHAR(255) NOT NULL,
    email          VARCHAR(255) NOT NULL UNIQUE,
    active         BOOLEAN NOT NULL DEFAULT true,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);
