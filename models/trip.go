package models

import (
	"net/http"
	"sync"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	rethink "github.com/dancannon/gorethink"
)

// Trip define an object at the API
type Trip struct {
	ID      string `gorethink:"id,omitempty"`
	Title   string
	Status  string
	Created time.Time
}

// Trips define a collection of objects at the API
type Trips struct {
	sync.RWMutex
	Conn *rethink.Session
}

//GetAllTrips returns all Trips for the loged user
func (t *Trips) GetAllTrips(w rest.ResponseWriter, r *rest.Request) {
	trips := []Trip{}
	// Fetch all the items from the database
	res, err := rethink.Table("Trip").OrderBy(rethink.Asc("Created")).Run(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = res.All(&trips)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&trips)
	res.Close()
}

//GetTrip returns an expecific Trip defined by the ID
func (t *Trips) GetTrip(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")

	res, err := rethink.Table("Trip").Get(id).Run(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if res.IsNil() {
		rest.NotFound(w, r)
		return
	}
	var trip *Trip
	res.One(&trip)
	w.WriteJson(trip)
	res.Close()
}

//PostTrip creates a new Trip
func (t *Trips) PostTrip(w rest.ResponseWriter, r *rest.Request) {
	trip := Trip{}
	err := r.DecodeJsonPayload(&trip)
	trip.Created = time.Now()
	result, err := rethink.Table("Trip").Insert(trip).RunWrite(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trip.ID = result.GeneratedKeys[0]
	w.WriteJson(&trip)
}
