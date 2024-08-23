package store

import "http-rest-api-go/internal/app/model"

// BookRepository ...
type BookRepository interface {
	Create(*model.Book) error
	FindAll() ([]*model.Book, error)
	Find(int) (*model.Book, error)
	FindByName(string) (*model.Book, error)
	Update(int, *model.UpdateBookInput) error
	Delete(int) error
}
