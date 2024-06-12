package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// API Server
type APIServer struct {
	listenAddress     string
	accountController AccountController
}

// "*" indicates the pointer (address of memory)
func NewAPIServer(listenAddress string, storage Storage) *APIServer {
	// "&" returns the pointer of instance.
	return &APIServer{
		listenAddress:     listenAddress,
		accountController: *NewAccountController(storage),
	}
}

func (api *APIServer) Run() {
	// Main API Router
	router := mux.NewRouter()

	api.accountController.Init(router)

	// API Up
	log.Printf("API running on: http://localhost%s", api.listenAddress)

	http.ListenAndServe(api.listenAddress, router)
}

// Utils
type APIError struct {
	Error string `json:"error"`
}

type APIMessage struct {
	Message string `json:"message"`
}

type APIToken struct {
	Token string `json:"token"`
	Type  string `json:"type"`
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
