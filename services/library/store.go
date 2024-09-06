package library

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/surfiniaburger/api-go/types"
)

type BookStore struct {
	db *sql.DB
}

func NewBookStore(db *sql.DB) *BookStore {
	return &BookStore{db: db}
}

func (s *BookStore) CreateBook(book types.CreateBookPayload) error {
	tagsJSON, err := json.Marshal(book.Tags)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		INSERT INTO books 
		(title, author, description, category, isbn, publishedDate, tags, fileUrl) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		book.Title, book.Author, book.Description, book.Category, book.ISBN, book.PublishedDate, tagsJSON, book.FileUrl)

	if err != nil {
		return err
	}

	return nil
}

func (s *BookStore) GetBookByID(bookID string) (*types.Book, error) {
	row := s.db.QueryRow("SELECT bookid, title, author, description, category, isbn, publishedDate, tags, fileUrl FROM books WHERE bookid = ?", bookID)

	var book types.Book
	var tags []byte // Store the raw JSON as []byte

	// Scan the row, including the raw JSON for tags
	err := row.Scan(&book.BookID, &book.Title, &book.Author, &book.Description, &book.Category, &book.ISBN, &book.PublishedDate, &tags, &book.FileUrl)
	if err != nil {
		return nil, err
	}

	// Decode the JSON tags into a []string
	if err := json.Unmarshal(tags, &book.Tags); err != nil {
		return nil, err
	}

	return &book, nil
}

func (s *BookStore) UpdateBook(bookID string, book types.UpdateBookPayload) error {
	// Convert the tags slice to JSON
	tagsJSON, err := json.Marshal(book.Tags)
	if err != nil {
		return err
	}

	// Execute the SQL query with the JSON string for tags
	_, err = s.db.Exec(`
        UPDATE books 
        SET title = ?, author = ?, description = ?, category = ?, isbn = ?, publishedDate = ?, tags = ?, fileUrl = ? 
        WHERE bookid = ?`,
		book.Title, book.Author, book.Description, book.Category, book.ISBN, book.PublishedDate, tagsJSON, book.FileUrl, bookID)

	if err != nil {
		return err
	}

	return nil
}

func (s *BookStore) DeleteBook(bookID string) error {
	result, err := s.db.Exec("DELETE FROM books WHERE bookid = ?", bookID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows // No book was found with the given ID
	}

	return nil
}

func (s *BookStore) GetAllBooks() ([]*types.Book, error) {
	rows, err := s.db.Query("SELECT bookid, title, author, description, category, isbn, publishedDate, tags, fileUrl FROM books")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	books := []*types.Book{}
	for rows.Next() {
		book := new(types.Book)
		var tagsJSON []byte

		err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Description, &book.Category, &book.ISBN, &book.PublishedDate, &tagsJSON, &book.FileUrl)
		if err != nil {
			return nil, err
		}

		// Unmarshal the tags JSON back into a []string
		err = json.Unmarshal(tagsJSON, &book.Tags)
		if err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, nil
}

func (s *BookStore) SearchBooks(searchTerm string) ([]types.Book, error) {
	// SQL query to search books by title, author, or description (you can modify it as needed)
	query := `
        SELECT bookid, title, author, description, category, isbn, publishedDate, tags, fileUrl
        FROM books
        WHERE title LIKE ? OR author LIKE ? OR description LIKE ? OR tags LIKE ? OR category LIKE ?
    `

	// Use wildcard matching for search terms
	searchPattern := fmt.Sprintf("%%%s%%", searchTerm)

	// Execute the query
	rows, err := s.db.Query(query, searchPattern, searchPattern, searchPattern, searchPattern, searchPattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []types.Book

	for rows.Next() {
		var book types.Book
		var tags []byte // Store the raw JSON as []byte

		// Scan each row
		if err := rows.Scan(&book.BookID, &book.Title, &book.Author, &book.Description, &book.Category, &book.ISBN, &book.PublishedDate, &tags, &book.FileUrl); err != nil {
			return nil, err
		}

		// Decode tags JSON
		if err := json.Unmarshal(tags, &book.Tags); err != nil {
			return nil, err
		}

		books = append(books, book)
	}

	return books, nil
}

// Post a review for a book
func (s *BookStore) PostReview(userID, bookID string, review types.ReviewPayload) error {
	_, err := s.db.Exec(`
        INSERT INTO reviews (userid, bookid, rating, comment, createdAt) 
        VALUES (?, ?, ?, ?, NOW())`,
		userID, bookID, review.Rating, review.Comment)
	return err
}

// Get reviews for a book
func (s *BookStore) GetReviews(bookID string) ([]types.Review, error) {
	rows, err := s.db.Query(`
        SELECT reviewid, userid, bookid, rating, comment, createdAt 
        FROM reviews WHERE bookid = ?`, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reviews := []types.Review{}
	for rows.Next() {
		var review types.Review
		if err := rows.Scan(&review.ReviewID, &review.UserID, &review.BookID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}

	return reviews, nil
}

// Delete a user review
func (s *BookStore) DeleteUserReview(userID, reviewID string) error {
	_, err := s.db.Exec(`
        DELETE FROM reviews WHERE reviewid = ? AND userid = ?`, reviewID, userID)
	return err
}

// Admin deletes a review
func (s *BookStore) DeleteReview(reviewID string) error {
	_, err := s.db.Exec(`DELETE FROM reviews WHERE reviewid = ?`, reviewID)
	return err
}

// Add a book to user's favorites list
func (s *BookStore) AddToFavorites(userID, bookID string) error {
	_, err := s.db.Exec(`
        INSERT INTO favorites (userid, bookid, createdAt) 
        VALUES (?, ?, NOW())`, userID, bookID)
	return err
}
