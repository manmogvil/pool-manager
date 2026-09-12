CREATE TABLE IF NOT EXISTS tickets (
    id              SERIAL PRIMARY KEY,
    draw_id         INTEGER REFERENCES draws(id) ON DELETE SET NULL,
    numbers         VARCHAR(100) NOT NULL,
    stars           VARCHAR(50),
    cost            DECIMAL(10, 2) NOT NULL,
    purchased_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    prize_tier      VARCHAR(50),
    prize_amount    DECIMAL(10, 2),
    matched_numbers INTEGER,
    matched_stars   INTEGER
);

CREATE INDEX idx_tickets_draw ON tickets(draw_id);
