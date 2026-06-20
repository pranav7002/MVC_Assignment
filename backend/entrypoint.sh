#!/bin/sh
set -e

DB_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"

echo "Running migrations..."
migrate -path ./migrations -database "$DB_URL" up

echo "Seeding users..."
psql "$DB_URL" <<'ENDSQL'

-- Kai (rich test user)
DO $$
DECLARE
    uid UUID;
    pw TEXT := '$2a$10$fiZSIc3tSM5kGXdtMX.ON.6YOrQWWMZihmppvw3ast4nzchFSlf/a';
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE username = 'kai') THEN
        INSERT INTO users (username, password_hash, trophies)
        VALUES ('kai', pw, 0) RETURNING id INTO uid;

        INSERT INTO village (user_id, town_hall_level, gold, elixir, max_gold, max_elixir)
        VALUES (uid, 1, 1000000, 1000000, 1000000, 1000000);

        INSERT INTO building_instance (user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) VALUES
        (uid, 'town_hall',        'Town Hall',        1, 8,  8,  4, false, 1500),
        (uid, 'training_grounds', 'Training Grounds', 1, 4,  12, 3, false, 400),
        (uid, 'storage',          'Gold Storage',     1, 4,  4,  3, false, 400),
        (uid, 'storage',          'Elixir Storage',   1, 13, 13, 3, false, 400);

        RAISE NOTICE 'Seeded kai (password: password)';
    END IF;
END $$;

-- 15 Bot players
DO $$
DECLARE
    uid UUID;
    pw TEXT := '$2a$10$fiZSIc3tSM5kGXdtMX.ON.6YOrQWWMZihmppvw3ast4nzchFSlf/a';
    cnt INT;
BEGIN
    SELECT COUNT(*) INTO cnt FROM users WHERE username LIKE 'player_%';
    IF cnt >= 15 THEN
        RAISE NOTICE 'Players already seeded, skipping.';
        RETURN;
    END IF;

    -- TH1 players
    FOR i IN 1..5 LOOP
        INSERT INTO users (username, password_hash, trophies)
        VALUES ('player_' || i, pw, i * 30) RETURNING id INTO uid;

        INSERT INTO village (user_id, town_hall_level, gold, elixir, max_gold, max_elixir)
        VALUES (uid, 1, 500, 500, 1500, 1500);

        INSERT INTO building_instance (user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) VALUES
        (uid, 'town_hall',        'Town Hall',        1, 8,  8,  4, false, 1500),
        (uid, 'defense',          'Cannon',           1, 3,  3,  2, false, 420),
        (uid, 'defense',          'Archer Tower',     1, 14, 3,  2, false, 400),
        (uid, 'resource',         'Gold Mine',        1, 3,  14, 2, false, 400),
        (uid, 'resource',         'Elixir Collector', 1, 14, 14, 2, false, 400),
        (uid, 'storage',          'Gold Storage',     1, 3,  7,  3, false, 400),
        (uid, 'storage',          'Elixir Storage',   1, 14, 7,  3, false, 400),
        (uid, 'training_grounds', 'Training Grounds', 1, 8,  14, 3, false, 400);
    END LOOP;

    -- TH2 players
    FOR i IN 6..10 LOOP
        INSERT INTO users (username, password_hash, trophies)
        VALUES ('player_' || i, pw, i * 30) RETURNING id INTO uid;

        INSERT INTO village (user_id, town_hall_level, gold, elixir, max_gold, max_elixir)
        VALUES (uid, 2, 1200, 1200, 3000, 3000);

        INSERT INTO building_instance (user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) VALUES
        (uid, 'town_hall',        'Town Hall',        2, 8,  8,  4, false, 1600),
        (uid, 'defense',          'Cannon',           2, 3,  3,  2, false, 470),
        (uid, 'defense',          'Cannon',           1, 15, 3,  2, false, 420),
        (uid, 'defense',          'Archer Tower',     2, 3,  15, 2, false, 460),
        (uid, 'resource',         'Gold Mine',        1, 15, 15, 2, false, 400),
        (uid, 'resource',         'Elixir Collector', 1, 13, 1,  2, false, 400),
        (uid, 'storage',          'Gold Storage',     1, 3,  7,  3, false, 400),
        (uid, 'storage',          'Elixir Storage',   1, 14, 7,  3, false, 400),
        (uid, 'training_grounds', 'Training Grounds', 1, 8,  14, 3, false, 400);
    END LOOP;

    -- TH3 players
    FOR i IN 11..15 LOOP
        INSERT INTO users (username, password_hash, trophies)
        VALUES ('player_' || i, pw, i * 30) RETURNING id INTO uid;

        INSERT INTO village (user_id, town_hall_level, gold, elixir, max_gold, max_elixir)
        VALUES (uid, 3, 3000, 3000, 6000, 6000);

        INSERT INTO building_instance (user_id, building_type, building_name, level, pos_x, pos_y, size, is_upgrading, hp) VALUES
        (uid, 'town_hall',        'Town Hall',        3, 8,  8,  4, false, 1700),
        (uid, 'defense',          'Cannon',           2, 3,  3,  2, false, 470),
        (uid, 'defense',          'Cannon',           2, 15, 3,  2, false, 470),
        (uid, 'defense',          'Archer Tower',     2, 3,  15, 2, false, 460),
        (uid, 'defense',          'Archer Tower',     1, 15, 15, 2, false, 400),
        (uid, 'defense',          'Mortar',           1, 1,  9,  2, false, 400),
        (uid, 'resource',         'Gold Mine',        1, 13, 1,  2, false, 400),
        (uid, 'resource',         'Elixir Collector', 1, 1,  1,  2, false, 400),
        (uid, 'storage',          'Gold Storage',     2, 1,  13, 3, false, 450),
        (uid, 'storage',          'Elixir Storage',   2, 16, 8,  3, false, 450),
        (uid, 'training_grounds', 'Training Grounds', 1, 8,  14, 3, false, 400);
    END LOOP;

    RAISE NOTICE 'Seeded player_1 to player_15 (password: password)';
END $$;

ENDSQL

echo "Starting server..."
go run ./cmd
