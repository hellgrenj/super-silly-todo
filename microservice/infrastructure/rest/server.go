package rest

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hellgrenj/super-silly-todo/microservice/domain/logic"
	"github.com/hellgrenj/super-silly-todo/microservice/domain/models"
)

// Server is the http server struct
type Server struct {
	router  *mux.Router
	service logic.Service
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	enableCors(&w) // TODO write as a midleware instead?
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	s.router.ServeHTTP(w, r)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT, PATCH")
	(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Access-Control-Allow-Methods, Authorization, X-Requested-With")
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v models.Ok) error {
	// TODO improve this function further, see https://www.alexedwards.net/blog/how-to-properly-parse-a-json-request-body

	r.Body = http.MaxBytesReader(w, r.Body, 1048576) // request body max 1 MB
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(v); err != nil {
		return err
	}
	return v.OK()
}

type operationResult struct {
	Result string `json:"result"`
}
type createdItemResult struct {
	ID   int64
	Name string
}

func (s *Server) addItemToList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	listID, err := strconv.Atoi(vars["listID"])
	if err != nil {
		json.NewEncoder(w).Encode(&operationResult{Result: "invalid list id"})
		return
	}
	var item models.Item
	if err := s.decode(w, r, &item); err != nil {
		fmt.Fprintf(w, "Something went wrong %v", err)
		return
	}
	fmt.Printf("Received item %v for list id %v\n", item, listID)
	id, name, err := s.service.AddItem(item, listID)
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(&createdItemResult{ID: id, Name: name})

}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(500)
	log.Println(err.Error())
	log.Fatal("Something went wrong!")
}

// NewServer returns a new http server
func NewServer(service logic.Service) *Server {
	s := &Server{router: mux.NewRouter(), service: service}
	s.router.HandleFunc("/todolist/{listID}/item", s.addItemToList).Methods("POST")
	return s
}
