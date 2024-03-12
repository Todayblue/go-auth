package handler

import (
	"go-chat/internals/core/domain"
	"go-chat/internals/core/ports"

	"github.com/gofiber/fiber/v2"
)

type BookHandler struct {
	bookService ports.BookService
}

func NewBookHandlers(bookService ports.BookService) *BookHandler {
	return &BookHandler{
		bookService: bookService,
	}
}

func (h *BookHandler) GetBooks(c *fiber.Ctx) error {
	books, err := h.bookService.GetBooks()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   books,
	})
}

func (h *BookHandler) CreateBook(c *fiber.Ctx) error {
	book := domain.Book{}
	// var req string
	err := c.BodyParser(&book)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "invalid request body",
		})
	}

	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "fail",
			"message": "title is required",
		})
	}

	result, err := h.bookService.CreateBook(book.Title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   result,
	})
}
