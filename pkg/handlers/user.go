package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/authentication"
	"github.com/emil-1003/UserManagementBackendGolang/pkg/models"
)

func Signup() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var body struct {
			Name     string
			Email    string
			Password string
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, fmt.Errorf("failed to read body: %w", err).Error(), http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), 14)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to hash password: %w", err).Error(), http.StatusBadRequest)
			return
		}

		if err = models.CreateUser(body.Name, body.Email, hashedPassword); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte("user was registered successfully"))
	}
}

func Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var body struct {
			Email    string
			Password string
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			fmt.Println(fmt.Errorf("failed to read body: %w", err).Error())
			http.Error(w, fmt.Errorf("failed to read body: %w", err).Error(), http.StatusBadRequest)
			return
		}

		user, err := models.GetUserByEmail(body.Email)
		if err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := models.AuthenticateUserPassword(user.Password, body.Password); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenString := authentication.CreateJwt(os.Getenv("TOKEN_SECRET_WORD"), user)

		if err := models.UpdateUserLastLogin(user); err != nil {
			fmt.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Write([]byte(tokenString))
	}
}

func GetUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		users, err := models.GetUsers()
		if err != nil {
			http.Error(w, fmt.Errorf("failed to get users: %w", err).Error(), http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(users)
	}
}

func DeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		idString := mux.Vars(r)["id"]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, fmt.Errorf("invalid id: %w", err).Error(), http.StatusBadRequest)
			return
		}

		err = models.DeleteUser(id)
		if err != nil {
			http.Error(w, fmt.Errorf("failed to delete user: %w", err).Error(), http.StatusNotFound)
			return
		}

		w.Write([]byte("user was successfully deleted"))
	}
}
