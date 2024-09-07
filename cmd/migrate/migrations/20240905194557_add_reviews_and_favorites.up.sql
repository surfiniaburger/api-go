CREATE TABLE IF NOT EXISTS reviews (
    reviewId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    bookId INT NOT NULL,
    rating INT CHECK (rating >= 1 AND rating <= 5), 
    comment TEXT,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (bookId) REFERENCES books(bookId)
);



CREATE TABLE IF NOT EXISTS favorites (
    favoriteId SERIAL PRIMARY KEY,
    userId INT NOT NULL,
    bookId INT NOT NULL,
    createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (userId) REFERENCES users(id),
    FOREIGN KEY (bookId) REFERENCES books(bookId),
    UNIQUE (userId, bookId)
);
