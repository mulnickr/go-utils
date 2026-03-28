package main

import (
	"net/http"
	"slices"

	"github.com/mulnickr/go-utils/serve"
)

func main() {
	// Initialize router
	router := serve.New()

	// Use middleware
	router.Use(DefaultAuth)

	// Group routes
	api := router.Group("/api/v1")

	// Create endpoints
	api.GET("/health", func(c *serve.Context) {
		c.JSON(http.StatusOK, serve.J{"health": "alive"})
	})

	// "/api/v1/books"
	books := api.Group("/books")
	books.GET("", getBooks)
	books.GET("/:id", getBook)
	books.POST("", createBook)
	books.PUT("/:id", updateBook)
	books.DELETE("/:id", deleteBook)
	// etc.

	// Nest groups
	auth := api.Group("/auth")
	// "/api/v1/auth/login"
	auth.POST("/login", login)

	// Start server
	router.ListenAndServe(":5000")
}

// DefaultAuth is a simple authentication middleware
func DefaultAuth(next serve.Handler) serve.HandlerFunc {
	// Authentication logic
	auth := true

	return serve.HandlerFunc(func(c *serve.Context) {
		if !auth {
			c.JSON(http.StatusUnauthorized, serve.J{"error": "not authorized"})
			// End the handler chain
			return
		}
		// Continue the handler chain
		next.ServeHTTP(c)
	})
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func login(ctx *serve.Context) {
	var body *LoginRequest
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, serve.J{"error": "malformed body", "msg": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, serve.J{"success": "login successful", "user": body.Username})
}

type Book struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var books = []*Book{
	{
		ID:     "1",
		Title:  "Where the Red Fern Grows",
		Author: "Wilson Rawls",
	},
	{
		ID:     "2",
		Title:  "1984",
		Author: "George Orwell",
	},
	{
		ID:     "3",
		Title:  "Animal Farm",
		Author: "George Orwell",
	},
}

func getBooks(ctx *serve.Context) {
	ctx.JSON(http.StatusOK, serve.J{"books": books})
}

func getBook(ctx *serve.Context) {
	id := ctx.Param("id")

	book := findBook(id)
	if book == nil {
		ctx.JSON(http.StatusBadRequest, serve.J{"error": "book not found"})
		return
	}

	ctx.JSON(http.StatusOK, serve.J{"result": book})
}

func findBook(id string) *Book {
	for _, book := range books {
		if book.ID == id {
			return book
		}
	}

	return nil
}

func createBook(ctx *serve.Context) {
	var body *Book
	if err := ctx.Bind(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, serve.J{"error": "malformed entity", "msg": err.Error()})
		return
	}

	books = append(books, body)
	ctx.JSON(http.StatusCreated, serve.J{"success": "book created"})
}

func updateBook(ctx *serve.Context) {
	var body Book
	if err := ctx.Bind(body); err != nil {
		ctx.JSON(http.StatusBadRequest, serve.J{"error": "malformed entity"})
		return
	}

	ctx.JSON(http.StatusNotImplemented, serve.J{"error": "not implemented"})
}

func deleteBook(ctx *serve.Context) {
	id := ctx.Param("id")

	for i, book := range books {
		if book.ID == id {
			books = slices.Delete(books, i, i+1)
			ctx.JSON(http.StatusOK, serve.J{"success": "book deleted"})
			return
		}
	}

	ctx.JSON(http.StatusBadRequest, serve.J{"error": "book now found"})
}
