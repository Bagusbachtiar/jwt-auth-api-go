package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// FakeStore fulfills the UserStore interface using a simple map in RAM
type FakeStore struct {
	users map[string]*User
}

func (f *FakeStore) CreateUser(ctx context.Context, email, passwordHash string) error {
	if _, exists := f.users[email]; exists {
		return errors.New("user already exists")
	}
	f.users[email] = &User{
		ID:           1,
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
	}
	return nil
}

func (f *FakeStore) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	u, exists := f.users[email]
	if !exists {
		return nil, errors.New("not found")
	}
	return u, nil
}

func TestRegisterAndLoginWithFakeStore(t *testing.T) {
	fakeDB := &FakeStore{users: make(map[string]*User)}
	server := &AuthServer{
		store:     fakeDB,
		jwtSecret: []byte("test-secret-key"),
	}

	// 1. Test Registration
	payload := []byte(`{"email":"test@example.com","password":"password123"}`)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(payload))
	rec := httptest.NewRecorder()
	server.handleRegister(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status 201, got %d", rec.Code)
	}

	// 2. Test Login
	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(payload))
	loginRec := httptest.NewRecorder()
	server.handleLogin(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", loginRec.Code)
	}

	var resp map[string]string
	json.NewDecoder(loginRec.Body).Decode(&resp)
	if resp["token"] == "" {
		t.Errorf("Expected JWT token in login response, got empty string")
	}
}
