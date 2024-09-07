CREATE TABLE IF NOT EXISTS books (
    bookId SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    isbn VARCHAR(20),
    publishedDate DATE,
    tags JSONB,
    fileUrl VARCHAR(255)
);
