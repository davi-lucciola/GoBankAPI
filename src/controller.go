package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Account Controller
type AccountController struct {
	storage Storage
}

func NewAccountController(storage Storage) *AccountController {
	return &AccountController{
		storage: storage,
	}
}

func (ac *AccountController) Init(router *mux.Router) {
	// Routes Definition
	router.HandleFunc("/account", makeHTTPHandleFunc(ac.handleAccount))
	router.HandleFunc("/account/{id}", authorizationMiddleware(makeHTTPHandleFunc(ac.handleAccountById)))
	router.HandleFunc("/transfer", authorizationMiddleware(makeHTTPHandleFunc(ac.transfer)))
	router.HandleFunc("/login", makeHTTPHandleFunc(ac.login))
}

func (ac *AccountController) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return ac.getAccounts(w, r)
	case "POST":
		return ac.createAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (ac *AccountController) handleAccountById(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return ac.getAccountById(w, r)
	case "DELETE":
		return ac.deleteAccount(w, r)
	default:
		return fmt.Errorf("method not allowed: %s", r.Method)
	}
}

func (ac *AccountController) getAccounts(w http.ResponseWriter, _ *http.Request) error {
	accounts, err := ac.storage.GetAccounts()

	if err != nil {
		return err
	}

	if len(accounts) == 0 {
		return WriteJSON(w, http.StatusNoContent, nil)
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (ac *AccountController) getAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		return fmt.Errorf("o identificador precisa ser um uuid")
	}

	accountId, err := uuid.Parse(r.Header.Get("x-account-id"))

	if err != nil {
		return fmt.Errorf("houve um erro no identificador do token")
	}

	if id != accountId {
		return fmt.Errorf("você não pode realizar essa operação")
	}

	account, err := ac.storage.GetAccountById(id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (ac *AccountController) createAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := CreateAccountRequest{}

	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}
	defer r.Body.Close()

	user, err := NewUser(createAccountReq.Username, createAccountReq.Password)

	if err != nil {
		return err
	}

	account := NewAccount(createAccountReq.Name, createAccountReq.Cpf, user.ID)
	account, err = ac.storage.CreateAccount(account, user)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, account)
}

func (ac *AccountController) deleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := uuid.Parse(mux.Vars(r)["id"])

	if err != nil {
		return fmt.Errorf("o identificador precisa ser um uuid")
	}

	accountId, err := uuid.Parse(r.Header.Get("x-account-id"))

	if err != nil {
		return fmt.Errorf("houve um erro no identificador do token")
	}

	if id != accountId {
		return fmt.Errorf("você não pode realizar essa operação")
	}

	if _, err := ac.storage.GetAccountById(id); err != nil {
		return err
	}

	ac.storage.DeleteAccount(id)
	return WriteJSON(w, http.StatusOK, APIMessage{Message: "Conta excluída com sucesso."})
}

func (ac *AccountController) transfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "PATCH" {
		return fmt.Errorf("method not allowed")
	}

	transferRequest := TransferRequest{}

	if err := json.NewDecoder(r.Body).Decode(&transferRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	accountId, err := uuid.Parse(r.Header.Get("x-account-id"))

	if err != nil {
		return fmt.Errorf("houve um erro no identificador do token")
	}

	accountRequest, err := ac.storage.GetAccountById(accountId)

	if err != nil {
		return err
	}

	if accountId != transferRequest.ToAccountId {
		if accountRequest.Balance < transferRequest.Amount {
			return fmt.Errorf("você não possui saldo suficiente para fazer essa operacao")
		}

		accountTo, err := ac.storage.GetAccountById(transferRequest.ToAccountId)

		if err != nil {
			return err
		}

		accountTo.Balance += transferRequest.Amount
		accountRequest.Balance -= transferRequest.Amount

		err = ac.storage.UpdateAccount(accountTo)

		if err != nil {
			return err
		}
	} else {
		accountRequest.Balance += transferRequest.Amount
	}

	err = ac.storage.UpdateAccount(accountRequest)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, APIMessage{
		Message: fmt.Sprintf("Valor %d depositado com sucesso.", transferRequest.Amount),
	})
}

func (ac *AccountController) login(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return fmt.Errorf("method not allowed")
	}

	loginRequest := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		return err
	}
	defer r.Body.Close()

	user, err := ac.storage.GetUserByUsername(loginRequest.Username)

	if err != nil {
		return err
	}

	tokenString, err := createJWT(user.AccountID)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, APIToken{
		Token: tokenString,
		Type:  "x-jwt-token",
	})
}
