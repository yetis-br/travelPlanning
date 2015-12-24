package models

import (
	"net/http"
	"time"

	"github.com/ant0ine/go-json-rest/rest"
	rethink "github.com/dancannon/gorethink"
)

// Place define an object at the API
type Place struct {
	ID      string    `gorethink:"id,omitempty"`
	Title   string    `gorethink:"title"`
	Status  string    `gorethink:"status"`
	Created time.Time `gorethink:"created"`
}

// Places define a collection of objects at the API
type Places struct {
	Conn *rethink.Session
}

//GetAllPlaces returns all Places for the loged user
func (t *Places) GetAllPlaces(w rest.ResponseWriter, r *rest.Request) {
	places := []Place{}

	//table := t.TableName
	// Fetch all the items from the database
	res, err := rethink.Table(tableName).OrderBy(rethink.Asc("Created")).Run(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = res.All(&places)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteJson(&places)
	res.Close()
}

//GetPlace returns an expecific Place defined by the ID
func (t *Places) GetPlace(w rest.ResponseWriter, r *rest.Request) {
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
	var place *Place
	res.One(&place)
	w.WriteJson(place)
	res.Close()
}

//PostPlace creates a new Place
func (t *Places) PostPlace(w rest.ResponseWriter, r *rest.Request) {
	place := Place{}
	err := r.DecodeJsonPayload(&place)
	place.Created = time.Now()
	result, err := rethink.Table(tableName).Insert(place).RunWrite(t.Conn)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	place.ID = result.GeneratedKeys[0]
	w.WriteHeader(http.StatusCreated)
	w.WriteJson(&place)
}
