package main

import (
	"net/http"

	"github.com/FreitasGabriel/fullcycle-api/configs"
	"github.com/FreitasGabriel/fullcycle-api/internal/entity"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/database"
	"github.com/FreitasGabriel/fullcycle-api/internal/infra/webserver/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	_, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("teste.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Product{})
	if err != nil {
		panic(err)
	}

	productDB := database.NewProduct(db)
	userDB := database.NewUser(db)
	productHandler := handler.NewProductHandler(productDB)
	userHandler := handler.NewUserHandler(userDB)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/products", productHandler.CreateProduct)
	r.Get("/products/{id}", productHandler.GetProduct)
	r.Get("/products", productHandler.GetAllProducts)
	r.Put("/products/{id}", productHandler.UpdateProduct)
	r.Delete("/products/{id}", productHandler.DeleteProduct)

	r.Post("/user", userHandler.CreateUser)
	// r.Get("/users", userHandler.GetAllUsers)

	http.ListenAndServe(":8000", r)
}
