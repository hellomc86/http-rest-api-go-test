package model

import (
	"errors"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Book ...
type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

// Validate ...
func (b *Book) Validate() error {
	return validation.ValidateStruct(
		b,
		validation.Field(&b.Title, validation.Required, validation.Length(1, 100)),
		validation.Field(&b.Author, validation.Required, validation.Length(1, 100)),
	)
}

type UpdateBookInput struct {
	Title  *string `json:"title"`
	Author *string `json:"author"`
}

func (i UpdateBookInput) Validate() error {
	if i.Title == nil && i.Author == nil {
		return errors.New("update structure has no values")
	}

	return nil
}
