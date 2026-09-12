CREATE TABLE IF NOT EXISTS lottery_games (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(100) NOT NULL UNIQUE,
    draw_days    VARCHAR(50) NOT NULL,
    ticket_price DECIMAL(10, 2) NOT NULL,
    active       BOOLEAN NOT NULL DEFAULT true
);
