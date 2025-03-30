package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

var validate = validator.New()

func todosGet(repo Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		todo := &Todo{ID: id}
		if _, err := repo.TodoGet(todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

func todosGetAll(repo Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		all, err := repo.TodoAll()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(all)
	}
}

func todosPost(repo Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := repo.TodoSet(&todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(todo)
	}
}

func todosDel(repo Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		todo := &Todo{ID: id}
		if err := repo.TodoDel(todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

func todosPut(repo Repo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		todo := &Todo{ID: id}
		if err := json.NewDecoder(r.Body).Decode(todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := validate.Struct(todo); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if _, err := repo.TodoGet(todo); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if err := repo.TodoSet(todo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(todo)
	}
}

func main() {
	mux := http.NewServeMux()
	repo := NewRepo()
	mux.HandleFunc("GET /todos", todosGetAll(repo))
	mux.HandleFunc("GET /todos/{id}", todosGet(repo))
	mux.HandleFunc("POST /todos", todosPost(repo))
	mux.HandleFunc("DELETE /todos/{id}", todosDel(repo))
	mux.HandleFunc("PUT /todos/{id}", todosPut(repo))

	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
