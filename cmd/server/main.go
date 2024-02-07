package main

import (
	"github.com/andre2ar/go-products/configs"
	"github.com/andre2ar/go-products/internal/entity"
	"github.com/andre2ar/go-products/internal/infra/database"
	"github.com/andre2ar/go-products/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func main() {
	_, err := configs.LoadConfig("./cmd/server")
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(sqlite.Open("go-products.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("Connected to the database")

	err = db.AutoMigrate(&entity.Product{}, &entity.User{})
	if err != nil {
		panic(err)
	}
	log.Println("Database migrated")

	productRepository := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productRepository)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/products", productHandler.GetProducts)
		r.Post("/products", productHandler.CreateProduct)
		r.Get("/products/{id}", productHandler.GetProduct)
		r.Put("/products/{id}", productHandler.UpdateProduct)
		r.Delete("/products/{id}", productHandler.DeleteProduct)
	})

	log.Fatalln(http.ListenAndServe(":8000", r))
}
