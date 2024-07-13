package server

import (
	"api-service/app/config"
	"api-service/app/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
	router.HandleFunc("/sources", handlers.GetAllSources).Methods("GET")
	router.HandleFunc("/sources", handlers.AddSource).Methods("POST")
	router.HandleFunc("/sources/{id}", handlers.GetUserNews).Methods("GET")
	router.HandleFunc("/sign-up", handlers.CreateUser).Methods("GET")
	router.HandleFunc("/subscribe", handlers.SubscribeUser).Methods("POST")
	router.HandleFunc("/unsubscribe", handlers.UnsubscribeUser).Methods("POST")
	//router.HandleFunc("/sign-up/admin", handlers.CreateAdmin).Methods("GET")

	return router
}
