package main

import (
	"log"
	"net/http"
	"os"

	"log/slog"

	"github.com/FreitasGabriel/fullcycle-api/configs"
	"github.com/FreitasGabriel/fullcycle-api/internal/entity"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/database"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/webserver/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Stargin application")

	logger.Info("Loading config")
	config, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logger.Info("Connecting to database")
	db, err := gorm.Open(sqlite.Open("teste.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	logger.Info("Running migrations")
	err = db.AutoMigrate(&entity.User{}, &entity.Product{})
	if err != nil {
		panic(err)
	}

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)
	productHandler := handler.NewProductHandler(productDB)
	userHandler := handler.NewUserHandler(userDB)

	logger.Info("Starting server")
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Use(LogRequest)
	r.Use(middleware.WithValue("jwt", config.TokenAuthKey))
	r.Use(middleware.WithValue("jwtExpiresIn", config.JWTExpiresIn))

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuthKey))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/", productHandler.GetAllProducts)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
	})

	r.Route(("/user"), func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Post("/generate_token", userHandler.GetJWT)
	})

	http.ListenAndServe(":8000", r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
