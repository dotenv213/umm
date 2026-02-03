# User Management Module (Go + SQLite)
A robust, persistent, and thread-safe user management system implemented in Go using an embedded SQLite database. This module is designed to perform CRUD operations without using an ORM, prioritizing performance, data integrity, and security.

## Project Overview
The goal of this project is to provide a standalone module for managing users with persistent storage. It handles:
- **C**reation of users with unique constraints.
- **R**etrieval of user details by ID or listing all users.
- **U**pdate of user information safely.
- **D**eletion of users.

**Key Features:**
- **No ORM:** Pure SQL for maximum control and performance.
- **Persistence:** Data survives application restarts (stored in `users.db`).
- **Security:** SQL Injection prevention using parameterized queries.
- **High Concurrency:** SQLite configured in **WAL (Write-Ahead Logging)** mode.
- **Durability:** Critical operations (Create) use **Transactions** to ensure atomicity.

---

## Architecture & Design

### Package Structure
The project follows the Standard Go Project Layout:
- `cmd/main.go`: The entry point and CLI demonstration of the module.
- `internal/userstore/`: Encapsulates all database logic. This code is private to the project and cannot be imported by external modules, enforcing clean separation of concerns.

### Database Schema
The database consists of a single `users` table designed for extensibility:

| Field | Type | Description |
|-------|------|-------------|
| `id` | `INTEGER` | Primary Key, Auto-incremented. |
| `username` | `TEXT` | Unique, Non-null. Uniquely identifies a user. |
| `email` | `TEXT` | Unique, Non-null. Used for communication. |
| `created_at` | `DATETIME` | Defaults to `CURRENT_TIMESTAMP`. Tracks registration time. |

### Persistence & Durability Approach
1.  **Storage:** Data is stored in a local file (`users.db`), not in memory.
2.  **WAL Mode:** `PRAGMA journal_mode = WAL;` is enabled to allow concurrent reads and writes, preventing database locks during high load.
3.  **Transactions:** Creation operations are wrapped in `BeginTx`, `Commit`, and `Rollback` patterns to ensure that data is never left in an inconsistent state if a crash occurs.
4.  **Foreign Keys:** `PRAGMA foreign_keys = ON;` is set to ensure future extensibility (e.g., adding a `posts` or `orders` table linked to users).

---

## 🚀 Build Instructions

### Prerequisites
1.  **Go**: Version 1.18 or higher.
2.  **GCC Compiler**: Required because `go-sqlite3` uses CGO.
    - *Linux:* Install `build-essential` or `gcc`.
    - *Windows:* Install MinGW or TDM-GCC.
    - *macOS:* Install Xcode Command Line Tools (`xcode-select --install`).

### Installation
Clone the repository and install dependencies:
```bash
git clone https://github.com/dotenv213/umm.git
cd umm/
go mod tidy
```

### Building the Project
To compile the application:
```bash
go build -o umm-cli cmd/main.go
```

### Running the Application
```bash
go run cmd/main.go
```
The application will create a users.db file in the root directory automatically.

---

## Testing Instructions
The project includes comprehensive unit tests covering >85% of the code, including success paths, failure scenarios, and edge cases.

### Running Tests
To execute all tests with verbose output:
```bash
go test -v ./internal/userstore/...
```

### Checking Coverage
To verify the code coverage requirement (Target: >80%):
```bash
go test -v -cover ./internal/userstore/...
```

Current Coverage: ~87%
```bash
=== RUN   TestCreateUser
--- PASS: TestCreateUser (0.00s)
=== RUN   TestCreateDuplicateUser
--- PASS: TestCreateDuplicateUser (0.00s)
=== RUN   TestGetByID
--- PASS: TestGetByID (0.00s)
=== RUN   TestGetUserNotFound
--- PASS: TestGetUserNotFound (0.00s)
=== RUN   TestListAllUsers
--- PASS: TestListAllUsers (0.00s)
=== RUN   TestUpdateUser
--- PASS: TestUpdateUser (0.00s)
=== RUN   TestUpdateNonExistUser
--- PASS: TestUpdateNonExistUser (0.00s)
=== RUN   TestDeleteUser
--- PASS: TestDeleteUser (0.00s)
=== RUN   TestDeleteNonExistUser
--- PASS: TestDeleteNonExistUser (0.00s)
=== RUN   TestStoreClose
--- PASS: TestStoreClose (0.00s)
=== RUN   TestListEmpty
--- PASS: TestListEmpty (0.00s)
=== RUN   TestGetByIDDetails
--- PASS: TestGetByIDDetails (0.00s)
=== RUN   TestNewDbError
--- PASS: TestNewDbError (0.00s)
=== RUN   TestOperationsOnClosedDB
--- PASS: TestOperationsOnClosedDB (0.00s)
PASS
coverage: 87.8% of statements
ok      github.com/dotenv213/umm/internal/userstore     0.009s  coverage: 87.8% of statements
```

### Tested Scenarios
1. **CRUD Operations:** Create, Get, Update, Delete.
2. **Edge Cases:** Fetching non-existent users, Deleting deleted users.
3. **Constraints:** Attempting to create users with duplicate usernames/emails.
4. **System Failures:** Operations on a closed database connection.
5. **Data Integrity:** Verifying fields are correctly stored and retrieved.

Tests use an in-memory database (:memory:) to ensure speed and isolation, preventing side effects on the persistent users.db file.

## Project Structure

```text
.
├── cmd/
│   └── main.go           # CLI entry point 
├── internal/
│   └── userstore/        # Core logic package
│       ├── model.go      # User struct definition
│       ├── store.go      # Interface definition
│       ├── sqlite.go     # SQLite implementation & SQL queries
│       ├── errors.go     # Custom error variables
│       └── store_test.go # Unit tests
├── go.mod                # Module definition
├── go.sum                # Checksums
├── .gitignore            # gitignore
└── README.md             # Documentation
```