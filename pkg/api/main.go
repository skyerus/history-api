package api

import (
	"database/sql"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // mysql driver
	"github.com/skyerus/history-api/pkg/logger"

	"github.com/gorilla/mux"
)

// App - Root component
type App struct {
	Router *mux.Router
}

type router struct {
	db *sql.DB
}

func newRouter(db *sql.DB) *router {
	return &router{db}
}

// Initialize - Initialize app
func (a *App) Initialize(db *sql.DB) {
	router := newRouter(db)
	a.Router = mux.NewRouter()
	a.Router.Use(cors)
	a.setRouters(router)
}

func (a *App) setRouters(router *router) {
	// Base routes
	a.Router.HandleFunc("/", healthCheck).Methods("GET", "OPTIONS")
	a.Router.HandleFunc("/api/random_article", router.randomArticle).Methods("GET", "OPTIONS")
}

// Run - Run the app
func (a *App) Run(host string) {
	srv := &http.Server{
		Handler:      a.Router,
		Addr:         host,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  18 * time.Second,
	}
	go logger.Log(srv.ListenAndServe())
}

// OpenDb ... OpenDb connection
func OpenDb() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DB_URL")+"?parseTime=true")
	if err != nil {
		go logger.Log(err)
		return db, err
	}
	return db, nil
}
