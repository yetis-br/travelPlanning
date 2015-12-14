package models

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/ant0ine/go-json-rest/rest"
)

// Trip define an object at the API
type Trip struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	/*
		UserID    string
		StartDate time.Time
		EndDate   time.Time
	*/
}

// Trips define a collection of objects at the API
type Trips struct {
	sync.RWMutex
	Store map[string]*Trip
}

//GetAllTrips returns all Trips for the loged user
func (t *Trips) GetAllTrips(w rest.ResponseWriter, r *rest.Request) {
	t.RLock()
	trips := make([]Trip, len(t.Store))
	i := 0
	for _, trip := range t.Store {
		trips[i] = *trip
		i++
	}
	t.RUnlock()
	w.WriteJson(&trips)
}

//GetTrip returns an expecific Trip defined by the ID
func (t *Trips) GetTrip(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	t.RLock()
	var trip *Trip
	if t.Store[id] != nil {
		trip = &Trip{}
		*trip = *t.Store[id]
	}
	t.RUnlock()
	if trip == nil {
		rest.NotFound(w, r)
		return
	}
	w.WriteJson(trip)
}

//PostTrip creates a new Trip
func (t *Trips) PostTrip(w rest.ResponseWriter, r *rest.Request) {
	trip := Trip{}
	err := r.DecodeJsonPayload(&trip)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Lock()
	id := fmt.Sprintf("%d", len(t.Store)) // stupid
	trip.ID = id
	t.Store[id] = &trip
	t.Unlock()
	w.WriteJson(&trip)
}
