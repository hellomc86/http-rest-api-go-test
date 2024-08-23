package sqlstore

import (
	"database/sql"
	"errors"

	"http-rest-api-go/internal/app/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestBook_Repository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := &Store{db: db}

	type mockBehavior func(book *model.Book, id int)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   *model.Book
		want    int
		wantErr bool
	}{
		{
			name: "Ok",
			input: &model.Book{
				Title:  "title",
				Author: "author",
			},
			want: 1,
			mock: func(book *model.Book, id int) {
				mock.ExpectBegin()
				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO books").
					WithArgs(book.Title, book.Author).WillReturnRows(rows)
				mock.ExpectCommit()
			},
		},
		{
			name: "Failed 2nd Insert",
			input: &model.Book{
				Title:  "title",
				Author: "author",
			},
			mock: func(book *model.Book, id int) {
				mock.ExpectBegin()

				sqlmock.NewRows([]string{"id"}).AddRow(id)

				mock.ExpectQuery("INSERT INTO books").
					WithArgs(book.Title, book.Author).WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			err := r.Book().Create(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBook_Repository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := &Store{db: db}

	type mockBehavior func(book *model.Book, id int)
	tests := []struct {
		name    string
		mock    func()
		want    []*model.Book
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "title1", "author1").
					AddRow(2, "title2", "author2").
					AddRow(3, "title3", "author3")

				mock.ExpectQuery("SELECT (.+) FROM books").WillReturnRows(rows)
			},
			want: []*model.Book{
				{1, "title1", "author1"},
				{2, "title2", "author2"},
				{3, "title3", "author3"},
			},
		},
		{
			name: "No Records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author"})

				mock.ExpectQuery("SELECT (.+) FROM books").WillReturnRows(rows)
			},
			want: []*model.Book{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Book().FindAll()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBook_Repository_GetById(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := &Store{db: db}

	tests := []struct {
		name    string
		mock    func()
		id      int
		want    *model.Book
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author"}).
					AddRow(1, "title1", "author1")

				mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).WillReturnRows(rows)
			},
			want: &model.Book{
				1, "title1", "author1",
			},
			id: 1,
		},
		{
			name: "Not Found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "author"})
				mock.ExpectQuery("SELECT (.+) FROM books").WithArgs(1).WillReturnRows(rows)
			},
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			got, err := r.Book().Find(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBook_Repository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := &Store{db: db}

	tests := []struct {
		name    string
		mock    func()
		id      int
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {

				mock.ExpectExec("Delete FROM books WHERE (.+)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
			},
			id:      1,
			wantErr: false,
		},
		{
			name: "Not Found",
			mock: func() {
				mock.ExpectExec("Delete FROM books WHERE (.+)").
					WithArgs(1).WillReturnError(sql.ErrNoRows)
			},
			id:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Book().Delete(tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBook_Repository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := &Store{db: db}

	type args struct {
		id    int
		input *model.UpdateBookInput
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		wantErr bool
	}{
		{
			name: "OK_AllFields",
			mock: func() {
				mock.ExpectExec("UPDATE books SET (.+) WHERE (.+)").
					WithArgs("new title", "new author", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				id: 1,
				input: &model.UpdateBookInput{
					Title:  stringPointer("new title"),
					Author: stringPointer("new author"),
				},
			},
		},
		{
			name: "OK_WithoutAuthor",
			mock: func() {
				mock.ExpectExec("UPDATE books SET (.+) WHERE (.+)").
					WithArgs("new title", 1).WillReturnResult(sqlmock.NewResult(0, 1))
			},
			input: args{
				id: 1,
				input: &model.UpdateBookInput{
					Title: stringPointer("new title"),
				},
			},
		},
		{
			name: "OK_NoInputFields",
			mock: func() {

			},
			input: args{
				id:    1,
				input: &model.UpdateBookInput{},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()

			err := r.Book().Update(tt.input.id, tt.input.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func stringPointer(s string) *string {
	return &s
}
