package sqlstore

import (
	"database/sql"

	"http-rest-api-go/internal/app/store"
)

// Store ...
type Store struct {
	db             *sql.DB
	bookRepository *BookRepository
}

// New ...
func New(db *sql.DB) (*Store, error) {

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS books(
		id bigserial PRIMARY KEY,
		title TEXT NOT NULL UNIQUE,
		author TEXT NOT NULL);
		`)

	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

// Book ...
func (s *Store) Book() store.BookRepository {
	if s.bookRepository != nil {
		return s.bookRepository
	}

	s.bookRepository = &BookRepository{
		store: s,
	}

	return s.bookRepository
}
