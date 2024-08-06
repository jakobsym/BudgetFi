package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jakobsym/BudgetFi/api/pkg/model"
	"github.com/joho/godotenv"
)

type Repository struct {
	sync.RWMutex
}

type uuid string

func New() *Repository {
	return &Repository{}
}

//TODO: GET

//TODO: DELETE

//TODO: PUT

// Creates a user upon registration
func (r *Repository) Post(_ context.Context, user *model.User) error {
	r.Lock()
	defer r.Unlock()
	db, err := MySqlConnection()
	defer db.Close()
	if err != nil {
		// err
	}
	//stmtIns, err := db.Prepare("INSERT INTO ")
	/*
		//TODO: This is an example
		// Prepare statement for inserting data
		stmtIns, err := db.Prepare("INSERT INTO squareNum VALUES( ?, ? )") // ? = placeholder
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer stmtIns.Close() // Close the statement when we leave main() / the program terminates
	*/
	return nil
}

// Util
// TODO: Set 'important settings' such as max # connectins, connection lifetime
func MySqlConnection() (*sql.DB, error) {
	db, err := sql.Open("mysql", loadDbEnv())
	if err != nil {
		// err
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		// err
	}
	fmt.Println("DB connected")
	return db, err
}

// Util
func loadDbEnv() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	mysqlCreds := os.Getenv("MYSQL_CREDENTIALS")
	return mysqlCreds
}
