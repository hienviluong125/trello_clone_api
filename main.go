package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hienviluong125/trello_clone_api/controllers"
	"github.com/hienviluong125/trello_clone_api/database"
	"github.com/hienviluong125/trello_clone_api/middlewares"
	"github.com/hienviluong125/trello_clone_api/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

func initialMigration() {
	var err error

	err = godotenv.Load()

	if err != nil {
		panic(err.Error())
	}

	var (
		dbUser = os.Getenv("DB_USER")
		dbPw   = os.Getenv("DB_PASSWORD")
		dbName = os.Getenv("DB_NAME")
	)

	dbUrl := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPw, dbName)
	database.DBConn, err = gorm.Open("postgres", dbUrl)

	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Database connected")
	// Migrate the schema
	database.DBConn.AutoMigrate(&models.User{})
	fmt.Println("Database Migrated")
}

func setJsonFormatMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// Migrations
	initialMigration()
	// Init routers
	router := mux.NewRouter()
	router.Use(setJsonFormatMiddleware)
	// User routers
	userController := controllers.UserController{}
	router.HandleFunc("/auth/register", userController.Register).Methods("POST")
	router.HandleFunc("/auth/login", userController.Login).Methods("POST")
	router.HandleFunc("/users/profile", middlewares.Authenticate(userController.Profile)).Methods("GET")
	router.HandleFunc("/users/profile", middlewares.Authenticate(userController.UpdateProfile)).Methods("PUT")

	loggedRouter := handlers.LoggingHandler(os.Stdout, router)
	http.ListenAndServe(":8080", loggedRouter)

	defer database.DBConn.Close()
}
