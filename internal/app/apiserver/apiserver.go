package apiserver

import (
	"database/sql"
	"log/slog"
	"net/http"

	"http-rest-api-go/internal/app/handler"
	"http-rest-api-go/internal/app/service"
	"http-rest-api-go/internal/app/store/sqlstore"
	"http-rest-api-go/internal/config"

	_ "github.com/lib/pq" // ...
)

// Start ...
func Start(config *config.Config, logger *slog.Logger) error {
	db, err := newDB(config.StoragePath)
	if err != nil {
		return err
	}

	defer db.Close()
	store, err := sqlstore.New(db)

	if err != nil {
		return err
	}

	services := service.NewService(store.Book())
	handlers := handler.NewHandler(services)

	srv := newServer(handlers.InitRoutes(), logger)

	return http.ListenAndServe(config.HTTPServer.Address, srv)
}

func newDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
