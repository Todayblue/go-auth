package repository

import (
	"go-chat/internals/core/domain"
)

func (b *DB) GetBooks() ([]*domain.Book, error) {
	var books []*domain.Book
	result := b.db.Find(&books)
	if result.Error != nil {
		return nil, result.Error
	}
	return books, nil
}

func (b *DB) CreateBook(title string) (*domain.Book, error) {
	book := &domain.Book{
		Title: title,
	}

	result := b.db.Create(&book)
	if result.Error != nil {
		return nil, result.Error
	}

	return book, nil
}
