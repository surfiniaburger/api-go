package library

import (
	"database/sql"
	"encoding/json"
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

func (s *BookStore) GetBookByID(bookID int) (*types.Book, error) {
	row := s.db.QueryRow(`
		SELECT bookId, title, author, description, category, isbn, publishedDate, tags, fileUrl 
		FROM books WHERE bookId = ?`, bookID)

	book := new(types.Book)
	var tagsJSON []byte
	err := row.Scan(&book.BookID, &book.Title, &book.Author, &book.Description, &book.Category, &book.ISBN, &book.PublishedDate, &tagsJSON, &book.FileUrl)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(tagsJSON, &book.Tags)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (s *BookStore) UpdateBook(bookID int, book types.UpdateBookPayload) error {
	tagsJSON, err := json.Marshal(book.Tags)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(`
		UPDATE books 
		SET title = ?, author = ?, description = ?, category = ?, isbn = ?, publishedDate = ?, tags = ?, fileUrl = ? 
		WHERE bookId = ?`,
		book.Title, book.Author, book.Description, book.Category, book.ISBN, book.PublishedDate, tagsJSON, book.FileUrl, bookID)

	if err != nil {
		return err
	}

	return nil
}
