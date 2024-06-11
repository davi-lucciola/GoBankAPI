package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(uuid.UUID) error
	UpdateAccount(*Account) error
	GetAccountById(uuid.UUID) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connStr := "host=db user=postgres dbname=go_bank_api password=password sslmode=disable"

	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createTable()
}

func (s *PostgresStore) createTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS account (
			id UUID PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			cpf VARCHAR(11) NOT NULL UNIQUE,
			number SERIAL NOT NULL,
			balance INT NOT NULL,
			created_at TIMESTAMP
		)`

	_, err := s.db.Exec(query)
	return err
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := Account{}
	err := rows.Scan(
		&account.ID,
		&account.Name,
		&account.Cpf,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := "SELECT * FROM account"

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}
	for rows.Next() {
		account, err := scanIntoAccount(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (s *PostgresStore) GetAccountById(accountID uuid.UUID) (*Account, error) {
	query := "SELECT * FROM account a WHERE a.id = $1"
	rows, err := s.db.Query(query, accountID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf(`conta '%v' não encontrada`, accountID)
}

func (s *PostgresStore) CreateAccount(account *Account) error {
	query := `
		INSERT INTO account (id, name, cpf, balance, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := s.db.Query(
		query,
		account.ID,
		account.Name,
		account.Cpf,
		account.Balance,
		account.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) DeleteAccount(accountID uuid.UUID) error {
	query := "DELETE FROM account WHERE id = $1"

	_, err := s.db.Query(query, accountID)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	return nil
}
