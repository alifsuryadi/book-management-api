package handlers

import (
	"book-management-api/models"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	DB *sql.DB
}

func NewCategoryHandler(db *sql.DB) *CategoryHandler {
	return &CategoryHandler{DB: db}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	rows, err := h.DB.Query(`
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories 
		ORDER BY id ASC
	`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Error:   err.Error(),
		})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var category models.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.CreatedAt,
			&category.CreatedBy,
			&category.ModifiedAt,
			&category.ModifiedBy,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan category",
				Error:   err.Error(),
			})
			return
		}
		categories = append(categories, category)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Categories retrieved successfully",
		Data:    categories,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
		})
		return
	}

	var category models.Category
	err = h.DB.QueryRow(`
		SELECT id, name, created_at, created_by, modified_at, modified_by 
		FROM categories 
		WHERE id = $1
	`, id).Scan(
		&category.ID,
		&category.Name,
		&category.CreatedAt,
		&category.CreatedBy,
		&category.ModifiedAt,
		&category.ModifiedBy,
	)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Category not found",
			Error:   "category with specified ID does not exist",
		})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to fetch category",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category retrieved successfully",
		Data:    category,
	})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	username, exists := c.Get("username")
	if !exists {
		username = "system"
	}

	err := h.DB.QueryRow(`
		INSERT INTO categories (name, created_by, modified_by) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at, modified_at
	`, category.Name, username, username).Scan(
		&category.ID,
		&category.CreatedAt,
		&category.ModifiedAt,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to create category",
			Error:   err.Error(),
		})
		return
	}

	usernameStr := username.(string)
	category.CreatedBy = &usernameStr
	category.ModifiedBy = &usernameStr

	c.JSON(http.StatusCreated, models.APIResponse{
		Success: true,
		Message: "Category created successfully",
		Data:    category,
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
		})
		return
	}

	// Check if category exists
	var exists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check category existence",
			Error:   err.Error(),
		})
		return
	}

	if !exists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Category not found",
			Error:   "category with specified ID does not exist",
		})
		return
	}

	result, err := h.DB.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to delete category",
			Error:   err.Error(),
		})
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Category not found",
			Error:   "category with specified ID does not exist",
		})
		return
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Category deleted successfully",
	})
}

func (h *CategoryHandler) GetBooksByCategory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.APIResponse{
			Success: false,
			Message: "Invalid category ID",
			Error:   err.Error(),
		})
		return
	}

	// Check if category exists
	var categoryExists bool
	err = h.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = $1)", id).Scan(&categoryExists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.APIResponse{
			Success: false,
			Message: "Failed to check category existence",
			Error:   err.Error(),
		})
		return
	}

	if !categoryExists {
		c.JSON(http.StatusNotFound, models.APIResponse{
			Success: false,
			Message: "Category not found",
			Error:   "category with specified ID does not exist",
		})
		return
	}

	rows, err := h.DB.Query(`
		SELECT b.id, b.title, b.description, b.image_url, b.release_year, 
			   b.price, b.total_page, b.thickness, b.category_id,
			   b.created_at, b.created_by, b.modified_at, b.modified_by,
			   c.name as category_name
		FROM books b
		LEFT JOIN categories c ON b.category_id = c.id
		WHERE b.category_id = $1
		ORDER BY b.id ASC
	`, id)
	
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
			&book.CategoryName,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.APIResponse{
				Success: false,
				Message: "Failed to scan book",
				Error:   err.Error(),
			})
			return
		}
		books = append(books, book)
	}

	c.JSON(http.StatusOK, models.APIResponse{
		Success: true,
		Message: "Books retrieved successfully",
		Data:    books,
	})
}