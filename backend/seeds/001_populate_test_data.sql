-- Seed data for Lottery Pool Manager
-- Run: docker exec -i lottery-postgres psql -U lottery -d lottery_pool < backend/migrations/006_seed_data.sql

-- Lottery Games
INSERT INTO lottery_games (name, draw_days, ticket_price, active) VALUES
    ('EuroMillones', 'Tue,Fri', 2.50, true),
    ('La Primitiva', 'Thu,Sun', 1.00, true),
    ('Bonoloto', 'Mon,Tue,Wed,Thu,Fri,Sat,Sun', 0.50, true),
    ('El Gordo de la Primitiva', 'Sun', 1.50, true),
    ('Lotería Nacional', 'Thu,Sat', 3.00, true)
ON CONFLICT (name) DO NOTHING;

-- Participants
INSERT INTO participants (name, email, active) VALUES
    ('Manuel', 'manuel@example.com', true),
    ('Laura', 'laura@example.com', true),
    ('Carlos', 'carlos@example.com', true),
    ('María', 'maria@example.com', true),
    ('Pedro', 'pedro@example.com', true)
ON CONFLICT (email) DO NOTHING;

-- Draws (EuroMillones)
INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 1, '2026-09-05', '5,12,23,34,45', '3,7', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 1 AND draw_date = '2026-09-05');

INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 1, '2026-09-02', '8,15,22,31,40', '2,9', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 1 AND draw_date = '2026-09-02');

INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 1, '2026-08-29', '3,11,18,27,38', '5,11', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 1 AND draw_date = '2026-08-29');

-- Draws (La Primitiva)
INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 2, '2026-09-07', '2,10,19,28,35,41', '', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 2 AND draw_date = '2026-09-07');

INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 2, '2026-09-04', '6,13,24,30,37,44', '', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 2 AND draw_date = '2026-09-04');

-- Draws (Bonoloto)
INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 3, '2026-09-08', '4,16,25,32,39', '', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 3 AND draw_date = '2026-09-08');

INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 3, '2026-09-07', '11,20,29,36,41', '', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 3 AND draw_date = '2026-09-07');

INSERT INTO draws (game_id, draw_date, result_numbers, result_stars, processed)
SELECT 3, '2026-09-06', '7,14,22,30,38', '', true
WHERE NOT EXISTS (SELECT 1 FROM draws WHERE game_id = 3 AND draw_date = '2026-09-06');

-- Contributions
INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 1, 1, 9, 2026, 25.00, true, 'BIZUM', 'EuroMillones septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 1 AND game_id = 1 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 1, 2, 9, 2026, 10.00, true, 'BIZUM', 'La Primitiva septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 1 AND game_id = 2 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 1, 3, 9, 2026, 5.00, true, 'CASH', 'Bonoloto septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 1 AND game_id = 3 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 2, 1, 9, 2026, 25.00, true, 'CASH', 'EuroMillones septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 2 AND game_id = 1 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 2, 2, 9, 2026, 10.00, false, 'CASH', ''
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 2 AND game_id = 2 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 3, 1, 9, 2026, 25.00, true, 'BIZUM', 'EuroMillones septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 3 AND game_id = 1 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 3, 3, 9, 2026, 5.00, true, 'CASH', 'Bonoloto septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 3 AND game_id = 3 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 4, 1, 9, 2026, 25.00, true, 'BIZUM', 'EuroMillones septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 4 AND game_id = 1 AND month = 9 AND year = 2026);

INSERT INTO contributions (participant_id, game_id, month, year, amount, paid, payment_method, comments)
SELECT 5, 2, 9, 2026, 10.00, true, 'CASH', 'La Primitiva septiembre'
WHERE NOT EXISTS (SELECT 1 FROM contributions WHERE participant_id = 5 AND game_id = 2 AND month = 9 AND year = 2026);

-- Tickets (EuroMillones draw_id = first draw inserted)
-- We use a subquery to find the draw_id dynamically
INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,23,34,45', '3,7', 2.50, '2026-09-03 10:00:00', '1st', 50000000.00, 5, 2
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,23,34,45');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,23,34,46', '3,7', 2.50, '2026-09-03 11:30:00', '2nd', 100000.00, 5, 1
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,23,34,46');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,23,34,45', '3,8', 2.50, '2026-09-04 09:15:00', '3rd', 5000.00, 5, 1
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,23,34,45' AND t.stars = '3,8');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '1,12,23,34,45', '3,7', 2.50, '2026-09-04 14:00:00', '4th', 100.00, 4, 2
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '1,12,23,34,45');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,23,34,50', '3,7', 2.50, '2026-09-05 08:00:00', '5th', 50.00, 4, 1
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,23,34,50');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '2,12,23,34,45', '3,7', 2.50, '2026-09-05 09:30:00', '', 0, 3, 2
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '2,12,23,34,45');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,23,40,45', '3,7', 2.50, '2026-09-05 10:45:00', '', 0, 3, 1
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,23,40,45');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '5,12,30,34,45', '3,7', 2.50, '2026-09-05 11:00:00', '', 0, 3, 1
FROM draws d WHERE d.game_id = 1 AND d.draw_date = '2026-09-05'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '5,12,30,34,45');

-- Tickets (La Primitiva)
INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '2,10,19,28,35,41', '', 1.00, '2026-09-05 12:00:00', '1st', 5000000.00, 6, 0
FROM draws d WHERE d.game_id = 2 AND d.draw_date = '2026-09-07'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '2,10,19,28,35,41');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '2,10,19,28,35,42', '', 1.00, '2026-09-06 10:00:00', '2nd', 100000.00, 5, 0
FROM draws d WHERE d.game_id = 2 AND d.draw_date = '2026-09-07'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '2,10,19,28,35,42');

INSERT INTO tickets (draw_id, numbers, stars, cost, purchased_at, prize_tier, prize_amount, matched_numbers, matched_stars)
SELECT d.id, '2,10,19,28,36,41', '', 1.00, '2026-09-06 11:30:00', '3rd', 1000.00, 5, 0
FROM draws d WHERE d.game_id = 2 AND d.draw_date = '2026-09-07'
AND NOT EXISTS (SELECT 1 FROM tickets t WHERE t.draw_id = d.id AND t.numbers = '2,10,19,28,36,41');
