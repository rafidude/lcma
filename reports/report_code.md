Here's the modern implementation of the legacy application using the target technology stack specified:

**main.go**
```go
package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

const (
	dbName     = "empl"
	dbUser     = "empl"
	dbPassword = "empl"
	dbHost     = "localhost"
	dbPort     = 5432
)

type App struct {
	db     *sqlx.DB
	router *chi.Mux
}

func main() {
	// Connect to database
	db, err := sqlx.Connect("pgx", "postgresql://"+dbUser+":"+dbPassword+"@"+dbHost+":"+strconv.Itoa(dbPort)+"/"+dbName)
	if err != nil {
		log.Fatal(err)
	}

	app := &App{db: db, router: chi.NewRouter()}

	// Define routes
	app.router.Get("/employees", app.listEmployees)
	app.router.Get("/employees/{id}", app.getEmployee)
	app.router.Post("/employees", app.addEmployee)
	app.router.Put("/employees/{id}", app.updateEmployee)
	app.router.Delete("/employees/{id}", app.deleteEmployee)

	// Start server
	http.ListenAndServe(":8080", app.router)
}
```

**employee.go**
```go
package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type Employee struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Age  int    `db:"age"`
}

func (app *App) listEmployees(w http.ResponseWriter, r *http.Request) {
	employees := make([]Employee, 0)

	err := app.db.Select(&employees, "SELECT * FROM employees")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, employees)
}

func (app *App) getEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	employee := Employee{}

	err := app.db.Get(&employee, "SELECT * FROM employees WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	render.JSON(w, r, employee)
}

func (app *App) addEmployee(w http.ResponseWriter, r *http.Request) {
	employee := Employee{}

	err := render.Bind(r, &employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = app.db.Exec("INSERT INTO employees (name, age) VALUES ($1, $2)", employee.Name, employee.Age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (app *App) updateEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	employee := Employee{}

	err := render.Bind(r, &employee)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = app.db.Exec("UPDATE employees SET name = $1, age = $2 WHERE id = $3", employee.Name, employee.Age, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (app *App) deleteEmployee(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := app.db.Exec("DELETE FROM employees WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
```

**schema.sql**
```sql
CREATE TABLE IF NOT EXISTS employees (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	age INTEGER NOT NULL
);
```

To run this application, you'll need to have Go, PostgreSQL, and the SQLX package installed on your system. You can install the SQLX package using the following command:
```bash
go get -u github.com/jmoiron/sqlx
```

Next, create a new PostgreSQL database and execute the schema.sql script to create the employees table.

Finally, build and run the application using the following commands:
```
go build main.go
./main
```

The application will be available at http://localhost:8080. You can use HTMX to send HTTP requests to the application and update the UI accordingly.