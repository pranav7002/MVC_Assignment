CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create ENUMs
CREATE TYPE building_type AS ENUM ('defense', 'resource', 'storage', 'training_grounds', 'town_hall');
CREATE TYPE battle_result AS ENUM ('LOSS', 'ONE_STAR', 'TWO_STARS', 'THREE_STARS');
CREATE TYPE resource_type AS ENUM ('gold', 'elixir');

-- Users Table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    xp INT NOT NULL DEFAULT 0,
    level INT NOT NULL DEFAULT 1,
    trophies INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Village Table
CREATE TABLE village (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    town_hall_level INT NOT NULL DEFAULT 1,
    gold INT NOT NULL DEFAULT 0,
    elixir INT NOT NULL DEFAULT 0,
    gold_last_collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    elixir_last_collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Troops Trained
CREATE TABLE troops_trained (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    troop_name VARCHAR(50) NOT NULL,
    quantity INT NOT NULL DEFAULT 0,
    UNIQUE(user_id, troop_name)
);
CREATE INDEX idx_troops_trained_user_id ON troops_trained(user_id);

-- Battles
CREATE TABLE battles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    attacker_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    defender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    gold_looted INT NOT NULL DEFAULT 0,
    elixir_looted INT NOT NULL DEFAULT 0,
    result battle_result NOT NULL,
    destruction_percentage INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Building Instance
CREATE TABLE building_instance (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    building_type building_type NOT NULL,
    building_name VARCHAR(50) NOT NULL,
    level INT NOT NULL DEFAULT 1,
    pos_x INT NOT NULL,
    pos_y INT NOT NULL,
    size INT NOT NULL,
    is_upgrading BOOLEAN NOT NULL DEFAULT false,
    hp INT NOT NULL
);
CREATE INDEX idx_building_instance_user_pos ON building_instance(user_id, pos_x, pos_y);

-- Config Tables 
CREATE TABLE game_progression_config (
    id BIGSERIAL PRIMARY KEY,
    town_hall_level INT NOT NULL,
    building_type building_type NOT NULL,
    building_name VARCHAR(50) NOT NULL,
    max_level BIGINT NOT NULL,
    max_built BIGINT NOT NULL
);

CREATE TABLE storage_config (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    resource_type resource_type NOT NULL,
    max_capacity INT NOT NULL,
    upgrade_cost INT NOT NULL,
    upgrade_cost_type resource_type NOT NULL,
    upgrade_duration_sec INT NOT NULL,
    max_hp INT NOT NULL,
    size INT NOT NULL DEFAULT 3,
    PRIMARY KEY (name, level)
);

CREATE TABLE resource_config (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    resource_type resource_type NOT NULL,
    max_capacity INT NOT NULL,
    resource_per_sec INT NOT NULL,
    upgrade_cost INT NOT NULL,
    upgrade_cost_type resource_type NOT NULL,
    upgrade_duration_sec INT NOT NULL,
    max_hp INT NOT NULL,
    size INT NOT NULL DEFAULT 3,
    PRIMARY KEY (name, level)
);

CREATE TABLE town_hall_config (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    upgrade_cost INT NOT NULL,
    upgrade_cost_type resource_type NOT NULL,
    upgrade_duration_sec INT NOT NULL,
    max_hp INT NOT NULL,
    size INT NOT NULL DEFAULT 4,
    PRIMARY KEY (name, level)
);

CREATE TABLE defense_config (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    upgrade_cost INT NOT NULL,
    upgrade_cost_type resource_type NOT NULL,
    upgrade_duration_sec INT NOT NULL,
    dps INT NOT NULL,
    max_hp INT NOT NULL,
    max_range INT NOT NULL,
    min_range INT NOT NULL DEFAULT 0,
    aoe_range INT NOT NULL DEFAULT 0,
    size INT NOT NULL DEFAULT 3,
    PRIMARY KEY (name, level)
);

CREATE TABLE training_grounds_config (
    name VARCHAR(50) NOT NULL,
    level INT NOT NULL,
    housing_space INT NOT NULL,
    upgrade_cost INT NOT NULL,
    upgrade_cost_type resource_type NOT NULL,
    upgrade_duration_sec INT NOT NULL,
    max_hp INT NOT NULL,
    size INT NOT NULL DEFAULT 3,
    PRIMARY KEY (name, level)
);

CREATE TABLE troop_config (
    name VARCHAR(50) PRIMARY KEY,
    dps INT NOT NULL,
    health INT NOT NULL,
    range INT NOT NULL,
    housing_space INT NOT NULL,
    training_cost INT NOT NULL
);
