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
	if err := json.NewDecoder
}