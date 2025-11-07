-- Add indexes for better query performance
-- Created at: 2024-10-25

-- Users table indexes
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_is_active ON users(is_active);
CREATE INDEX IF NOT EXISTS idx_users_deleted_at ON users(deleted_at);

-- Books table indexes
CREATE INDEX IF NOT EXISTS idx_books_isbn ON books(isbn);
CREATE INDEX IF NOT EXISTS idx_books_author ON books(author);
CREATE INDEX IF NOT EXISTS idx_books_title ON books(title);
CREATE INDEX IF NOT EXISTS idx_books_year ON books(year);
CREATE INDEX IF NOT EXISTS idx_books_deleted_at ON books(deleted_at);

-- Composite indexes for common queries
CREATE INDEX IF NOT EXISTS idx_users_email_is_active ON users(email, is_active);
CREATE INDEX IF NOT EXISTS idx_books_author_year ON books(author, year);
