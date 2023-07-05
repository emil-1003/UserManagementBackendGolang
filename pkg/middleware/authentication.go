package middleware

import (
	"net/http"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/authentication"
)

func AuthOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check if token is valid
		_, ok := authentication.GetToken(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func AdminOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Check if token is valid
		token, ok := authentication.GetToken(r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if token role is admin
		if !authentication.IsUserAdmin(token) {
			http.Error(w, "not admin", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}
