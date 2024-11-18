package handlers

import (
	"database/sql"
	"encoding/json"
	"go-mysql-api/model"
	"net/http"
	"github.com/gorilla/mux"
)

func GetUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, name, email,phone,address FROM users")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var users []models.User
		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			users = append(users, user)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func CreateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if user.Name == "" || user.Email == "" || user.Phone == "" || user.Address == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "All fields are required",
			})
			return
		}

		queryCheck := "SELECT COUNT(*) FROM users WHERE email = ?"
		var count int
		if err := db.QueryRow(queryCheck, user.Email).Scan(&count); err != nil {
			http.Error(w, "Error checking email", http.StatusInternalServerError)
			
			return
		}

		if count > 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)

			json.NewEncoder(w).Encode(map[string]string{
				"error": "Email already exists",
			})

			return
		}


		query := "INSERT INTO users (name, email, phone, address) VALUES (?, ?, ?, ?)"
		result, err := db.Exec(query, user.Name, user.Email, user.Phone, user.Address)
		if err != nil {

			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		id, _ := result.LastInsertId()
		user.ID = int(id)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}

func UpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		vars := mux.Vars(r)
		id := vars["id"]

		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		query := "UPDATE users SET name = ?, email = ?, phone = ?, address = ? WHERE id = ?"
		_, err := db.Exec(query, user.Name, user.Email, user.Phone, user.Address, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User updated successfully"})
	}
}

func DeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		vars := mux.Vars(r)
		id := vars["id"]

		query := "DELETE FROM users WHERE id = ?"
		_, err := db.Exec(query, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
	}
}


