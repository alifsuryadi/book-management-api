# Book Management API

This is a RESTful API for managing books and their categories, built with Go using Gin framework and PostgreSQL database.

## Features

- 📚 Book management (CRUD operations)
- 📂 Category management (CRUD operations)
- 🔐 JWT Authentication
- 🗄️ PostgreSQL database
- 🚀 Railway deployment ready
- ✅ Input validation
- 📖 Comprehensive documentation

## Tech Stack

- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **Authentication**: JWT (JSON Web Token)
- **Migration**: sql-migrate
- **Deployment**: Railway

## Project Structure

```
book-management-api/
├── config/
│   └── config.go          # Configuration management
├── database/
│   ├── database.go        # Database connection
│   └── migrations/        # Database migration files
├── handlers/
│   ├── users.go          # User authentication handlers
│   ├── categories.go     # Category handlers
│   └── books.go          # Book handlers
├── middleware/
│   └── auth.go           # Authentication middleware
├── models/
│   └── models.go         # Data models and structs
├── routes/
│   └── routes.go         # API routes configuration
├── go.mod
├── go.sum
├── main.go
└── README.md
```

## Installation & Setup

### Prerequisites

- Go 1.21 or higher
- PostgreSQL database
- Git

### Local Development

1. **Clone the repository**

   ```bash
   git clone https://github.com/alifsuryadi/book-management-api
   cd book-management-api
   ```

2. **Install dependencies**

   ```bash
   go mod download
   ```

3. **Set up PostgreSQL Database**

   Make sure PostgreSQL is installed and running on your system. Then create a database:

   ```bash
   # Connect to PostgreSQL as superuser
   psql -U postgres

   # Create database and user (in PostgreSQL prompt)
   CREATE DATABASE book_management;
   CREATE USER bookuser WITH PASSWORD 'bookpass123';
   GRANT ALL PRIVILEGES ON DATABASE book_management TO bookuser;
   \q
   ```

4. **Set up environment variables**

   Copy the example environment file:

   ```bash
   cp .env.example .env
   ```

   Then edit `.env` file with your database credentials:

   ```bash
   DATABASE_URL=postgres://postgres:your_password@localhost:5432/book_management?sslmode=disable
   JWT_SECRET="your-secret-key-here"
   ENVIRONMENT=development
   ```

   Or export them directly:

   ```bash
   export DATABASE_URL="postgres://bookuser:bookpass123@localhost:5432/book_management?sslmode=disable"
   export JWT_SECRET="your-secret-key-here"
   export ENVIRONMENT="development"
   ```

5. **Run the application**
   ```bash
   go run main.go
   ```

The server will start on port 8080 (or the PORT environment variable if set).

## Database Schema

### Users Table

- `id` (integer, primary key)
- `username` (varchar, unique)
- `password` (varchar, hashed)
- `created_at` (timestamp)
- `created_by` (varchar)
- `modified_at` (timestamp)
- `modified_by` (varchar)

### Categories Table

- `id` (integer, primary key)
- `name` (varchar)
- `created_at` (timestamp)
- `created_by` (varchar)
- `modified_at` (timestamp)
- `modified_by` (varchar)

### Books Table

- `id` (integer, primary key)
- `title` (varchar)
- `description` (text)
- `image_url` (varchar)
- `release_year` (integer, min: 1980, max: 2024)
- `price` (integer)
- `total_page` (integer)
- `thickness` (varchar, auto-calculated: "tebal" if >100 pages, "tipis" if d100 pages)
- `category_id` (integer, foreign key)
- `created_at` (timestamp)
- `created_by` (varchar)
- `modified_at` (timestamp)
- `modified_by` (varchar)

## API Endpoints

### Authentication

#### Login

- **POST** `/api/users/login`
- **Description**: Authenticate user and get JWT token
- **Request Body**:
  ```json
  {
    "username": "admin",
    "password": "admin123"
  }
  ```
- **Response**:
  ```json
  {
    "success": true,
    "message": "Login successful",
    "data": {
      "token": "jwt_token_here",
      "user": {
        "id": 1,
        "username": "admin",
        "created_at": "2024-01-01T00:00:00Z"
      }
    }
  }
  ```

### Categories

All category endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### Get All Categories

- **GET** `/api/categories`
- **Description**: Retrieve all categories

#### Create Category

- **POST** `/api/categories`
- **Request Body**:
  ```json
  {
    "name": "Science Fiction"
  }
  ```

#### Get Category by ID

- **GET** `/api/categories/:id`
- **Description**: Retrieve specific category

#### Delete Category

- **DELETE** `/api/categories/:id`
- **Description**: Delete specific category

#### Get Books by Category

- **GET** `/api/categories/:id/books`
- **Description**: Retrieve all books in a specific category

### Books

All book endpoints require JWT authentication via `Authorization: Bearer <token>` header.

#### Get All Books

- **GET** `/api/books`
- **Description**: Retrieve all books with category information

#### Create Book

- **POST** `/api/books`
- **Request Body**:
  ```json
  {
    "title": "The Great Gatsby",
    "description": "A classic American novel",
    "image_url": "https://example.com/image.jpg",
    "release_year": 1925,
    "price": 15000,
    "total_page": 180,
    "category_id": 1
  }
  ```
- **Note**: The `thickness` field is automatically calculated:
  - "tebal" if `total_page` > 100
  - "tipis" if `total_page` d 100

#### Get Book by ID

- **GET** `/api/books/:id`
- **Description**: Retrieve specific book with category information

#### Delete Book

- **DELETE** `/api/books/:id`
- **Description**: Delete specific book

### Health Check

- **GET** `/health`
- **Description**: Check API health status

## Validation Rules

### Books

- `title`: Required
- `release_year`: Required, must be between 1980 and 2024
- `price`: Required, must be positive integer
- `total_page`: Required, must be positive integer
- `category_id`: Optional, must exist in categories table if provided

### Categories

- `name`: Required

## Error Handling

The API returns consistent error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error information"
}
```

Common HTTP status codes:

- `200`: Success
- `201`: Created
- `400`: Bad Request (validation errors)
- `401`: Unauthorized (authentication required)
- `404`: Not Found
- `500`: Internal Server Error

## Authentication

This API uses JWT (JSON Web Token) for authentication. After logging in, include the token in the Authorization header:

```
Authorization: Bearer your_jwt_token_here
```

### Default User

- Username: `admin`
- Password: `admin123`

### Admin Seeding

To create the default admin user, use the `/api/seed/admin` endpoint:

```bash
# Seed admin user
curl -X POST http://localhost:8080/api/seed/admin \
  -H "Content-Type: application/json"
```

This will create an admin user with:
- Username: `admin`
- Password: `admin123`

**Note**: This endpoint should only be used once during initial setup. It will return an error if an admin user already exists.

### Password Reset

To reset a user's password, use the `/api/reset-password` endpoint:

```bash
# Reset password
curl -X POST http://localhost:8080/api/reset-password \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "new_password": "newpassword123"
  }'
```

**Security Note**: This endpoint should be protected in production environments with additional authentication or administrative privileges.

## Deployment to Railway

1. **Create Railway account** at [railway.app](https://railway.app)

2. **Create new project** and add PostgreSQL database

3. **Set environment variables**:

   - `DATABASE_URL`: (automatically set by Railway PostgreSQL)
   - `JWT_SECRET`: Your secret key
   - `ENVIRONMENT`: production

4. **Deploy**:

   ```bash
   # Connect to Railway
   railway login
   railway link

   # Deploy
   railway up
   ```

5. **Your API will be available at**: `https://your-app-name.railway.app`

## Example Usage

### 1. Login to get token

```bash
curl -X POST https://your-api-url.railway.app/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

### 2. Create a category

```bash
curl -X POST https://your-api-url.railway.app/api/categories \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Programming"
  }'
```

### 3. Create a book

```bash
curl -X POST https://your-api-url.railway.app/api/books \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "title": "Clean Code",
    "description": "A Handbook of Agile Software Craftsmanship",
    "release_year": 2008,
    "price": 45000,
    "total_page": 464,
    "category_id": 1
  }'
```

### 4. Get all books

```bash
curl -X GET https://your-api-url.railway.app/api/books \
  -H "Authorization: Bearer YOUR_TOKEN"
```

## Environment Variables

| Variable       | Description                  | Default                                                                       |
| -------------- | ---------------------------- | ----------------------------------------------------------------------------- |
| `DATABASE_URL` | PostgreSQL connection string | `postgres://postgres:password@localhost:5432/book_management?sslmode=disable` |
| `JWT_SECRET`   | Secret key for JWT tokens    | `your-secret-key-change-this-in-production`                                   |
| `ENVIRONMENT`  | Application environment      | `development`                                                                 |
| `PORT`         | Server port                  | `8080`                                                                        |

## Development

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Building for Production

```bash
go build -o book-management-api
./book-management-api
```

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## License

This project is created for educational purposes as part of Golang Bootcamp Quiz.

## Support

For questions or issues, please create an issue in the repository or contact the development team.
