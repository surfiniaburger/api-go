CREATE TABLE IF NOT EXISTS books (
    bookId INT PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,
    category VARCHAR(100) NOT NULL,
    isbn VARCHAR(20),
    publishedDate DATE,
    tags JSON,
    fileUrl VARCHAR(255)
);
