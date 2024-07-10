package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Egyszerűsített felhasználói adatok
var users = map[string]string{"user": "password"}
var tokens = map[string]string{}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/protected", authenticate(protectedHandler))
	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok || users[username] != password {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Token generálása (egyszerűsített példa)
	token := "Bearer " + username + "-token"
	tokens[username] = token

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		for _, v := range tokens {
			if v == "Bearer "+token {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
	}
}

func protectedHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Access to protected resource granted")
}
