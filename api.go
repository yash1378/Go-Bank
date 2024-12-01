package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
)

var a int
var mu *sync.Mutex

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))
	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferreq := new(TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferreq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferreq)
}

func (s *APIServer) handleGetAccountByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		id := mux.Vars(r)["id"]
		fmt.Println(id)
		// Convert string to int
		num := Str_to_Int(id)

		account, err := s.store.GetAccountByID(num)
		if err != nil {
			return err
		}
		return WriteJSON(w, http.StatusOK, account)

	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var Obj CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&Obj); err != nil {
		return err
	}

	account := NewAccount(Obj.FirstName, Obj.LastName)
	fmt.Println(1)
	account.ID = a
	a = a + 1
	if err := s.store.CreateAccount(account); err != nil {
		fmt.Println(2)
		return err
	}
	fmt.Println(3)

	return WriteJSON(w, http.StatusOK, account)

}

func (s *APIServer) handleGetAll(w http.ResponseWriter, r *http.Request) error {
	return nil

}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Println(id)
	// Convert string to int
	num := Str_to_Int(id)

	if err := s.store.DeleteAccount(num); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": num})

}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	return nil

}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error
type apiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle the error here only
			WriteJSON(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}

func Str_to_Int(st string) int {
	num, err := strconv.Atoi(st)
	if err != nil {
		fmt.Println("Error converting string to int:", err)
		panic(err)
	}

	return num

}
