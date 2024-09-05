package types

import (
	"time"
)

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"` // Add this field
	CreatedAt time.Time `json:"createdAt"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price"`
	// note that this isn't the best way to handle quantity
	// because it's not atomic (in ACID), but it's good enough for this example
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"createdAt"`
}

type CartCheckoutItem struct {
	ProductID int `json:"productID"`
	Quantity  int `json:"quantity"`
}

type Order struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userID"`
	Total     float64   `json:"total"`
	Status    string    `json:"status"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"createdAt"`
}

type OrderItem struct {
	ID        int       `json:"id"`
	OrderID   int       `json:"orderID"`
	ProductID int       `json:"productID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
}

type ProductStore interface {
	GetProductByID(id int) (*Product, error)
	GetProductsByID(ids []int) ([]Product, error)
	GetProducts() ([]*Product, error)
	CreateProduct(CreateProductPayload) error
	UpdateProduct(Product) error
}

type BookStore interface {
	CreateBook(book CreateBookPayload) error
	GetBookByID(bookID string) (*Book, error) // Ensure it's string
	UpdateBook(bookID string, book UpdateBookPayload) error
	DeleteBook(bookID string) error
	GetAllBooks() ([]*Book, error)
	SearchBooks(searchTerm string) ([]Book, error)
	// Add more methods as needed
}

type OrderStore interface {
	CreateOrder(Order) (int, error)
	CreateOrderItem(OrderItem) error
}
type CreateProductPayload struct {
	Name        string  `json:"name" validate:"required"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Price       float64 `json:"price" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
}

type RegisterUserPayload struct {
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=3,max=130"`
	Role      string `json:"role"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type CartCheckoutPayload struct {
	Items []CartCheckoutItem `json:"items" validate:"required"`
}

type CreateBookPayload struct {
	BookID        string   `json:"bookId" validate:"required"`
	Title         string   `json:"title" validate:"required"`
	Author        string   `json:"author" validate:"required"`
	Description   string   `json:"description"`
	Category      string   `json:"category" validate:"required"`
	ISBN          string   `json:"isbn"`
	PublishedDate string   `json:"publishedDate"`
	Tags          []string `json:"tags"`
	FileUrl       string   `json:"fileUrl"`
}

type Book struct {
	BookID        string   `json:"bookId"`
	Title         string   `json:"title"`
	Author        string   `json:"author"`
	Description   string   `json:"description"`
	Category      string   `json:"category"`
	ISBN          string   `json:"isbn"`
	PublishedDate string   `json:"publishedDate"`
	Tags          []string `json:"tags"`
	FileUrl       string   `json:"fileUrl"`
}

type UpdateBookPayload struct {
	Title         string   `json:"title" validate:"required"`
	Author        string   `json:"author" validate:"required"`
	Description   string   `json:"description"`
	Category      string   `json:"category" validate:"required"`
	ISBN          string   `json:"isbn"`
	PublishedDate string   `json:"publishedDate"`
	Tags          []string `json:"tags"`
	FileUrl       string   `json:"fileUrl"`
}
