package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/controller"
	customMiddleware "github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/repository"
	"github.com/pranav7002/MVC_Assignment/internal/services"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)              // rate limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // rate limiting and analytics
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api", func(r chi.Router) {
		// PUBLIC
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", app.authController.RegisterHandler)
			r.Post("/auth/login", app.authController.LoginHandler)
		})

		// PROTECTED
		r.Group(func(r chi.Router) {
			r.Use(customMiddleware.JWTMiddleware)

			// VILLAGE BUILDING
			r.Get("/buildings", app.villageController.BuildingHandler)
			r.Get("/village", app.villageController.VillageHandler)
			r.Post("/buildings", app.villageController.BuildingCreationHandler)
			r.Put("/buildings/{buildingID}/move", app.villageController.BuildingPositionHandler)
			r.Put("/buildings/{buildingID}/upgrade", app.villageController.BuildingUpgradeHandler)

			// ECONOMY
			r.Post("/economy/collect", app.economyController.ResourceCollectionHandler)

			// TROOPS
			r.Post("/troops/train", app.troopController.TrainTroopHandler)
			r.Get("/troops", app.troopController.TroopHandler)
			r.Delete("/troops/{troopName}", app.troopController.TroopDeleteHandler)

			// BATTLE 
			r.Post("/battle/{defendersID}", app.battleController.RunSimulation)
		})
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}

func (app *application) connectDB() (*pgxpool.Pool, error) {
	dbpool, err := pgxpool.New(context.Background(), app.config.db.dsn)
	if err != nil {
		return nil, err
	}

	log.Printf("db connected successfully")
	return dbpool, nil
}

type application struct {
	config config
	dbpool *pgxpool.Pool

	authController    *controller.AuthController
	villageController *controller.VillageController
	economyController *controller.EconomyController
	troopController   *controller.TroopController
	battleController  *controller.BattleController
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) hydrate(secretKey []byte) {
	dbpool := app.dbpool

	userRepo := &repository.UserRepository{DB: dbpool}
	authService := &services.AuthService{
		UserRepo: userRepo,
		SecretKey: secretKey,
	}
	authController := &controller.AuthController{AuthService: authService}

	villageRepo := &repository.VillageRepository{DB: dbpool}
	configRepo := &repository.ConfigRepository{DB: dbpool}
	villageService := &services.VillageService{
		VillageRepo: villageRepo,
		ConfigRepo:  configRepo,
	}
	villageController := &controller.VillageController{VillageService: villageService}

	economyService := &services.EconomyService{
		VillageRepo: villageRepo,
		UserRepo: userRepo,
		ConfigRepo: configRepo,
	}
	economyController := &controller.EconomyController{EconomyService: economyService}

	troopRepo := &repository.TroopRepository{DB: dbpool}
	troopService := &services.TroopService{
		TroopRepo:   troopRepo,
		VillageRepo: villageRepo,
		ConfigRepo:  configRepo,
	}
	troopController := &controller.TroopController{TroopService: troopService}

	battleRepo := &repository.BattleRepository{DB: dbpool}
	battleService := &services.BattleService{
		BattleRepo:  battleRepo,
		VillageRepo: villageRepo,
		ConfigRepo:  configRepo,
	}
	battleController := &controller.BattleController{
		BattleService:  battleService,
		VillageService: villageService,
	}

	app.authController = authController
	app.villageController = villageController
	app.economyController = economyController
	app.troopController = troopController
	app.battleController = battleController
}
