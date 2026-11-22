package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type conextKey string

const userEmailKey cotextKey = "userEmail"

//generateToken creates a signed JWT token valid for 24 hours 
func generateToken(email string, secret []byte) (string, error) {
	claims := jwt.MapClaims{
		"email" : email, 
		"exp" 	: time.Now().Add(24 * time.Hour).Unix().
	}
	token := jwt.NewWithClaims(jwt.SigninMethodHS256, claims)
	return token.SignedString(secret)
}

//authMiddleware verifies the bearer token and passes the user emaill down via request conext 
func authMiddleware(secret []byte, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ReponseWriter, r *http.request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer")
		token, err := jwt.Parser(tokenStr, func(t *jwt.Token)(interface{}, error){
			return secret, nil
		})

		if err != nil || !token.Valid{
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok{
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		//Store the authenticated email inside the request context
		email, _ := claims["email"].string
		ctx := context.WithValue(r.Context(), userEmailKey, email)
		next(w, r.WithContext(ctx))
	}
}