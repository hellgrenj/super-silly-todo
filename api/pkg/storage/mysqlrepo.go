package storage

import (
	"errors"
	"fmt"
	"os"
	"time"

	mysqlErrors "github.com/go-mysql/errors"

	// used by sqlx library
	_ "github.com/go-sql-driver/mysql"
	"github.com/hellgrenj/super-silly-todo/api/pkg/domainerrors"
	"github.com/hellgrenj/super-silly-todo/api/pkg/todo"
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

// AddList adds a new todolist in the database and returns the id for the created list or an error
func (m *MySQL) AddList(list todo.List) (int64, error) {

	result, err := m.db.Exec("INSERT INTO list (name) VALUES (?)", list.Name)
	if err != nil {
		if ok, myerr := mysqlErrors.Error(err); ok {

			if myerr == mysqlErrors.ErrDupeKey {
				expl := fmt.Sprintf("A list with name %s already exists", list.Name)
				return 0, domainerrors.NewConflictError(expl, err)
			}
			return 0, err
		}

		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

// AddItem adds an item in the database (for the todo list with the provided listID) and returns the created item id or an error
func (m *MySQL) AddItem(item todo.Item, listID int) (int64, error) {
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

// GetAllLists gets all todo lists (with items) from database
func (m *MySQL) GetAllLists() ([]todo.List, error) {

	rows, err := m.db.Queryx(`SELECT list.*, COALESCE(item.id, 0), COALESCE(item.name, ''), COALESCE(item.done, false) FROM list LEFT JOIN item ON (item.list_id = list.id)`)
	if err != nil {
		return nil, err
	}
	allLists, err := sqlxRowsToTodoLists(*rows)
	if err != nil {
		return nil, err
	}
	return allLists, err

	//return nil, errors.New("Not implemented")
}

// GetListByID finds a list by id
func (m *MySQL) GetListByID(id int) (todo.List, error) {
	rows, err := m.db.Queryx("SELECT list.id, list.name, COALESCE (item.id, 0), COALESCE (item.name, ''), COALESCE (item.done, false) FROM list LEFT JOIN item ON (item.list_id = list.id) WHERE list.id = ?", id)
	if err != nil {
		return todo.List{}, err
	}
	lists, err := sqlxRowsToTodoLists(*rows)
	if err != nil {
		return todo.List{}, err
	}
	if len(lists) == 0 {
		return todo.List{}, domainerrors.NewNotFoundError(fmt.Sprintf("list with id %v", id), errors.New("list not found"))
	}
	return lists[0], err

}

// DeleteListByID removes a list by id
func (m *MySQL) DeleteListByID(id int) error {
	rows, err := m.db.Queryx("DELETE FROM list WHERE list.id = ? ", id)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}

// DeleteItemByID removes an item by id
func (m *MySQL) DeleteItemByID(id int) error {
	rows, err := m.db.Queryx("DELETE FROM item WHERE item.id = ? ", id)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil
}

// UpdateItemDone updates a single todolist items done field
func (m *MySQL) UpdateItemDone(id int, done bool) error {
	rows, err := m.db.Queryx("UPDATE item SET done = ? WHERE item.id = ?", done, id)
	defer rows.Close()
	if err != nil {
		return err
	}
	return nil

}
func sqlxRowsToTodoLists(rows sqlx.Rows) ([]todo.List, error) {
	var allLists []todo.List
	var dataset = make(map[int]*todo.List)

	for rows.Next() {
		l := todo.List{}
		i := todo.Item{}
		err := rows.Scan(&l.ID, &l.Name, &i.ID, &i.Name, &i.Done)
		if err != nil {
			return nil, err
		}

		if _, ok := dataset[l.ID]; ok {
			if itemNotNull(i) {
				dataset[l.ID].Items = append(dataset[l.ID].Items, i)
			}

		} else {
			if itemNotNull(i) {
				l.Items = append(l.Items, i)
			}
			dataset[l.ID] = &l
		}

	}
	// map to list
	for _, v := range dataset {
		allLists = append(allLists, *v)
	}
	return allLists, nil
}
