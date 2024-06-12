package main

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	UserID    uuid.UUID `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(name string, cpf string, userID uuid.UUID) *Account {
	return &Account{
		ID:        uuid.New(),
		Name:      name,
		Cpf:       cpf,
		UserID:    userID,
		CreatedAt: time.Now().UTC(),
	}
}

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	AccountID uuid.UUID `json:"accountId"`
	CreatedAt time.Time `json:"createdAt"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return fmt.Sprintf("%x", bytes), err
}

func NewUser(username string, password string) (*User, error) {
	password, err := hashPassword(password)

	return &User{
		ID:        uuid.New(),
		Username:  username,
		Password:  password,
		CreatedAt: time.Now().UTC(),
	}, err
}

func (u *User) CheckPasswordHash(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func createTables(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username VARCHAR(100) NOT NULL UNIQUE,
			password VARCHAR(255) NOT NULL,
			created_at TIMESTAMP
		)
	`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("error to create users table: %v", err)
	}

	query = `
		CREATE TABLE IF NOT EXISTS accounts (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			number SERIAL NOT NULL,
			balance INT NOT NULL,
			user_id UUID NOT NULL,
			created_at TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users (id) 
		)`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("error to create accounts table: %v", err)
	}

	return nil
}
