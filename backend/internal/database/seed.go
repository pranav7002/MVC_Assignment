package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Seed(pool *pgxpool.Pool) {
	ctx := context.Background()
	var count int

	// check if troop_config is already seeded
	err := pool.QueryRow(ctx, "SELECT COUNT(*) FROM troop_config").Scan(&count)
	if err != nil {
		log.Println("Failed to count troops: ", err)
		return
	}

	// seed if the database is empty
	if count == 0 {
		log.Println("Seeding Database with Game Configuration...")

		troopQuery := `
			INSERT INTO troop_config (name, dps, health, range, speed, housing_space, training_cost) VALUES 
			('Barbarian', 9, 45, 1, 16, 1, 25),
			('Archer', 8, 22, 4, 24, 1, 50),
			('Goblin', 11, 25, 1, 32, 1, 25),
			('Giant', 12, 400, 1, 12, 4, 150),
			('Wall Breaker', 10, 20, 1, 24, 2, 200);
		`
		if _, err = pool.Exec(ctx, troopQuery); err != nil {
			log.Println("Failed to seed troops: ", err)
		}

		thQuery := `
			INSERT INTO town_hall_config (name, level, upgrade_cost, upgrade_cost_type, upgrade_duration_sec, max_hp) VALUES 
			('Town Hall', 1, 0, 'gold', 0, 1500),
			('Town Hall', 2, 1000, 'gold', 60, 1600),
			('Town Hall', 3, 4000, 'gold', 3600, 1850),
			('Town Hall', 4, 15000, 'gold', 14400, 2100);
		`
		if _, err = pool.Exec(ctx, thQuery); err != nil {
			log.Println("Failed to seed town halls: ", err)
		}

		defQuery := `
			INSERT INTO defense_config (name, level, upgrade_cost, upgrade_cost_type, upgrade_duration_sec, dps, max_hp, range) VALUES 
			('Cannon', 1, 250, 'gold', 10, 9, 420, 9),
			('Cannon', 2, 1000, 'gold', 900, 11, 470, 9),
			('Cannon', 3, 4000, 'gold', 3600, 15, 540, 9),
			('Cannon', 4, 16000, 'gold', 14400, 19, 620, 9),
			
			('Archer Tower', 1, 1000, 'gold', 900, 11, 400, 10),
			('Archer Tower', 2, 2000, 'gold', 1800, 15, 460, 10),
			('Archer Tower', 3, 5000, 'gold', 3600, 19, 520, 10),
			('Archer Tower', 4, 20000, 'gold', 14400, 25, 590, 10),
			
			('Mortar', 1, 8000, 'gold', 7200, 4, 400, 11),
			('Mortar', 2, 32000, 'gold', 43200, 5, 450, 11),
			('Mortar', 3, 120000, 'gold', 86400, 6, 500, 11),
			('Mortar', 4, 180000, 'gold', 172800, 7, 550, 11);
		`
		if _, err = pool.Exec(ctx, defQuery); err != nil {
			log.Println("Failed to seed defences: ", err)
		}

		resQuery := `
			INSERT INTO resource_config (name, level, resource_type, max_capacity, resource_per_sec, upgrade_cost, upgrade_cost_type, upgrade_duration_sec, max_hp) VALUES 
			('Gold Mine', 1, 'gold', 500, 3, 150, 'elixir', 10, 400),
			('Gold Mine', 2, 'gold', 1000, 6, 300, 'elixir', 60, 450),
			('Gold Mine', 3, 'gold', 1500, 10, 700, 'elixir', 900, 500),
			('Gold Mine', 4, 'gold', 2500, 13, 1400, 'elixir', 3600, 550),
			
			('Elixir Collector', 1, 'elixir', 500, 3, 150, 'gold', 10, 400),
			('Elixir Collector', 2, 'elixir', 1000, 6, 300, 'gold', 60, 450),
			('Elixir Collector', 3, 'elixir', 1500, 10, 700, 'gold', 900, 500),
			('Elixir Collector', 4, 'elixir', 2500, 13, 1400, 'gold', 3600, 550);
		`
		if _, err = pool.Exec(ctx, resQuery); err != nil {
			log.Println("Failed to seed resources: ", err)
		}

		storeQuery := `
			INSERT INTO storage_config (name, level, resource_type, max_capacity, upgrade_cost, upgrade_cost_type, upgrade_duration_sec, max_hp) VALUES 
			('Gold Storage', 1, 'gold', 1500, 300, 'elixir', 60, 400),
			('Gold Storage', 2, 'gold', 3000, 750, 'elixir', 1800, 450),
			('Gold Storage', 3, 'gold', 6000, 1500, 'elixir', 7200, 500),
			('Gold Storage', 4, 'gold', 12000, 3000, 'elixir', 14400, 550),
			
			('Elixir Storage', 1, 'elixir', 1500, 300, 'gold', 60, 400),
			('Elixir Storage', 2, 'elixir', 3000, 750, 'gold', 1800, 450),
			('Elixir Storage', 3, 'elixir', 6000, 1500, 'gold', 7200, 500),
			('Elixir Storage', 4, 'elixir', 12000, 3000, 'gold', 14400, 550);
		`
		if _, err = pool.Exec(ctx, storeQuery); err != nil {
			log.Println("Failed to seed storages: ", err)
		}

		progQuery := `
			INSERT INTO game_progression_config (town_hall_level, building_type, building_name, max_level, max_built) VALUES 
			-- TH1
			(1, 'defense', 'Cannon', 2, 2),
			(1, 'defense', 'Archer Tower', 0, 0),
			(1, 'defense', 'Mortar', 0, 0),
			(1, 'resource', 'Gold Mine', 2, 1),
			(1, 'resource', 'Elixir Collector', 2, 1),
			(1, 'storage', 'Gold Storage', 2, 1),
			(1, 'storage', 'Elixir Storage', 2, 1),
			(1, 'training_grounds', 'Training Grounds', 0, 0),
			
			-- TH2
			(2, 'defense', 'Cannon', 3, 2),
			(2, 'defense', 'Archer Tower', 2, 1),
			(2, 'defense', 'Mortar', 0, 0),
			(2, 'resource', 'Gold Mine', 4, 2),
			(2, 'resource', 'Elixir Collector', 4, 2),
			(2, 'storage', 'Gold Storage', 4, 1),
			(2, 'storage', 'Elixir Storage', 4, 1),
			(2, 'training_grounds', 'Training Grounds', 1, 1),
			
			-- TH3
			(3, 'defense', 'Cannon', 4, 2),
			(3, 'defense', 'Archer Tower', 3, 1),
			(3, 'defense', 'Mortar', 1, 1),
			(3, 'resource', 'Gold Mine', 6, 3),
			(3, 'resource', 'Elixir Collector', 6, 3),
			(3, 'storage', 'Gold Storage', 6, 2),
			(3, 'storage', 'Elixir Storage', 6, 2),
			(3, 'training_grounds', 'Training Grounds', 2, 1),
			
			-- TH4
			(4, 'defense', 'Cannon', 5, 2),
			(4, 'defense', 'Archer Tower', 4, 2),
			(4, 'defense', 'Mortar', 2, 1),
			(4, 'resource', 'Gold Mine', 8, 4),
			(4, 'resource', 'Elixir Collector', 8, 4),
			(4, 'storage', 'Gold Storage', 8, 2),
			(4, 'storage', 'Elixir Storage', 8, 2),
			(4, 'training_grounds', 'Training Grounds', 3, 2);
		`
		if _, err = pool.Exec(ctx, progQuery); err != nil {
			log.Println("Failed to seed progression rules: ", err)
		}

		log.Println("Database successfully seeded with base game configuration!")
	}
}
