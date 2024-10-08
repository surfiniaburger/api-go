// api.go
package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/surfiniaburger/api-go/services/cart"
	"github.com/surfiniaburger/api-go/services/library"
	"github.com/surfiniaburger/api-go/services/order"
	"github.com/surfiniaburger/api-go/services/product"
	"github.com/surfiniaburger/api-go/services/user"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/v1").Subrouter()

	userStore := user.NewStore(s.db)
	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoutes(subrouter)

	productStore := product.NewStore(s.db)
	productHandler := product.NewHandler(productStore, userStore)
	productHandler.RegisterRoutes(subrouter)

	orderStore := order.NewStore(s.db)

	bookStore, err := library.NewBookStore(s.db)
	if err != nil {
		log.Fatalf("Failed to create BookStore: %v", err)
	}
	bookHandler := library.NewBookHandler(bookStore, userStore)
	bookHandler.RegisterRoutes(subrouter)

	cartHandler := cart.NewHandler(productStore, orderStore, userStore)
	cartHandler.RegisterRoutes(subrouter)

	// Serve static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("static")))

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
