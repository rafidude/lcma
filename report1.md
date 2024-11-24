Legacy Codebase Documentation
==========================

### Architecture Overview and Main Components

The legacy codebase is built using Flask, a micro web framework written in Python. The application consists of several key components:

1.  **Routes**: The application defines several routes for listing employees, adding new employees, editing existing employees, and deleting employees.
2.  **Database Interactions**: The application uses the `psycopg` library to interact with a PostgreSQL database. Database operations are performed using raw SQL queries.
3.  **Templates**: The application uses Jinja2 templates for rendering HTML templates.
4.  **Static Files**: The application serves static files directly.

### Key Business Logic and Workflows

1.  **Employee Listing**: The application retrieves a list of employees from the database and displays them in a table.
2.  **Employee Creation**: The application accepts user input for employee name and age, validates the input, and inserts the new employee into the database.
3.  **Employee Editing**: The application retrieves an existing employee from the database, accepts user input for updated employee details, and updates the employee record in the database.
4.  **Employee Deletion**: The application accepts an employee ID, retrieves the corresponding employee record from the database, and deletes the record.

### Database Schema and Relationships

The database schema consists of a single table named `employees` with the following columns:

| Column Name | Data Type |
| :--------- | :------- |
| id         | integer   |
| name       | varchar   |
| age        | integer   |

There are no explicit relationships defined between tables in the database schema.

### External Dependencies and Integrations

The application depends on the following external libraries:

1.  `Flask` for building the web application
2.  `psycopg` for interacting with the PostgreSQL database
3.  `dotenv` for loading environment variables from a `.env` file

Technical Debt and Potential Issues
=====================================

### Security Vulnerabilities

1.  **SQL Injection**: The application uses raw SQL queries, which makes it vulnerable to SQL injection attacks.
2.  **Cross-Site Scripting (XSS)**: The application does not perform input validation or sanitization for user input, making it vulnerable to XSS attacks.

### Performance Bottlenecks

1.  **Database Connection**: The application establishes a new database connection for each request, which can lead to performance issues under high load.
2.  **Query Optimization**: The application uses simple queries that may not be optimized for performance.

### Maintainability Concerns

1.  **Code Organization**: The application code is not organized into separate modules or packages, making it difficult to maintain.
2.  **Lack of Comments**: The application code does not include comments or documentation, making it difficult for new developers to understand the codebase.

### Outdated Dependencies or Deprecated Features

1.  **Flask Version**: The application uses an outdated version of Flask, which may not include the latest security patches and features.
2.  **psycopg Version**: The application uses an outdated version of psycopg, which may not include the latest security patches and features.

### Missing Error Handling or Edge Cases

1.  **Database Errors**: The application does not handle database errors, which can lead to unexpected behavior or crashes.
2.  **Input Validation**: The application does not perform input validation, which can lead to unexpected behavior or crashes.

### Code Smells and Anti-Patterns

1.  **Global Database Connection**: The application uses a global database connection, which can lead to performance issues and make the code harder to maintain.

Modern Implementation
====================

### Project Structure

The modern implementation will follow the standard Go project structure:
```markdown
project/
├── cmd/
│   └── main.go
├── internal/
│   ├── database/
│   │   └── employees.go
│   ├── handlers/
│   │   ├── employees.go
│   │   └── index.go
│   └── models/
│       └── employee.go
├── pkg/
│   └── utils/
│       └── error.go
├── public/
├── templates/
│   ├── employees.html
│   ├── index.html
│   └── layout.html
├── tests/
├── vendor/
├── go.mod
├── go.sum
└── main.go
```

### Database Configuration

We will use the `sqlx` package to interact with the database. We will define a `database` package that exports a `DB` struct, which will be used to access the database.
```go
// internal/database/employees.go
package database

import (
	"database/sql"
	"fmt"

	"example.com/project/internal/models"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	*sqlx.DB
}

func NewDB(connection string) (*DB, error) {
	db, err := sqlx.Connect("postgres", connection)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (db *DB) GetEmployees() ([]models.Employee, error) {
	employees := []models.Employee{}

	err := db.Select(&employees, "SELECT * FROM employees")
	if err != nil {
		return nil, err
	}

	return employees, nil
}

func (db *DB) CreateEmployee(employee *models.Employee) error {
	_, err := db.NamedExec("INSERT INTO employees (name, age) VALUES (:name, :age)", employee)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) UpdateEmployee(employee *models.Employee) error {
	_, err := db.NamedExec("UPDATE employees SET name = :name, age = :age WHERE id = :id", employee)
	if err != nil {
		return err
	}

	return nil
}

func (db *DB) DeleteEmployee(id int) error {
	_, err := db.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
```

### HTTP Handlers

We will use the `chi` package to define HTTP handlers. We will define a `handlers` package that exports an `EmployeesHandler` struct, which will be used to handle HTTP requests.
```go
// internal/handlers/employees.go
package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/project/internal/database"
	"example.com/project/internal/models"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type EmployeesHandler struct {
	db *database.DB
}

func NewEmployeesHandler(db *database.DB) *EmployeesHandler {
	return &EmployeesHandler{db}
}

func (h *EmployeesHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	employees, err := h.db.GetEmployees()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(employees)
}

func (h *EmployeesHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	var employee models.Employee
	err := json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.CreateEmployee(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *EmployeesHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var employee models.Employee
	err = json.NewDecoder(r.Body).Decode(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.UpdateEmployee(&employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *EmployeesHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.db.DeleteEmployee(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
```

### Main Function

We will define a `main` function that creates a new instance of the `EmployeesHandler` struct and defines the HTTP routes.
```go
// main.go
package main

import (
	"fmt"
	"log"

	"example.com/project/internal/database"
	"example.com/project/internal/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	db, err := database.NewDB("user:password@localhost/database")
	if err != nil {
		log.Fatal(err)
	}

	handler := handlers.NewEmployeesHandler(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/employees", handler.GetEmployees)
	r.Post("/employees", handler.CreateEmployee)
	r.Put("/employees/{id}", handler.UpdateEmployee)
	r.Delete("/employees/{id}", handler.DeleteEmployee)

	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
```

### templates

We will use the `HTMX` package to define HTML templates. We will define a `layout.html` template that includes a navigation bar and a `content` block.
```html
<!-- layout.html -->
<!DOCTYPE html>
<html>
<head>
	<title>Employees</title>
</head>
<body>
	<nav>
		<ul>
			<li><a href="{{url "/employees"}}">Employees</a></li>
		</ul>
	</nav>
	<div id="content">
		{{block "content" .}}{{end}}
	</div>
</body>
</html>
```

We will define an `employees.html` template that includes a table with a list of employees.
```html
<!-- employees.html -->
{{define "content"}}
<h1>Employees</h1>
<table>
	<thead>
		<tr>
			<th>ID</th>
			<th>Name</th>
			<th>Age</th>
		</tr>
	</thead>
	<tbody>
		{{range .Employees}}
		<tr>
			<td>{{.ID}}</td>
			<td>{{.Name}}</td>
			<td>{{.Age}}</td>
		</tr>
		{{end}}
	</tbody>
</table>
{{end}}
```

We will define an `index.html` template that includes a link to the `employees` page.
```html
<!-- index.html -->
{{define "content"}}
<h1>Home</h1>
<p><a href="{{url "/employees"}}">View Employees</a></p>
{{end}}
```