package main

import (
	"fmt"
	"net/http"

	"github.com/zwzn/zeus"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title>Tests</title>
</head>
<body>
	Hello World
</body>
</html>`)
}

type User struct {
	ID    int    `storm:"id"`     // primary key
	Group string `storm:"index"`  // this field will be indexed
	Email string `storm:"unique"` // this field will be indexed with a unique constraint
	Name  string // this field will not be indexed
	Age   int    `storm:"index"`
}

func main() {
	z, _ := zeus.Open(
		&User{},
	)
	http.HandleFunc("/", greet)
	http.HandleFunc("/api", z.Handle())
	http.ListenAndServe(":8080", nil)
}
