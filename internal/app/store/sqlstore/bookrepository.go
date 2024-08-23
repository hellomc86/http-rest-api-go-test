package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"

	"http-rest-api-go/internal/app/model"
	"http-rest-api-go/internal/app/store"
)

// BookRepository ...
type BookRepository struct {
	store *Store
}

// Create ...
func (r *BookRepository) Create(b *model.Book) error {
	if err := b.Validate(); err != nil {
		return err
	}

	tx, err := r.store.db.Begin()
	if err != nil {
		return err
	}

	row := tx.QueryRow(
		"INSERT INTO books (title, author) VALUES ($1, $2) RETURNING id",
		b.Title,
		b.Author,
	)
	err = row.Scan(&b.ID)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()
}

// Find ...
func (r *BookRepository) FindAll() ([]*model.Book, error) {
	books := []*model.Book{}

	rows, err := r.store.db.Query("SELECT * FROM books")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	for rows.Next() {
		b := &model.Book{}
		rows.Scan(&b.ID, &b.Title, &b.Author)
		books = append(books, b)
	}

	return books, nil
}

// Find ...
func (r *BookRepository) Find(id int) (*model.Book, error) {
	b := &model.Book{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, author FROM books WHERE id = $1",
		id,
	).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return b, nil
}

// FindByName ...
func (r *BookRepository) FindByName(title string) (*model.Book, error) {
	b := &model.Book{}
	if err := r.store.db.QueryRow(
		"SELECT id, title, author FROM books WHERE title = $1",
		title,
	).Scan(
		&b.ID,
		&b.Title,
		&b.Author,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return b, nil
}

// Update ...
func (r *BookRepository) Update(id int, b *model.UpdateBookInput) error {

	if err := b.Validate(); err != nil {
		return err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if b.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, b.Title)
		argId++
	}

	if b.Author != nil {
		setValues = append(setValues, fmt.Sprintf("author=$%d", argId))
		args = append(args, b.Author)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf(`UPDATE books SET %s WHERE id = $%d`,
		setQuery, argId)
	args = append(args, id)

	_, err := r.store.db.Exec(query, args...)
	return err
}

// Delete ...
func (r *BookRepository) Delete(id int) error {

	if _, err := r.store.db.Exec("Delete FROM books WHERE id=$1", id); err != nil {

		return err
	}

	return nil
}
