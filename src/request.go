package main

import "github.com/google/uuid"

type CreateAccountRequest struct {
	Name     string `json:"name"`
	Cpf      string `json:"cpf"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type TransferRequest struct {
	ToAccountId uuid.UUID `json:"toAccountId"`
	Amount      int64     `json:"amount"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
