package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hellgrenj/super-silly-todo/microservice/domain/logic"
	"github.com/hellgrenj/super-silly-todo/microservice/infrastructure/persistence"
	"github.com/hellgrenj/super-silly-todo/microservice/infrastructure/rest"
)

func main() {
	defer func() { //catch or finally
		if err := recover(); err != nil { //catch
			fmt.Fprintf(os.Stderr, "Exception: %v\n", err)
			os.Exit(1)
		}
	}()
	repo := persistence.NewMySQLRepo(1)
	service := logic.NewService(repo)
	server := rest.NewServer(service)
	log.Println("microservice up and running at port 5000")
	log.Fatal(http.ListenAndServe(":5000", server))
}
