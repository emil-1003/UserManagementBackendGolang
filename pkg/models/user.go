package models

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/emil-1003/UserManagementBackendGolang/pkg/database"
)

type User struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
	Role     Role   `json:"role"`
	Created  string `json:"created"`
	Login    string `json:"login"`
}

func GetUsers() ([]User, error) {
	var users []User

	rows, err := database.DB.Query(`
		SELECT users.*, role.name AS role_name 
		FROM users
		JOIN role ON users.role_id = role.id;
	`)
	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var u User
		var r Role

		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &r.ID, &u.Created, &u.Login, &r.Name); err != nil {
			return users, err
		}

		u.Role = r

		createdParsedTime, _ := time.Parse(time.RFC3339, u.Created)
		u.Created = createdParsedTime.Format("2006-01-02 15:04")

		loginParsedTime, _ := time.Parse(time.RFC3339, u.Login)
		u.Login = loginParsedTime.Format("2006-01-02 15:04")

		users = append(users, u)
	}

	return users, err
}

func GetUserByEmail(email string) (User, error) {
	var u User
	var r Role
	result := database.DB.QueryRow(`
		SELECT users.*, role.name AS role_name 
		FROM users
		JOIN role ON users.role_id = role.id
		WHERE users.email = ?;
	`, email)
	err := result.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &r.ID, &u.Created, &u.Login, &r.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return u, fmt.Errorf("invalid email or password")
		}
		return u, fmt.Errorf("failed to scan user: %w", err)
	}

	u.Role = r

	return u, err
}

func UpdateUserLastLogin(user User) error {
	_, err := database.DB.Exec(`
		UPDATE users
		SET last_login = NOW() WHERE email = ?
		AND password = ?
	`, user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("failed to update last_login for user: %w", err)
	}

	return err
}

func AuthenticateUserPassword(userPassword, bodyPassword string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(bodyPassword)); err != nil {
		return fmt.Errorf("invalid email or password: %w", err)
	}

	return nil
}

func CreateUser(name string, email string, hashedPassword []byte) error {
	_, err := database.DB.Exec(`
		INSERT INTO users (name, email, password, role_id)
		VALUES (?, ?, ?, 1)
	`, name, email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return err
}

func DeleteUser(id int) error {

	_, err := database.DB.Exec(`
		DELETE FROM users
		WHERE id = ?
	`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return err
}
