// library/routes.go
package library

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/surfiniaburger/api-go/services/auth"
	"github.com/surfiniaburger/api-go/types"
	"github.com/surfiniaburger/api-go/utils"
)

type BookHandler struct {
	bookStore types.BookStore
	userStore types.UserStore
}

func NewBookHandler(bookStore types.BookStore, userStore types.UserStore) *BookHandler {
	return &BookHandler{bookStore: bookStore, userStore: userStore}
}

func (h *BookHandler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/admin/library", auth.WithJWTAuth(h.handleCreateBook, h.userStore, "admin")).Methods(http.MethodPost)
	router.HandleFunc("/admin/library/{id}", auth.WithJWTAuth(h.handleUpdateBook, h.userStore, "admin")).Methods(http.MethodPut)
	router.HandleFunc("/admin/library/{id}", auth.WithJWTAuth(h.handleDeleteBook, h.userStore, "admin")).Methods(http.MethodDelete)
	router.HandleFunc("/admin/library", auth.WithJWTAuth(h.handleGetAllBooks, h.userStore, "admin")).Methods(http.MethodGet)

	// User endpoint to get all books
	router.HandleFunc("/library", auth.WithJWTAuth(h.handleGetAllBooksForUsers, h.userStore, "user")).Methods(http.MethodGet)
	// Route for users to get a specific book by bookid
	router.HandleFunc("/library/{bookid}", h.handleGetBookByID).Methods(http.MethodGet)

	// Search the books in the library
	router.HandleFunc("/library/search", h.handleSearchBooks).Methods(http.MethodGet)

	// Route for users to post a review on a book
	router.HandleFunc("/library/{bookid}/reviews", auth.WithJWTAuth(h.handlePostReview, h.userStore, "user")).Methods(http.MethodPost)

	// Route for users to get reviews for a book
	router.HandleFunc("/library/{bookid}/reviews", h.handleGetReviews).Methods(http.MethodGet)

	// Route for users to delete their own review
	router.HandleFunc("/library/reviews/{reviewid}", auth.WithJWTAuth(h.handleDeleteUserReview, h.userStore, "user")).Methods(http.MethodDelete)

	// Admin route to delete any review
	router.HandleFunc("/admin/library/{bookid}/reviews/{reviewid}", auth.WithJWTAuth(h.handleDeleteAdminReview, h.userStore, "admin")).Methods(http.MethodDelete)

	// Route for users to add books to their favorites list
	router.HandleFunc("/library/user/lists", auth.WithJWTAuth(h.handleAddToFavorites, h.userStore, "user")).Methods(http.MethodPost)
}

func (h *BookHandler) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	var book types.CreateBookPayload
	if err := utils.ParseJSON(r, &book); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(book); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.bookStore.CreateBook(book)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, book)
}

func (h *BookHandler) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	var book types.UpdateBookPayload
	if err := utils.ParseJSON(r, &book); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	if err := utils.Validate.Struct(book); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteError(w, http.StatusBadRequest, fmt.Errorf("invalid payload: %v", errors))
		return
	}

	err := h.bookStore.UpdateBook(bookID, book)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, book)
}

func (h *BookHandler) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["id"]

	err := h.bookStore.DeleteBook(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("book not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "book deleted successfully"})
}

func (h *BookHandler) handleGetAllBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookStore.GetAllBooks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, books)
}

func (h *BookHandler) handleGetAllBooksForUsers(w http.ResponseWriter, r *http.Request) {
	books, err := h.bookStore.GetAllBooks()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, books)
}

func (h *BookHandler) handleGetBookByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookid"] // Extract bookid from the request URL

	// Get the book details from the BookStore
	book, err := h.bookStore.GetBookByID(bookID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("book not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	// Respond with the book details
	utils.WriteJSON(w, http.StatusOK, book)
}

func (h *BookHandler) handleSearchBooks(w http.ResponseWriter, r *http.Request) {
	// Get search term from query parameters
	searchTerm := r.URL.Query().Get("search")
	if searchTerm == "" {
		// Create an error from the string message
		err := fmt.Errorf("search term is required")
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	// Call the BookStore method to search books
	books, err := h.bookStore.SearchBooks(searchTerm)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	// Return the list of books
	utils.WriteJSON(w, http.StatusOK, books)
}

// POST /library/{bookid}/reviews - Post a review for a book
func (h *BookHandler) handlePostReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookid"]

	var review types.ReviewPayload
	if err := utils.ParseJSON(r, &review); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context()) // Get user ID from JWT token
	userIDStr := strconv.Itoa(userID)                // Convert userID from int to string

	err := h.bookStore.PostReview(userIDStr, bookID, review)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusCreated, map[string]string{"message": "review posted"})
}

// GET /library/{bookid}/reviews - Get reviews for a book
func (h *BookHandler) handleGetReviews(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookID := vars["bookid"]

	reviews, err := h.bookStore.GetReviews(bookID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, reviews)
}

// DELETE /library/reviews/{reviewid} - Delete a user review
func (h *BookHandler) handleDeleteUserReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["reviewid"]

	userID := auth.GetUserIDFromContext(r.Context())
	userIDStr := strconv.Itoa(userID) // Convert userID from int to string

	err := h.bookStore.DeleteUserReview(userIDStr, reviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("review not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "review deleted"})
}

// DELETE /admin/library/{bookid}/reviews/{reviewid} - Admin deletes a review
func (h *BookHandler) handleDeleteAdminReview(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID := vars["reviewid"]

	err := h.bookStore.DeleteReview(reviewID)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.WriteError(w, http.StatusNotFound, fmt.Errorf("review not found"))
		} else {
			utils.WriteError(w, http.StatusInternalServerError, err)
		}
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "review deleted by admin"})
}

// POST /library/user/lists - Add a book to user's favorites list
func (h *BookHandler) handleAddToFavorites(w http.ResponseWriter, r *http.Request) {
	var favorite types.FavoritePayload
	if err := utils.ParseJSON(r, &favorite); err != nil {
		utils.WriteError(w, http.StatusBadRequest, err)
		return
	}

	userID := auth.GetUserIDFromContext(r.Context())
	userIDStr := strconv.Itoa(userID) // Convert userID from int to string

	err := h.bookStore.AddToFavorites(userIDStr, favorite.BookID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSON(w, http.StatusOK, map[string]string{"message": "book added to favorites"})
}
