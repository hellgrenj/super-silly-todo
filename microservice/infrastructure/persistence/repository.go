package persistence

import (
	"fmt"
	"os"
	"time"

	// used by sqlx library
	_ "github.com/go-sql-driver/mysql"
	"github.com/hellgrenj/super-silly-todo/microservice/domain/models"
	"github.com/jmoiron/sqlx"
)

// MySQL is the struct for the database
type MySQL struct {
	db *sqlx.DB
}

// NewMySQLRepo sets up connection pool (retries connection max 5 times over 20 seconds, depending on what you pass in as connectionAttempt)
func NewMySQLRepo(connectionAttempt int) *MySQL {
	connectionString := os.Getenv("MYSQL_CONNECTION_STRING")
	if len(connectionString) == 0 {
		connectionString = "mysql://root:example@localhost/silly"
	}

	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to connect to database: %v\n", err)
		if connectionAttempt < 5 {
			connectionAttempt++
			fmt.Printf("Trying again in 4 seconds attempt %v of 5\n", connectionAttempt)
			time.Sleep(4 * time.Second)
			return NewMySQLRepo(connectionAttempt)
		}
		os.Exit(1)
	}

	return &MySQL{
		db: db,
	}
}

// AddItem adds an item in the database (for the todo list with the provided listID) and returns the created item id or an error
func (m *MySQL) AddItem(item models.Item, listID int) (int64, error) {
	result, err := m.db.Exec("INSERT INTO item (name, done, list_id) VALUES (?, ?, ?)", item.Name, item.Done, listID)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}
