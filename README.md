# 📚 Student API

A RESTful API built in Go for managing student records. Supports multiple databases (PostgreSQL, MongoDB, SQLite) through a clean and modular architecture.

---

## 🚀 Features

- 🔁 Full CRUD operations for student records
- 🗄️ Configurable Multi-database support (PostgreSQL, MongoDB, SQLite) via environment variables
- 🧱 Modular code structure with reusable components
- ✅ Built with standard REST conventions

---

## 🧠 Tech Stack

- **Language**: Go (Golang)
- **Databases**: PostgreSQL, MongoDB, SQLite
- **Architecture**: Clean layered design with interfaces
- **Libraries**: `net/http`, custom middlewares, Go standard libraries

---

## 📁 Project Structure

student-api/
├── cmd/
│   └── student-api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── http/
│   │   └── handler/         # HTTP request handlers
│   ├── storage/
│   │   ├── postgreSql/      # PostgreSQL implementation
│   │   ├── mongodb/         # MongoDB implementation
│   │   ├── sqlite/          # SQLite implementation
│   │   └── storage.go       # Storage interface
│   ├── types/               # Data structures and models
│   └── utils/
│       └── response/        # Standardized response utilities
├── go.mod                   # Go module dependencies
└── go.sum                   # Dependency checksums

