package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthServer struct {
	store     UserStore // Depends on the INTERFACE!
	jwtSecret []byte
}

func (s *AuthServer) handleRegister(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password required", http.StatusBadRequest)
		return
	}

	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	if err := s.store.CreateUser(r.Context(), req.Email, string(hashedBytes)); err != nil {
		log.Printf("❌ Register DB error: %v", err)
		http.Error(w, "User already exists or database error", http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (s *AuthServer) handleLogin(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	user, err := s.store.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Compare the plain password with the stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := generateToken(user.Email, s.jwtSecret)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (s *AuthServer) handleProfile(w http.ResponseWriter, r *http.Request) {
	// Extract the authenticated user's email from context
	email, _ := r.Context().Value(userEmailKey).(string)
	user, err := s.store.GetUserByEmail(r.Context(), email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func main() {
	_ = godotenv.Load()
	dbPool, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Database connection error: %v", err)
	}
	defer dbPool.Close()

	// Plug PostgresStore into our UserStore interface
	pgStore := NewPostgresStore(dbPool)
	server := &AuthServer{
		store:     pgStore,
		jwtSecret: []byte(os.Getenv("JWT_SECRET")),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", server.handleRegister)
	mux.HandleFunc("POST /login", server.handleLogin)
	mux.HandleFunc("GET /profile", authMiddleware(server.jwtSecret, server.handleProfile))

	log.Println("🔐 JWT Auth API running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
