package models

import (
	"time"
)

type User struct {
	ID         int       `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `json:"-" db:"password"` // Don't include in JSON responses
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	ModifiedBy string    `json:"modified_by" db:"modified_by"`
}

type Category struct {
	ID         int       `json:"id" db:"id"`
	Name       string    `json:"name" db:"name" binding:"required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	CreatedBy  string    `json:"created_by" db:"created_by"`
	ModifiedAt time.Time `json:"modified_at" db:"modified_at"`
	ModifiedBy string    `json:"modified_by" db:"modified_by"`
}

type Book struct {
	ID          int       `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
	ImageURL    string    `json:"image_url" db:"image_url"`
	ReleaseYear int       `json:"release_year" db:"release_year" binding:"required,min=1980,max=2024"`
	Price       int       `json:"price" db:"price" binding:"required,min=0"`
	TotalPage   int       `json:"total_page" db:"total_page" binding:"required,min=1"`
	Thickness   string    `json:"thickness" db:"thickness"`
	CategoryID  *int      `json:"category_id" db:"category_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	CreatedBy   string    `json:"created_by" db:"created_by"`
	ModifiedAt  time.Time `json:"modified_at" db:"modified_at"`
	ModifiedBy  string    `json:"modified_by" db:"modified_by"`
	
	// For joined queries
	CategoryName string `json:"category_name,omitempty" db:"category_name"`
}

type BookInput struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year" binding:"required,min=1980,max=2024"`
	Price       int    `json:"price" binding:"required,min=0"`
	TotalPage   int    `json:"total_page" binding:"required,min=1"`
	CategoryID  *int   `json:"category_id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}