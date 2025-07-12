# Book Rental System API

A simple RESTful API for managing a book rental system using Go and Gin.

## Prerequisites

- Go 1.18 or higher

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd book-rental-system-go
```

2. Install dependencies:
```bash
go mod tidy
```

## Running the Application

### Development Mode
```bash
go run main.go
```

The server will start on `http://localhost:8081`.

### Production Build
```bash
go build -o book-rental-system main.go
./book-rental-system
```

## API Endpoints

### Books
- `POST /books` – Add a new book
  - Body: `{"title": "Book Title", "author": "Author Name"}`
- `GET /books` – List all books

### Users
- `POST /users` – Create a new user
  - Body: `{"name": "User Name", "email": "user@example.com"}`

### Rentals
- `POST /rent` – Rent a book to a user
  - Body: `{"user_id": "uuid", "book_id": "uuid"}`
- `POST /return` – Return a rented book
  - Body: `{"user_id": "uuid", "book_id": "uuid"}`

## Testing

### Run All Tests
```bash
go test ./...
```

### Run Specific Test Files
```bash
go test ./handlers
go test ./handlers -v
```

### Run Individual Test Functions
```bash
go test ./handlers -run TestAddBook
go test ./handlers -run TestRentBook
```

### Run Tests with Coverage
```bash
go test ./handlers -cover
go test ./handlers -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## Example Usage

### Using curl

1. Create a user:
```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Martin PJ", "email": "martin@example.com"}'
```

2. Add a book:
```bash
curl -X POST http://localhost:8080/books \
  -H "Content-Type: application/json" \
  -d '{"title": "Life of Pi", "author": "Yann Martel"}'
```

3. Rent a book (replace UUIDs with actual IDs from previous responses):
```bash
curl -X POST http://localhost:8080/rent \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-uuid", "book_id": "book-uuid"}'
```

4. List all books:
```bash
curl http://localhost:8080/books
```

5. Return a book:
```bash
curl -X POST http://localhost:8080/return \
  -H "Content-Type: application/json" \
  -d '{"user_id": "user-uuid", "book_id": "book-uuid"}'
```
