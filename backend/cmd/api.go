package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pranav7002/MVC_Assignment/internal/controller"
	customMiddleware "github.com/pranav7002/MVC_Assignment/internal/middleware"
	"github.com/pranav7002/MVC_Assignment/internal/models"
	"github.com/pranav7002/MVC_Assignment/internal/repository"
	"github.com/pranav7002/MVC_Assignment/internal/services"
)

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://127.0.0.1:3000"}, // Next.js port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)              // rate limiting
	r.Use(middleware.ClientIPFromRemoteAddr) // rate limiting and analytics
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

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

			// SHOP
			r.Get("/shop/buildings", app.shopController.ShopBuildingsHandler)
			r.Get("/shop/troops", app.shopController.ShopTroopsHandler)

			// TROOPS
			r.Post("/troops/train", app.troopController.TrainTroopHandler)
			r.Get("/troops", app.troopController.TroopHandler)
			r.Delete("/troops/{troopName}", app.troopController.TroopDeleteHandler)

			// BATTLE 
			r.Get("/battle/ws/{defendersID}", app.battleController.HandleWebSocket)
		})
	})

	return r
}

func (app *application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
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
	shopController    *controller.ShopController
}

type config struct {
	addr string
	db   dbConfig
}

type dbConfig struct {
	dsn string
}

func (app *application) hydrate(secretKey []byte) {
	var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,              
    WriteBufferSize: 1024,             
    CheckOrigin: func(r *http.Request) bool {
        return true             
    },
	}

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
		WSUpgrader: upgrader,
		BattleManager: &models.BattleManager{Mu: new(sync.Mutex)}, 
	}

	app.authController = authController
	app.villageController = villageController
	app.economyController = economyController
	app.troopController = troopController
	app.battleController = battleController

	shopRepo := &repository.ShopRepository{DB: dbpool}
	shopService := &services.ShopService{
		ShopRepo:    shopRepo,
		VillageRepo: villageRepo,
	}
	app.shopController = &controller.ShopController{ShopService: shopService}
}
