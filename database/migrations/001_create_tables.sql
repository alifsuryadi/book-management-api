-- +migrate Up

-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(255)
);

-- Create categories table
CREATE TABLE IF NOT EXISTS categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(255)
);

-- Create books table
CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    image_url VARCHAR(255),
    release_year INTEGER,
    price INTEGER,
    total_page INTEGER,
    thickness VARCHAR(50),
    category_id INTEGER REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by VARCHAR(255),
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_by VARCHAR(255)
);

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_books_category_id ON books(category_id);
CREATE INDEX IF NOT EXISTS idx_books_release_year ON books(release_year);
CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);

-- Insert default admin user (password: admin123)
INSERT INTO users (username, password, created_by) VALUES 
('admin', '$2a$10$6v4IMCcL8NZT/JSF9fSwr.fQx/zWzsBGJF8N9m3bj5RRS4HSPpjUy', 'system');

-- Insert sample categories
INSERT INTO categories (name, created_by) VALUES 
('Fiction', 'system'),
('Non-Fiction', 'system'),
('Science', 'system'),
('Technology', 'system'),
('History', 'system');

-- +migrate Down

DROP TABLE IF EXISTS books;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS users;