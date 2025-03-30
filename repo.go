package main

import (
	"sync"
)

type Todo struct {
	ID    string `json:"id" validate:"required,uuid"`
	Title string `json:"title" validate:"required"`
	Done  bool   `json:"done"`
}

type Repo interface {
	TodoGet(todo *Todo) (*Todo, error)
	TodoAll() ([]*Todo, error)
	TodoSet(todo *Todo) error
	TodoDel(todo *Todo) error
}

func NewRepo() Repo {
	return &RepoImpl{
		db: make(map[string]*Todo),
	}
}

type RepoImpl struct {
	mu sync.RWMutex
	db map[string]*Todo
}

func (r *RepoImpl) TodoGet(todo *Todo) (*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.db[todo.ID], nil
}

func (r *RepoImpl) TodoSet(todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.db[todo.ID] = todo
	return nil
}

func (r *RepoImpl) TodoDel(todo *Todo) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.db, todo.ID)
	return nil
}

func (r *RepoImpl) TodoAll() ([]*Todo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	all := make([]*Todo, 0, len(r.db))
	for _, todo := range r.db {
		all = append(all, todo)
	}
	return all, nil
}
