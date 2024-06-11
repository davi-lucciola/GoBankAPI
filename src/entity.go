package main

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cpf       string    `json:"cpf"`
	Number    int64     `json:"number"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(name string, cpf string) *Account {
	return &Account{
		ID:        uuid.New(),
		Name:      name,
		Cpf:       cpf,
		CreatedAt: time.Now().UTC(),
	}
}
