CREATE TABLE IF NOT EXISTS draws (
    id              SERIAL PRIMARY KEY,
    game_id         INTEGER REFERENCES lottery_games(id) ON DELETE SET NULL,
    draw_date       DATE NOT NULL,
    result_numbers  VARCHAR(100),
    result_stars    VARCHAR(50),
    processed       BOOLEAN NOT NULL DEFAULT false,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_draws_game_date ON draws(game_id, draw_date);
