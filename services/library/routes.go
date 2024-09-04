// library/routes.go
package library

import (
	"fmt"
	"net/http"

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
