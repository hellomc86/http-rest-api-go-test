package service

import (
	"http-rest-api-go/internal/app/model"
	"http-rest-api-go/internal/app/store"
)

type BookService struct {
	repo store.BookRepository
}

func NewBookService(repo store.BookRepository) *BookService {
	return &BookService{repo: repo}
}

func (s *BookService) Create(book *model.Book) error {

	return s.repo.Create(book)
}

func (s *BookService) GetAll() ([]*model.Book, error) {
	return s.repo.FindAll()
}

func (s *BookService) GetById(Id int) (*model.Book, error) {
	return s.repo.Find(Id)
}

func (s *BookService) Delete(Id int) error {
	return s.repo.Delete(Id)
}

func (s *BookService) Update(Id int, input *model.UpdateBookInput) error {
	return s.repo.Update(Id, input)
}
