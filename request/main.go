package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	tasks = make(map[int]Task)
	counter = 1
)

type Task struct {
	Text string `josn:"text"`
}

type addResponse struct{
	TaskId int`josn:"id"`
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	resMess := "alive"

	err := json.NewEncoder(w).Encode(resMess)

	if err != nil {
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	}
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fmt.Println("add")

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var task Task

	json.Unmarshal(body, &task)

	tasks[counter] = task

	err = json.NewEncoder(w).Encode(addResponse{TaskId:counter})
	counter++

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func main() {
	fmt.Println("Hello")

	r := mux.NewRouter()

	// check liveness
	r.HandleFunc("/healthz", handleHealthz).Methods("GET")

	// add data in list
	r.HandleFunc("/add", handleAdd).Methods("POST")

	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":1234", handler))
}
