package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccountByID(int) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type Postgresstore struct {
	db *sql.DB
}

func NewPostgresStore() (*Postgresstore, error) {
	// Connect to the PostgreSQL database
	connStr := "user=postgres dbname=postgres host=172.17.0.2 password=gobank sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("connected")
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Postgresstore{
		db: db,
	}, nil
}

func (s *Postgresstore) Init() error {
	fmt.Println("creating table ...")
	return s.CreateAccountTable()
}

func (s *Postgresstore) CreateAccountTable() error {
	query := `CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50),
		number BIGINT,  -- Change to VARCHAR or appropriate type
		balance BIGINT ,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Default to current timestamp
	)`

	fmt.Println("Creating account table...")
	_, err := s.db.Exec(query)
	if err != nil {
		fmt.Println("Error creating table:", err)
		return err
	}
	fmt.Println("Account table created successfully.")
	return nil
}

func (s *Postgresstore) CreateAccount(account *Account) error {

	query := `INSERT INTO account (id, first_name, last_name, number, balance, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	// Execute the query with values from the Account struct
	_, err := s.db.Exec(query, account.ID, account.FirstName, account.LastName, account.Number, account.Balance, account.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *Postgresstore) UpdateAccount(*Account) error {
	return nil
}

func (s *Postgresstore) DeleteAccount(id int) error {
	_, err := s.db.Query("delete from account where id=$1", id)
	return err
}

func (s *Postgresstore) GetAccountByID(id int) (*Account, error) {
	query := "SELECT * FROM account WHERE id=$1"

	row := s.db.QueryRow(query, id)
	account, err := scanAccount(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("account with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to retrieve account: %v", err)
	}

	return account, nil
}

func (s *Postgresstore) GetAccounts() ([]*Account, error) {
	rows, err := s.db.Query("SELECT * FROM account")
	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)
		account, err := scanAccounts(rows)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}
	return accounts, nil
}

// Helper function to scan a row into an Account struct

func scanAccount(row *sql.Row) (*Account, error) {
	account := &Account{}
	err := row.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}

// Helper function to scan a row into an Account struct
func scanAccounts(rows *sql.Rows) (*Account, error) {
	account := &Account{}
	err := rows.Scan(
		&account.ID,
		&account.FirstName,
		&account.LastName,
		&account.Number,
		&account.Balance,
		&account.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return account, nil
}
