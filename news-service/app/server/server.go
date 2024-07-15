package server

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"news_service/app/config"
	"news_service/app/handlers"
)

type Server struct {
	httpServer *http.Server
}

func NewServer() error {
	log.Println("Starting server on Port:", config.Port)
	return http.ListenAndServe(config.Port, InitRoutes())
}

func InitRoutes() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/sources", handlers.GetSourcesHandler).Methods("GET")
	router.HandleFunc("/sources", handlers.AddSourceHandler).Methods("POST")
	router.HandleFunc("/sources/{id}/news", handlers.GetUserNews).Methods("GET")

	return router
}
