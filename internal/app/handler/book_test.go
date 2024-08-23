package handler

import (
	"bytes"
	"errors"
	"http-rest-api-go/internal/app/model"
	"http-rest-api-go/internal/app/service"
	mock_service "http-rest-api-go/internal/app/service/mocks"
	"strings"

	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_handleBooksCreate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockBookItem, book *model.Book)

	tests := []struct {
		name                 string
		inputBody            string
		inputBook            *model.Book
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "title", "author": "author"}`,
			inputBook: &model.Book{Title: "title", Author: "author"},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.Book) {
				r.EXPECT().Create(book).Return(nil)
			},
			expectedStatusCode:   201,
			expectedResponseBody: `{"id":0,"title":"title","author":"author"}`,
		},
		{
			name:      "Wrong Input",
			inputBody: `{"title": "title"}`,
			inputBook: &model.Book{Title: "title"},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.Book) {
				r.EXPECT().Create(book).Return(errors.New("author: cannot be blank."))
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"error":"author: cannot be blank."}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"title": "title", "author": "author"}`,
			inputBook: &model.Book{
				Title:  "title",
				Author: "author",
			},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.Book) {
				r.EXPECT().Create(book).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookItem(c)
			test.mockBehavior(repo, test.inputBook)

			service := &service.Service{BookItem: repo}
			handler := Handler{service}

			// Init Endpoint
			r := mux.NewRouter()

			r.HandleFunc("/books", handler.handleBooksCreate()).Methods("POST")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/books",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Trim(w.Body.String(), "\n"), test.expectedResponseBody)
		})
	}
}

func TestHandler_handleBooksGetAll(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockBookItem)

	tests := []struct {
		name                 string
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "Ok",
			mockBehavior: func(r *mock_service.MockBookItem) {
				r.EXPECT().GetAll().Return([]*model.Book{{Title: "title", Author: "author"}}, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `[{"id":0,"title":"title","author":"author"}]`,
		},
		{
			name: "Service Error",
			mockBehavior: func(r *mock_service.MockBookItem) {
				r.EXPECT().GetAll().Return([]*model.Book{}, errors.New("something went wrong"))
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookItem(c)
			test.mockBehavior(repo)

			service := &service.Service{BookItem: repo}
			handler := Handler{service}

			// Init Endpoint
			r := mux.NewRouter()

			r.HandleFunc("/books", handler.handleBooksGetAll()).Methods("GET")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/books",
				bytes.NewBufferString(""))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Trim(w.Body.String(), "\n"), test.expectedResponseBody)
		})
	}
}

func TestHandler_handleBooksUpdate(t *testing.T) {
	// Init Test Table
	type mockBehavior func(r *mock_service.MockBookItem, book *model.UpdateBookInput)

	title := "title"
	author := "author"

	tests := []struct {
		name                 string
		inputBody            string
		inputBook            *model.UpdateBookInput
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "Ok",
			inputBody: `{"title": "title", "author": "author"}`,
			inputBook: &model.UpdateBookInput{Title: &title, Author: &author},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.UpdateBookInput) {
				r.EXPECT().Update(1, book).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:      "Only title",
			inputBody: `{"title": "title"}`,
			inputBook: &model.UpdateBookInput{Title: &title},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.UpdateBookInput) {
				r.EXPECT().Update(1, book).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:      "Only author",
			inputBody: `{"author": "author"}`,
			inputBook: &model.UpdateBookInput{Author: &author},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.UpdateBookInput) {
				r.EXPECT().Update(1, book).Return(nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: ``,
		},
		{
			name:      "empty input",
			inputBody: `{}`,
			inputBook: &model.UpdateBookInput{},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.UpdateBookInput) {
				r.EXPECT().Update(1, book).Return(errors.New("empty input"))
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"error":"empty input"}`,
		},
		{
			name:      "Service Error",
			inputBody: `{"title": "title", "author": "author"}`,
			inputBook: &model.UpdateBookInput{
				Title:  &title,
				Author: &author,
			},
			mockBehavior: func(r *mock_service.MockBookItem, book *model.UpdateBookInput) {
				r.EXPECT().Update(1, book).Return(errors.New("something went wrong"))
			},
			expectedStatusCode:   422,
			expectedResponseBody: `{"error":"something went wrong"}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Init Dependencies
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockBookItem(c)
			test.mockBehavior(repo, test.inputBook)

			service := &service.Service{BookItem: repo}
			handler := Handler{service}

			// Init Endpoint
			r := mux.NewRouter()

			r.HandleFunc("/books/{id}", handler.handleBooksPut()).Methods("PUT")

			// Create Request
			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/books/1",
				bytes.NewBufferString(test.inputBody))

			// Make Request
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, w.Code, test.expectedStatusCode)
			assert.Equal(t, strings.Trim(w.Body.String(), "\n"), test.expectedResponseBody)
		})
	}
}
