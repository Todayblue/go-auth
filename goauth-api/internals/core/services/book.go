package services

import (
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"
)

type BookService struct {
	repo ports.BookRepository
}

func NewBookService(repo ports.BookRepository) *BookService {
	return &BookService{
		repo: repo,
	}
}

func (b *BookService) GetBooks() ([]*domain.Book, error) {
	return b.repo.GetBooks()
}

func (b *BookService) CreateBook(title string) (*domain.Book, error) {
	return b.repo.CreateBook(title)
}
