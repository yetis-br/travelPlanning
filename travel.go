package models

import (
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

// Travel define an object at the API
type Travel struct {
	ID     string
	Title  string
	UserID string
}

// Travels define a collection of objects at the API
type Travels struct {
	sync.RWMutex
	Store map[string]*Travel
}

//GetAllTravels returns all travels registered
func (t *Travels) GetAllTravels(w rest.ResponseWriter, r *rest.Request) {
	t.RLock()
	travels := make([]Travel, len(t.Store))
	i := 0
	for _, travel := range t.Store {
		travels[i] = *travel
		i++
	}
	t.RUnlock()
	w.WriteJson(&travels)
}
