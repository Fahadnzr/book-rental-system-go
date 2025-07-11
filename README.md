# Book Rental System API

A simple RESTful API for managing a book rental system using Go and Gin.

## Prerequisites

- Go 1.18 or higher

## Install dependencies

```
go mod tidy
```

## Run the server

```
go run main.go
```

The server will start on `http://localhost:8080`.

## API Endpoints

- `POST /books` – Add a new book
- `GET /books` – List all books
- `POST /users` – Create a new user
- `POST /rent` – Rent a book to a user
- `POST /return` – Return a rented book

## Run tests

```
go test ./handlers
```
