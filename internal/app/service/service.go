package service

import (
	"http-rest-api-go/internal/app/model"
	"http-rest-api-go/internal/app/store"
)

//go:generate mockgen -source=service.go -destination=mocks/mock.go

type BookItem interface {
	Create(book *model.Book) error
	GetAll() ([]*model.Book, error)
	GetById(Id int) (*model.Book, error)
	Delete(Id int) error
	Update(Id int, input *model.UpdateBookInput) error
}

type Service struct {
	BookItem
}

func NewService(repos store.BookRepository) *Service {
	return &Service{
		BookItem: NewBookService(repos),
	}
}
