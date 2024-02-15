package main

import (
	"context"
	"errors"
	"github.com/andre2ar/go-products/configs"
	_ "github.com/andre2ar/go-products/docs"
	"github.com/andre2ar/go-products/internal/entity"
	"github.com/andre2ar/go-products/internal/infra/database"
	"github.com/andre2ar/go-products/internal/infra/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/swaggo/http-swagger/v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title           Go Products
// @version         1.0
// @description     Product API with authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   André Alvim Ribeiro
// @contact.url    https://www.linkedin.com/in/andre2ar/
// @contact.email  andre2ar@outlook.com

// @host      localhost:8000
// @BasePath  /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	config, err := configs.LoadConfig("./cmd/server")
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

	userRepository := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userRepository)

	log.Println("Documentation can be found on " + config.DocsUrl + "/api/v1/docs/index.html")

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.WithValue("Jwt", config.TokenAuth))
	r.Use(middleware.WithValue("JwtExpiresIn", config.JWTExpiresIn))

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL(config.DocsUrl+"/api/v1/docs/doc.json")))

		r.Post("/sessions", userHandler.CreateSession)

		r.Post("/users", userHandler.CreateUser)

		r.Route("/products", func(r chi.Router) {
			r.Use(jwtauth.Verifier(config.TokenAuth))
			r.Use(jwtauth.Authenticator(config.TokenAuth))

			r.Get("/", productHandler.GetProducts)
			r.Post("/", productHandler.CreateProduct)
			r.Get("/{id}", productHandler.GetProduct)
			r.Put("/{id}", productHandler.UpdateProduct)
			r.Delete("/{id}", productHandler.DeleteProduct)
		})
	})

	startServer(r)
}

func startServer(r *chi.Mux) {
	server := &http.Server{
		Addr:    ":8000",
		Handler: r,
	}

	go func() {
		err := server.ListenAndServe()

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()

	log.Println("Server started at: http://localhost:8000")

	WaitForTerminateSignal()

	GracefullyShutdown(server)
}

func WaitForTerminateSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop
}

func GracefullyShutdown(server *http.Server) {
	ctx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	log.Println("Shutting down server in upt to 5 seconds...")
	err := server.Shutdown(ctx)
	if err != nil {
		log.Fatalf("Could not gracefully shutdown server: %v\n", err)
	}

	log.Println("Server stopped")
}
