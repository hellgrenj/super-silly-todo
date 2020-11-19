package logic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/hellgrenj/super-silly-todo/microservice/domain/models"
)

// Repository provides access to todolist repository.
type Repository interface {
	AddItem(item models.Item, listID int) (int64, error)
}
type service struct {
	r Repository
}

// Service provides this microservices features
type Service interface {
	AddItem(item models.Item, listID int) (int64, string, error)
}

// NewService creates a todo service with the necessary dependencies
func NewService(r Repository) Service {
	return &service{r}
}

func (s *service) AddItem(item models.Item, listID int) (int64, string, error) {
	r1 := make(chan APIResponse)
	r2 := make(chan APIResponse)
	r3 := make(chan APIResponse)
	go fetchInsightFromThirdPartyAPI(r1, item.Name)
	go fetchInsightFromThirdPartyAPI(r2, item.Name)
	go fetchInsightFromThirdPartyAPI(r3, item.Name)

	var r APIResponse

	select { // take the first one and move on....
	case r = <-r1:
		fmt.Println("received from r1", r)
	case r = <-r2:
		fmt.Println("received from r2", r)
	case r = <-r3:
		fmt.Println("received from r3", r)
	}

	fmt.Printf("\nFetched response for %v with the message %v (SleepTime %v)\n", item.Name, r.Message, r.SleepTime)
	if r.Message != "" {
		item.Name = fmt.Sprintf("%v (%v)", item.Name, r.Message)
	}

	createdID, err := s.r.AddItem(item, listID)
	return createdID, item.Name, err
}

// APIResponse is the struct for the resonse from the third party API
type APIResponse struct {
	SleepTime   string
	Message     string
	CurrentTime string
}

// with simulated response times (see SleepTime...)
func fetchInsightFromThirdPartyAPI(c chan APIResponse, itemName string) {
	var client = &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Get(fmt.Sprintf("https://johanhellgren.se/api/test?itemName=%v", itemName))
	if err != nil {
		panic(err) // add err channel
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var r APIResponse
	err2 := json.Unmarshal(body, &r)
	if err2 != nil {
		fmt.Println(err2)
		return
	}
	c <- r
}
