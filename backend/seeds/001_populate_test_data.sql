-- Seed: 5 users (admin + 4 regular, 1 inactive), 3 games, 6 draws, contributions, tickets
-- Password for all users: password123

-- Users (admin active, 3 regular active, 1 inactive for testing)
INSERT INTO users (name, email, password_hash, role, active) VALUES
    ('Admin',    'admin@example.com',    '$2a$10$FKWw7tPmb/UMOnbwWruqdOQbBykHgNwdAUepcQ1tZi8dA2nKp1KDS', 'admin', true),
    ('Alice',    'alice@example.com',    '$2a$10$FKWw7tPmb/UMOnbwWruqdOQbBykHgNwdAUepcQ1tZi8dA2nKp1KDS', 'user',  true),
    ('Bob',      'bob@example.com',      '$2a$10$FKWw7tPmb/UMOnbwWruqdOQbBykHgNwdAUepcQ1tZi8dA2nKp1KDS', 'user',  true),
    ('Carol',    'carol@example.com',    '$2a$10$FKWw7tPmb/UMOnbwWruqdOQbBykHgNwdAUepcQ1tZi8dA2nKp1KDS', 'user',  true),
    ('Dave',     'dave@example.com',     '$2a$10$FKWw7tPmb/UMOnbwWruqdOQbBykHgNwdAUepcQ1tZi8dA2nKp1KDS', 'user',  false);

-- Lottery games
INSERT INTO lottery_games (name, draw_days, ticket_price, active) VALUES
    ('EuroMillones',  'Tuesday, Friday',     2.50, true),
    ('La Primitiva',  'Thursday, Saturday',  1.00, true),
    ('El Gordo',      'Thursday, Sunday',    1.50, true);

-- Draws (2 per game, one with results, one without)
INSERT INTO draws (game_id, draw_date, draw_id_api, result_numbers, result_stars, processed) VALUES
    (1, '2026-09-01', 'EM-2026-09-01', '03,15,22,34,47', '04,09', true),
    (1, '2026-09-08', 'EM-2026-09-08', NULL,              NULL,    false),
    (2, '2026-09-04', 'LP-2026-09-04', '01,08,19,33,41', '07',    true),
    (2, '2026-09-11', 'LP-2026-09-11', NULL,              NULL,    false),
    (3, '2026-09-03', 'EG-2026-09-03', '05,12,28,36,44', '02',    true),
    (3, '2026-09-10', 'EG-2026-09-10', NULL,              NULL,    false);

-- Contributions (active users only: Alice=2, Bob=3, Carol=4)
INSERT INTO contributions (user_id, game_id, month, year, amount, paid, payment_date, payment_method, comments) VALUES
    (2, 1, 9, 2026, 5.00,  true,  '2026-09-01', 'CASH',  'September EuroMillones'),
    (2, 2, 9, 2026, 2.00,  true,  '2026-09-01', 'BIZUM', 'September Primitiva'),
    (3, 1, 9, 2026, 5.00,  true,  '2026-09-02', 'CASH',  'September EuroMillones'),
    (3, 3, 9, 2026, 3.00,  false, NULL,         'BIZUM', 'Pending El Gordo'),
    (4, 1, 9, 2026, 5.00,  false, NULL,         'CASH',  'September EuroMillones'),
    (4, 2, 9, 2026, 2.00,  true,  '2026-09-03', 'BIZUM', 'September Primitiva');

-- Tickets (1 per draw)
INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars) VALUES
    (1, '03,15,22,34,47', '04,09', 5.00, '2026-09-01', '1st',     50000.00, 5, 2),
    (2, '07,14,21,30,42', '03,08', 2.50, '2026-09-07', NULL,      NULL,     NULL, NULL),
    (3, '01,08,19,33,41', '07',    1.00, '2026-09-04', '2nd',     1200.00,  5, 1),
    (4, '05,12,25,38,44', '02',    1.00, '2026-09-10', NULL,      NULL,     NULL, NULL),
    (5, '05,12,28,36,44', '02',    1.50, '2026-09-03', '3rd',     500.00,   4, 1),
    (6, '10,20,30,40,50', '01',    1.50, '2026-09-09', NULL,      NULL,     NULL, NULL);
