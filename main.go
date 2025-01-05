package main

import (
	"fmt"
	"go-login-register/database"
	"go-login-register/handlers"
	"net/http"
)

func main() {
	// MongoDB bağlantısı
	database.Connect()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	handlers.UserCollection = database.Client.Database("go_login_register").Collection("users")

	http.HandleFunc("/", handlers.IndexHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/dashboard", handlers.DashboardHandler)
	http.HandleFunc("/upload", handlers.UploadHandler)
	http.HandleFunc("/logout", handlers.LogoutHandler)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
