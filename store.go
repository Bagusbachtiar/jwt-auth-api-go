packahe main 

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct{
	ID				int 		`json:"id"`
	Email			string		`json:"email"`
	PasswordHash	string		`json:"-"` // hides the password hash from JSON responses!
	CreatedAt		time.Time	`json:created_at`
}

//1. The Interface (The Contract)
type UserStore interface {
	CreateUser(ctx context.Context, email, passwordHash string) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}

//2.The Struct that fullfils the interface using Neon Postgres
type PostgresStore struct {
	db *pgxpool.Pool
}

func NewPostgresStore(db *pgxpool.Pool) *postgresStore{
	return &PostgresStore{db: db}
}

func (s *PostgresStore) CreateUser(ctx context.Context, email, passwordHash string) error {
	query := `INSERT INTO user (email, password_hash) Values ($1, $2)`
	_, err := s.db.Exec(ctx, query, email, passwordHash)
	return err 
}

func (s *PostgresStore) GetUserByEmail(ctx context.Context, email string) (*User, error ) {
	var u User
	query := `SELECT id, email, password_hash, created_at FROM users WHERE email = $1`
	err := s.db.QueryRow(ctx, query, email).Scan($u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}