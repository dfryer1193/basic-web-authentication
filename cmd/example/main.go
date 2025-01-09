package main

import (
	"github.com/dfryer1193/basic-web-authentication/handlers"
	"github.com/dfryer1193/basic-web-authentication/storage"
	"log"
	"net/http"
)

func main() {
	userStore := storage.NewInMemoryUserStore()
	authHandler := handlers.NewUserAwareHandler("example", userStore)

	http.HandleFunc("/register", authHandler.RegisterHandler)
	http.HandleFunc("/login", authHandler.LoginHandler)
	http.HandleFunc("/change-password", authHandler.PasswordUpdateHandler)
	http.HandleFunc("/welcome", authHandler.WelcomeHandler)

	// Serve static files
	fs := http.FileServer(http.Dir("./web/static"))
	http.Handle("/", fs)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
