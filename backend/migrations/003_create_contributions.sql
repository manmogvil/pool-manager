CREATE TABLE IF NOT EXISTS contributions (
    id             SERIAL PRIMARY KEY,
    participant_id INTEGER NOT NULL REFERENCES participants(id),
    game_id        INTEGER NOT NULL REFERENCES lottery_games(id),
    month          INTEGER NOT NULL CHECK (month >= 1 AND month <= 12),
    year           INTEGER NOT NULL CHECK (year >= 2020),
    amount         DECIMAL(10, 2) NOT NULL,
    paid           BOOLEAN NOT NULL DEFAULT false,
    payment_date   DATE,
    payment_method VARCHAR(50) NOT NULL CHECK (payment_method IN ('CASH', 'BIZUM')),
    comments       TEXT,
    created_at     TIMESTAMP NOT NULL DEFAULT NOW()
);