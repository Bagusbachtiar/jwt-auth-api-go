package main 

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/vs/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email		string `json:"email"`
	Password 	string `json:"password"`
}

type AuthServer struct {
	store		UserStore//Depend on the Interface
	jwtSecret	[]byte
}

func (s *AuthServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err !=nil || req.Email == "" || req.Password == ""{
			http.Error(w, "Email and passowrd required", http.StatusBadRequest)
			return
	}
	
	//Hash the password securely with bcrypt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash Password", http.StatusInternalServerError)
		return
	}

	if err := s.store.CreateUser(r.Context(), req.Email, string(hashedBytes)); err != nil {
		http.Error(w, "User already exists or database error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User Registered Succesfully"})
}