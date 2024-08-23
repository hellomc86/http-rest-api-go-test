package handler

import (
	"http-rest-api-go/internal/app/service"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Handler struct {
	service *service.Service
}



func NewHandler(services *service.Service) *Handler {
	return &Handler{service: services}
}

func (h *Handler) InitRoutes() *mux.Router {
	router := mux.NewRouter()

	router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	router.HandleFunc("/books", h.handleBooksCreate()).Methods("POST")
	router.HandleFunc("/books/", h.handleBooksGetAll()).Methods("GET")
	router.HandleFunc("/books/{id}", h.handleBooksGet()).Methods("GET")
	router.HandleFunc("/books/{id}", h.handleBooksPut()).Methods("PUT")
	router.HandleFunc("/books/{id}", h.handleBooksDelete()).Methods("Delete")
	return router
}
