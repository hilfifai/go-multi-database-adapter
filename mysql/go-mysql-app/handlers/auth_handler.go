package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"go-postgres-app/db"
	"go-postgres-app/models"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var u models.User
	json.NewDecoder(r.Body).Decode(&u)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	_, err := db.DB.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)",
		u.Name, u.Email, string(hashedPassword)) // MySQL pakai tanda tanya (?)
	if err != nil {
		http.Error(w, "Register error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Registered"})
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	json.NewDecoder(r.Body).Decode(&req)

	var dbUser models.User
	var hashed string
	err := db.DB.QueryRow("SELECT id, name, email, password FROM users WHERE email = ?", req.Email).
		Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &hashed) // MySQL pakai ?

	if err == sql.ErrNoRows || bcrypt.CompareHashAndPassword([]byte(hashed), []byte(req.Password)) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Buat token JWT di sini
	claims := models.JWTClaims{
		UserID: dbUser.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("secret"))
	if err != nil {
		http.Error(w, "Token creation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": signedToken})
}
