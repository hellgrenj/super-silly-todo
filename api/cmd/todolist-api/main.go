package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hellgrenj/super-silly-todo/api/pkg/todo"

	"github.com/hellgrenj/super-silly-todo/api/pkg/storage"

	"github.com/hellgrenj/super-silly-todo/api/pkg/http/rest"
)

func main() {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	// first init db repo
	repository := storage.NewMySQLRepo(1)
	// then construct todo service (passing in db dependency)
	service := todo.NewService(repository) // repository qualifies as a todo repository .. meets interface contract
	// then constructing rest server (passing in service dependency)
	server := rest.NewServer(service)
	// then start rest service
	log.Fatal(http.ListenAndServe(":4000", server))
}
