# Create CRUD Api with Go

GET : http://localhost:8080/books

---

POST : http://localhost:8080/books

```json
{
  "title": "New Book Title",
  "author": "Author Name"
}

```

---

PUT : http://localhost:8080/books/{id}

```json
{
  "title": "Updated Book Title",
  "author": "Updated Author"
}
```

---

DELETE : http://localhost:8080/books/{id}

---
## How to run ?

```bash
go run main.go
```
