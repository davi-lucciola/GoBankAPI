package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account, *User) (*Account, error)
	DeleteAccount(uuid.UUID) error
	UpdateAccount(*Account) error
	GetAccountById(uuid.UUID) (*Account, error)
	GetAccounts() ([]*Account, error)
	GetUserByUsername(string) (*User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStorage() (*PostgresStore, error) {
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
	return createTables(s.db)
}

func scanIntoAccount(rows *sql.Rows) (*Account, error) {
	account := Account{}
	err := rows.Scan(
		&account.ID,
		&account.Name,
		&account.Cpf,
		&account.Number,
		&account.Balance,
		&account.UserID,
		&account.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, err
}

func scanIntoUser(rows *sql.Rows) (*User, error) {
	user := User{}
	err := rows.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.AccountID,
	)

	if err != nil {
		return nil, err
	}

	return &user, err
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := "SELECT * FROM accounts"

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
	query := "SELECT * FROM accounts a WHERE a.id = $1"
	rows, err := s.db.Query(query, accountID)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoAccount(rows)
	}

	return nil, fmt.Errorf(`conta '%v' não encontrada`, accountID)
}

func (s *PostgresStore) GetUserByUsername(username string) (*User, error) {
	query := `
		SELECT u.*, a.id FROM accounts a
		INNER JOIN users u ON a.user_id = u.id 
		WHERE u.username = $1
	`

	rows, err := s.db.Query(query, username)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		return scanIntoUser(rows)
	}

	return nil, fmt.Errorf(`usuário '%v' não encontrado`, username)
}

func (s *PostgresStore) CreateAccount(account *Account, user *User) (*Account, error) {
	tran, err := s.db.Begin()

	if err != nil {
		return nil, err
	}

	// Create User
	query := `
		INSERT INTO users (id, username, password, created_at)
		VALUES ($1, $2, $3, $4)
	`

	rows, err := tran.Query(
		query,
		user.ID,
		user.Username,
		user.Password,
		user.CreatedAt,
	)
	rows.Close()

	if err != nil {
		tran.Rollback()
		return nil, err
	}

	// Create Account
	query = `
		INSERT INTO accounts (id, name, cpf, balance, user_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING *
	`

	rows, err = tran.Query(
		query,
		account.ID,
		account.Name,
		account.Cpf,
		account.Balance,
		account.UserID,
		account.CreatedAt,
	)

	if err != nil {
		tran.Rollback()
		return nil, err
	}

	for rows.Next() {
		account, err = scanIntoAccount(rows)
	}

	if err != nil {
		tran.Rollback()
		return nil, err
	}

	tran.Commit()
	return account, err
}

func (s *PostgresStore) DeleteAccount(accountID uuid.UUID) error {
	query := "DELETE FROM accounts WHERE id = $1"

	_, err := s.db.Query(query, accountID)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) UpdateAccount(account *Account) error {
	query := `
		UPDATE accounts
		SET name = $1, cpf = $2, balance = $3
		WHERE id = $4
	`

	_, err := s.db.Query(
		query,
		account.Name,
		account.Cpf,
		account.Balance,
		account.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
