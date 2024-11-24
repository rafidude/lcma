```json
{
  "project_name": "EmployeeManagementSystem",
  "file_structure": {
    "main": {
      "go.mod": "",
      "main.go": ""
    },
    "internal": {
      "database": {
        "database.go": "",
        "employee.go": "",
        "employee_repository.go": ""
      },
      "handlers": {
        "employee_handler.go": ""
      },
      "models": {
        "employee.go": ""
      },
      "server": {
        "server.go": ""
      },
      "templates": {
        "base.html": "",
        "list.html": "",
        "edit.html": "",
        "add.html": ""
      }
    }
  }
}
```

```go
// go.mod
module EmployeeManagementSystem

go 1.19

require (
	github.com/go-chi/chi/v5 v5.0.4
	github.com/go-chi/cors v1.2.0
	github.com/jmoiron/sqlx v1.3.4
	github.com/lib/pq v1.10.7
	github.com/mattn/go-sqlite3 v1.14.13
)
```

```go
// database/database.go
package database

import (
	"database/sql"
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	conn *sqlx.DB
}

func NewDatabase(dsn string) (*Database, error) {
	conn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return &Database{conn: conn}, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}
```

```go
// database/employee.go
package database

import (
	"database/sql"
	"errors"
)

type EmployeeRepository struct {
	db *Database
}

func NewEmployeeRepository(db *Database) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

func (e *EmployeeRepository) FindAll() ([]*models.Employee, error) {
	employees := []*models.Employee{}
	err := e.db.conn.Select(&employees, "SELECT * FROM employees")
	return employees, err
}

func (e *EmployeeRepository) FindOne(id int) (*models.Employee, error) {
	employee := &models.Employee{}
	err := e.db.conn.Get(employee, "SELECT * FROM employees WHERE id=$1", id)
	if err != nil {
		return nil, err
	}
	return employee, nil
}

func (e *EmployeeRepository) Create(employee *models.Employee) error {
	_, err := e.db.conn.Exec("INSERT INTO employees (name, age) VALUES ($1, $2)", employee.Name, employee.Age)
	return err
}

func (e *EmployeeRepository) Update(employee *models.Employee) error {
	_, err := e.db.conn.Exec("UPDATE employees SET name=$1, age=$2 WHERE id=$3", employee.Name, employee.Age, employee.ID)
	return err
}

func (e *EmployeeRepository) Delete(id int) error {
	_, err := e.db.conn.Exec("DELETE FROM employees WHERE id=$1", id)
	return err
}
```

```go
// handlers/employee_handler.go
package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mattn/go-sqlite3"
	"EmployeeManagementSystem/internal/database"
	"EmployeeManagementSystem/internal/models"
)

type EmployeeHandler struct {
	repo *database.EmployeeRepository
}

func NewEmployeeHandler(repo *database.EmployeeRepository) *EmployeeHandler {
	return &EmployeeHandler{repo: repo}
}

func (e *EmployeeHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	employees, err := e.repo.FindAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employees)
}

func (e *EmployeeHandler) FindOne(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	employee, err := e.repo.FindOne(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

func (e *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	employee := &models.Employee{}
	err := json.NewDecoder(r.Body).Decode(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = e.repo.Create(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

func (e *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	employee := &models.Employee{}
	err := json.NewDecoder(r.Body).Decode(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = e.repo.Update(employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(employee)
}

func (e *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := e.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
```

```go
// models/employee.go
package models

type Employee struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}
```

```go
// server/server.go
package server

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"EmployeeManagementSystem/internal/database"
	"EmployeeManagementSystem/internal/handlers"
	"EmployeeManagementSystem/internal/templates"
)

type Server struct {
	db  *database.Database
	repo *database.EmployeeRepository
}

func NewServer(db *database.Database) *Server {
	return &Server{db: db, repo: database.NewEmployeeRepository(db)}
}

func (s *Server) Close() error {
	return s.db.Close()
}

func (s *Server) GetRouter() *chi.Mux {
	r := chi.NewRouter()
	cors := cors.New(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:  []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)
	r.Get("/employees", s.employeeHandler.FindAll)
	r.Get("/employees/{id}", s.employeeHandler.FindOne)
	r.Post("/employees", s.employeeHandler.Create)
	r.Put("/employees/{id}", s.employeeHandler.Update)
	r.Delete("/employees/{id}", s.employeeHandler.Delete)
	r.Get("/employees", s.employeeTemplateHandler.FindAll)
	r.Get("/employees/{id}", s.employeeTemplateHandler.FindOne)
	r.Get("/employees/add", s.employeeTemplateHandler.Add)
	r.Get("/employees/{id}/edit", s.employeeTemplateHandler.Edit)
	r.Post("/employees", s.employeeTemplateHandler.Create)
	r.Put("/employees/{id}", s.employeeTemplateHandler.Update)
	r.Delete("/employees/{id}", s.employeeTemplateHandler.Delete)
	return r
}
```

```go
// templates/base.html
<!DOCTYPE html>
<html>
<head>
    <title>Employee Management System</title>
    <script src="https://unpkg.com/htmx.org@1.7.0/dist/htmx.min.js"></script>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
</head>
<body>
    <div class="container mx-auto p-4 mt-4">
        {{block "content" .}}{{end}}
    </div>
</body>
</html>
```

```go
// templates/list.html
{{ define "content" }}
    <h1 class="text-3xl font-bold mb-4">Employees</h1>
    <table class="table-auto w-full mb-4">
        <thead>
            <tr>
                <th class="px-4 py-2">ID</th>
                <th class="px-4 py-2">Name</th>
                <th class="px-4 py-2">Age</th>
                <th class="px-4 py-2">Actions</th>
            </tr>
        </thead>
        <tbody>
            {{ range .employees }}
            <tr>
                <td class="px-4 py-2">{{ .ID }}</td>
                <td class="px-4 py-2">{{ .Name }}</td>
                <td class="px-4 py-2">{{ .Age }}</td>
                <td class="px-4 py-2">
                    <a href="/employees/{{ .ID }}" class="text-blue-500 hover:text-blue-700">Edit</a>
                    <a href="/employees/{{ .ID }}/delete" class="text-red-500 hover:text-red-700">Delete</a>
                </td>
            </tr>
            {{ end }}
        </tbody>
    </table>
    <a href="/employees/add" class="text-blue-500 hover:text-blue-700">Add New Employee</a>
{{ end }}
```

```go
// templates/add.html
{{ define "content" }}
    <h1 class="text-3xl font-bold mb-4">Add New Employee</h1>
    <form class="w-full max-w-lg" hx-post="/employees">
        <div class="flex flex-wrap -mx-3 mb-6">
            <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="grid-first-name">
                    Name
                </label>
                <input class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-first-name" type="text" name="name">
            </div>
            <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="grid-last-name">
                    Age
                </label>
                <input class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-last-name" type="number" name="age">
            </div>
        </div>
        <div class="flex justify-end">
            <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
                Add Employee
            </button>
        </div>
    </form>
{{ end }}
```

```go
// templates/edit.html
{{ define "content" }}
    <h1 class="text-3xl font-bold mb-4">Edit Employee</h1>
    <form class="w-full max-w-lg" hx-put="/employees/{{ .employee.ID }}">
        <div class="flex flex-wrap -mx-3 mb-6">
            <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="grid-first-name">
                    Name
                </label>
                <input class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-first-name" type="text" name="name" value="{{ .employee.Name }}">
            </div>
            <div class="w-full md:w-1/2 px-3 mb-6 md:mb-0">
                <label class="block uppercase tracking-wide text-gray-700 text-xs font-bold mb-2" for="grid-last-name">
                    Age
                </label>
                <input class="appearance-none block w-full bg-gray-200 text-gray-700 border border-gray-200 rounded py-3 px-4 mb-3 leading-tight focus:outline-none focus:bg-white focus:border-gray-500" id="grid-last-name" type="number" name="age" value="{{ .employee.Age }}">
            </div>
        </div>
        <div class="flex justify-end">
            <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
                Update Employee
            </button>
        </div>
    </form>
{{ end }}
```

```go
// templates/delete.html
{{ define "content" }}
    <h1 class="text-3xl font-bold mb-4">Delete Employee</h1>
    <p class="mb-4">Are you sure you want to delete this employee?</p>
    <form class="w-full max-w-lg" hx-delete="/employees/{{ .employee.ID }}">
        <div class="flex justify-end">
            <button class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline" type="submit">
                Delete Employee
            </button>
        </div>
    </form>
{{ end }}
```