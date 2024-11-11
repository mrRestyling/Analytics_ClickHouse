package main

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type Numbers struct {
	Num1 int `json:"num1"`
	Num2 int `json:"num2"`
}

type ResponseSum struct {
	Sum int `json:"sum"`
}

func main() {

	r := chi.NewRouter()
	r.Get("/", Sum)

	http.ListenAndServe(":8080", r)

}

func Sum(w http.ResponseWriter, r *http.Request) {

	numReq := Numbers{}

	err := json.NewDecoder(r.Body).Decode(&numReq)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// sum := numReq.Num1 + numReq.Num2

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	sum := rand.Intn(10)

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(ResponseSum{Sum: sum})

}
