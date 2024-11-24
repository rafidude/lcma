Below is the modern UI implementation using the target technology stack:

**index.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>Employees List</title>
    <link href="https://cdn.tailwindcss.com" rel="stylesheet">
</head>
<body>
    <h1 class="text-3xl">Employees</h1>
    <a href="/add" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Add New Employee</a>
    <div class="flex justify-center">
        <table class="table-auto">
            <thead class="bg-gray-100">
                <tr>
                    <th class="px-4 py-2">ID</th>
                    <th class="px-4 py-2">Name</th>
                    <th class="px-4 py-2">Age</th>
                    <th class="px-4 py-2">Actions</th>
                </tr>
            </thead>
            <tbody hx-target="this" hx-swap="innerHTML" hx-get="/api/employees">
                <!-- table rows will be loaded dynamically -->
            </tbody>
        </table>
    </div>

    <script src="https://unpkg.com/htmx.org@1.7.0/dist/htmx.min.js"></script>
</body>
</html>
```

**add.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>Add Employee</title>
    <link href="https://cdn.tailwindcss.com" rel="stylesheet">
</head>
<body>
    <h1 class="text-3xl">Add New Employee</h1>
    <form hx-post="/api/employees" hx-target="#result" hx-swap="innerHTML" hx-trigger="submit">
        <!-- using hx-post, hx-target, hx-swap and hx-trigger attributes for HTMX -->
        <p>Name: <input type="text" name="name" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"></p>
        <p>Age: <input type="number" name="age" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"></p>
        <p><input type="submit" value="Add Employee" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"></p>
    </form>
    <div id="result"></div>
    <a href="/" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Back to List</a>

    <script src="https://unpkg.com/htmx.org@1.7.0/dist/htmx.min.js"></script>
</body>
</html>
```

**edit.html**
```html
<!DOCTYPE html>
<html>
<head>
    <title>Edit Employee</title>
    <link href="https://cdn.tailwindcss.com" rel="stylesheet">
</head>
<body>
    <h1 class="text-3xl">Edit Employee</h1>
    <form hx-post="/api/employees/{{ .id }}" hx-target="#result" hx-swap="innerHTML" hx-trigger="submit">
        <p>Name: <input type="text" name="name" value="{{ .name }}" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"></p>
        <p>Age: <input type="number" name="age" value="{{ .age }}" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"></p>
        <p><input type="submit" value="Update Employee" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"></p>
    </form>
    <div id="result"></div>
    <a href="/" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Back to List</a>

    <script src="https://unpkg.com/htmx.org@1.7.0/dist/htmx.min.js"></script>
</body>
</html>
```

Note that I have assumed Golang will be used for API creation which will be used by HTMX. In this example, `/api/employees` and `/api/employees/{id}` API endpoints will be used by `index.html` and `edit.html` respectively.

The above example assumes you will be using Go-Chi for API routing and PostgreSQL for database. You can modify the code as per your requirement.

Below is a basic example of how you might implement the API using Go-Chi:

```go
package main

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi"
    "github.com/go-chi/chi/middleware"
    "golang.org/x/exp/slog"

    _ "github.com/lib/pq"
)

type Employee struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
    Age  int    `json:"age"`
}

var db *sql.DB

func getEmployees(w http.ResponseWriter, r *http.Request) {
    rows, err := db.Query("SELECT * FROM employees")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var employees []Employee
    for rows.Next() {
        var employee Employee
        err = rows.Scan(&employee.ID, &employee.Name, &employee.Age)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        employees = append(employees, employee)
    }

    json.NewEncoder(w).Encode(employees)
}

func addEmployee(w http.ResponseWriter, r *http.Request) {
    var employee Employee
    err := json.NewDecoder(r.Body).Decode(&employee)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("INSERT INTO employees (name, age) VALUES ($1, $2)", employee.Name, employee.Age)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func editEmployee(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    var employee Employee
    err := json.NewDecoder(r.Body).Decode(&employee)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    _, err = db.Exec("UPDATE employees SET name = $1, age = $2 WHERE id = $3", employee.Name, employee.Age, id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
}

func main() {
    var err error
    db, err = sql.Open("postgres", "user=myuser dbname=mydb password=mypass sslmode=disable")
    if err != nil {
        slog.Fatal(err)
    }
    defer db.Close()

    router := chi.NewRouter()
    router.Use(middleware.Logger, middleware.Recoverer)

    router.Get("/api/employees", getEmployees)
    router.Post("/api/employees/{id}", editEmployee)
    router.Post("/api/employees", addEmployee)

    http.ListenAndServe(":8080", router)
}
```