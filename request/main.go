package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	tasks   = make(map[int]Task)
	counter = 1
)

type Task struct {
	Text string `josn:"text"`
}

type addResponse struct {
	TaskId int `josn:"id"`
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

	err = json.NewEncoder(w).Encode(addResponse{TaskId: counter})
	counter++

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func handleList(w http.ResponseWriter, r *http.Request) {

	listResponse := make(map[int]interface{})

	for i, task := range tasks {
		listResponse[i] = task.Text
	}

	fmt.Println(listResponse)

	err := json.NewEncoder(w).Encode(listResponse)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Println("delete")
	vars := mux.Vars(r)
	ou := vars["key"]
	fmt.Println(ou)
	originalURL, err := url.PathUnescape(ou)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idx, err := strconv.Atoi(originalURL)

	delete(tasks, idx)
	err = json.NewEncoder(w).Encode(originalURL)

}

func main() {
	fmt.Println("Hello")

	r := mux.NewRouter()

	// check liveness
	r.HandleFunc("/healthz", handleHealthz).Methods("GET")

	// add data in list
	r.HandleFunc("/add", handleAdd).Methods("POST")

	// list data in list
	r.HandleFunc("/list", handleList).Methods("GET")

	// delete data in list by key
	r.HandleFunc("/delete/{key}", handleDelete).Methods("POST")

	handler := cors.Default().Handler(r)
	log.Fatal(http.ListenAndServe(":1234", handler))
}
