package handlers

import (
	"book-management-api/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	DB *sql.DB
}

func NewBookHandler(db *sql.DB) *BookHandler {
	return &BookHandler{DB: db}
}

func (h *BookHandler) GetAll(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
			   b.price, b.total_page, b.thickness, b.category_id,
			   b.created_at, b.created_by, b.modified_at, b.modified_by,
			   c.name as category_name
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		ORDER BY b.id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch books",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		var categoryName sql.NullString
		
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Description,
			&book.ImageURL,
			&book.ReleaseYear,
			&book.Price,
			&book.TotalPage,
			&book.Thickness,
			&book.CategoryID,
			&book.CreatedAt,
			&book.CreatedBy,
			&book.ModifiedAt,
			&book.ModifiedBy,
			&categoryName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan book",
				Error:   err.Error(),
			})
			return
		}
		
		if categoryName.Valid {
			book.CategoryName = categoryName.String
		}
		
		books = append(books, book)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}

func (h *BookHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   err.Error(),
		})
		return
	}

	var book models.Book
	var categoryName sql.NullString
	
	err = h.DB.QueryRow(`
		SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
			   b.price, b.total_page, b.thickness, b.category_id,
			   b.created_at, b.created_by, b.modified_at, b.modified_by,
			   c.name as category_name
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.id = $1
	`, id).Scan(
		&book.ID,
		&book.Title,
		&book.Description,
		&book.ImageURL,
		&book.ReleaseYear,
		&book.Price,
		&book.TotalPage,
		&book.Thickness,
		&book.CategoryID,
		&book.CreatedAt,
		&book.CreatedBy,
		&book.ModifiedAt,
		&book.ModifiedBy,
		&categoryName,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   "book with specified ID does not exist",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch book",
			Error:   err.Error(),
		})
		return
	}

	if categoryName.Valid {
		book.CategoryName = categoryName.String
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Book retrieved successfully",
		Data:    book,
	})
}

func (h *BookHandler) Create(c *gin.Context) {
	var bookInput models.BookInput
	if err := c.ShouldBindJSON(&bookInput); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// Validate category exists if provided
	if bookInput.CategoryID != nil {
		var categoryExists bool
		err := h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", *bookInput.CategoryID).Scan(&categoryExists)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to validate category",
				Error:   err.Error(),
			})
			return
		}

		if !categoryExists {
			c.JSON(http.StatusBadRequest, models.APIResponse{
				Success: false,
				Message: "Invalid category ID",
				Error:   "category with specified ID does not exist",
			})
			return
		}
	}

	// Determine thickness based on total_page
	thickness := "tipis"
	if bookInput.TotalPage > 100 {
		thickness = "tebal"
	}

	username, exists := c.Get("username")
	if !exists {
		username = "system"
	}

	var book models.Book
	err := h.DB.QueryRow(`
		INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_by, modified_by) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
		RETURNING id, created_at, modified_at
	`, bookInput.Title, bookInput.Description, bookInput.ImageURL, bookInput.ReleaseYear, 
	   bookInput.Price, bookInput.TotalPage, thickness, bookInput.CategoryID, username, username).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.ModifiedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create book",
			Error:   err.Error(),
		})
		return
	}

	// Populate the book struct
	book.Title = bookInput.Title
	book.Description = bookInput.Description
	book.ImageURL = bookInput.ImageURL
	book.ReleaseYear = bookInput.ReleaseYear
	book.Price = bookInput.Price
	book.TotalPage = bookInput.TotalPage
	book.Thickness = thickness
	book.CategoryID = bookInput.CategoryID
	usernameStr := username.(string)
	book.CreatedBy = &usernameStr
	book.ModifiedBy = &usernameStr

	// Get category name if category_id is provided
	if bookInput.CategoryID != nil {
		var categoryName string
		err := h.DB.QueryRow("SELECT name FROM categories WHERE id = $1", *bookInput.CategoryID).Scan(&categoryName)
		if err == nil {
			book.CategoryName = categoryName
		}
	}

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Book created successfully",
		Data:    book,
	})
}

func (h *BookHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid book ID",
			Error:   err.Error(),
		})
		return
	}

	// Check if book exists
	var exists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM books WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check book existence",
			Error:   err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   "book with specified ID does not exist",
		})
		return
	}

	result, err := h.DB.Exec("DELETE FROM books WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete book",
			Error:   err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Book not found",
			Error:   "book with specified ID does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Book deleted successfully",
	})
}