package store

// Store ...
type Store interface {
	Book() BookRepository
}
