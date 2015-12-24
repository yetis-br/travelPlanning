package models

import (
	"net/http"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/ant0ine/go-json-rest/rest"
	rethink "github.com/dancannon/gorethink"
)

// Trip define an object at the API
type Trip struct {
	ID          string      `gorethink:"id,omitempty"`
	Title       string      `gorethink:"title"`
	Status      string      `gorethink:"status"`
	Created     time.Time   `gorethink:"created"`
	Updated     time.Time   `gorethink:"updated"`
	TotalPlaces int         `gorethink:"-" json:"totalPlaces"`
	Places      []TripPlace `gorethink:"places"`
}

// TripPlace define an place associated to the trip at the API
type TripPlace struct {
	ID     string `gorethink:"place_id" json:"place_id"`
	Status string `gorethink:"place_status" json:"place_status"`
	Order  int    `gorethink:"place_order" json:"place_order"`
}

// Trips define a collection of objects at the API
type Trips struct {
	Conn      *rethink.Session
	TableName string
}

const tableName = "Trip"

//GetAllTrips returns all Trips for the loged user
func (t *Trips) GetAllTrips(w rest.ResponseWriter, r *rest.Request) {
	trips := []Trip{}

	//table := t.TableName
	// Fetch all the items from the database
	res, err := rethink.Table(tableName).OrderBy(rethink.Asc("Created")).Run(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = res.All(&trips)

	//Include total places in the
	for _, trip := range trips {
		trip.TotalPlaces = len(trip.Places)
	}

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

	res, err := rethink.Table(tableName).Get(id).Run(t.Conn)
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
	trip.TotalPlaces = len(trip.Places)

	w.WriteJson(trip)
	res.Close()
}

//PostTrip creates a new Trip
func (t *Trips) PostTrip(w rest.ResponseWriter, r *rest.Request) {
	trip := Trip{}
	err := r.DecodeJsonPayload(&trip)
	trip.Created = time.Now()
	trip.Updated = time.Now()
	result, err := rethink.Table(tableName).Insert(trip).RunWrite(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trip.ID = result.GeneratedKeys[0]
	trip.TotalPlaces = len(trip.Places)

	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&trip)
}

//UpdateTripPlaces update trip places
func (t *Trips) UpdateTripPlaces(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	res, err := rethink.Table(tableName).Get(id).Run(t.Conn)
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

	tripPlaces := []TripPlace{}
	err = r.DecodeJsonPayload(&tripPlaces)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trip.Updated = time.Now()
	trip.Places = tripPlaces
	trip.TotalPlaces = len(tripPlaces)

	_, err = rethink.Table(tableName).Get(id).Update(trip).RunWrite(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&trip)
}

//DeleteTripPlaces delete trip places
func (t *Trips) DeleteTripPlaces(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	placeID := r.PathParam("place")
	index, _ := strconv.Atoi(placeID)

	_, err := rethink.Table(tableName).Get(id).Update(map[string]interface{}{
		"places": rethink.Row.Field("places").DeleteAt(index),
	}).RunWrite(t.Conn)

	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		log.Warning(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
