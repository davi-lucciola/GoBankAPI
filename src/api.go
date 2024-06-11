package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIError struct {
	Error string
}

type APIMessage struct {
	Message string
}

type APIServer struct {
	listenAddress string
	storage       Storage
}

// "*" indicates the pointer (address of memory)
func NewAPIServer(listenAddress string, storage Storage) *APIServer {
	// "&" returns the pointer of instance.
	return &APIServer{
		storage:       storage,
		listenAddress: listenAddress,
	}
}

func WriteJSON(w http.ResponseWriter, statusCode int, body any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(body)
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func (s *APIServer) Run() {
	// Main API Router
	router := mux.NewRouter()

	// Routes Definition
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleAccountById))

	// API Up
	log.Printf("API running on: http://localhost%s", s.listenAddress)

	http.ListenAndServe(s.listenAddress, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.getAccounts(w, r)
	case "POST":
		return s.createAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.getAccountById(w, r)
	case "DELETE":
		return s.deleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (s *APIServer) getAccounts(w http.ResponseWriter, r *http.Request) error {
	// account := NewAccount("Davi")
	accounts, err := s.storage.GetAccounts()

	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		return WriteJSON(w, http.StatusNoContent, nil)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) getAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		return fmt.Errorf("o identificador precisa ser um uuid")
	}

	account, err := s.storage.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) createAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.Name, createAccountReq.Cpf)
	err := s.storage.CreateAccount(account)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (s *APIServer) deleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		return fmt.Errorf("o identificador precisa ser um uuid")
	}

	if _, err := s.storage.GetAccountById(id); err != nil {
		return err
	}

	s.storage.DeleteAccount(id)
	return WriteJSON(w, http.StatusOK, APIMessage{Message: "Conta excluída com sucesso."})
}

func (s *APIServer) transfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
