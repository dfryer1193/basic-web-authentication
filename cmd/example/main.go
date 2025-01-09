package main

import (
	"github.com/dfryer1193/basic-web-authentication/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/welcome", handlers.WelcomeHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
