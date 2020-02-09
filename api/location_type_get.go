package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/felipeluna/google-maps-places-go/client"
	"github.com/felipeluna/google-maps-places-go/db"
)

// Receives a request with ?place=string&location=string, checks for values in the database
// if can't find records, tries to get the coordinates of place and fires a request
// to return the list of location.
func HandleGetLocationGetLocationType(placeDatabase db.PlaceDatabase, apiKey string, database db.Database) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		place, done := validatePlaceQuery(r, w)
		if done {
			return
		}
		lt, done := validateLocationTypeQuery(r, w, place)
		if done {
			return
		}

		locations := database.GetAllByType(lt, place)
		// if there are items in the database return
		// google api free tier too much "OVER_QUERY_LIMIT"
		if len(locations) > 1 {
			rawJson, err := json.Marshal(locations)
			if err != nil {
				log.Println(fmt.Sprintf("error marshalling response  %v", err))
				w.Header().Add("Content-type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
			}
			w.Header().Add("Content-type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(rawJson)
		} else {
			if getListLocationFromRequest(placeDatabase, place, apiKey, w, lt, database) {
				return
			}

		}
	}
}

func getListLocationFromRequest(placeDatabase db.PlaceDatabase, place string, apiKey string, w http.ResponseWriter, lt db.LocationType, database db.Database) bool {
	// based on name get coordinates
	dbPlace, done := getPlace(placeDatabase, strings.ToLower(place), apiKey, w)
	// if true error happened and response was sent
	if done {
		return true
	}

	// get locations and fill database
	err := client.GetListOfLocations(lt, dbPlace, apiKey, database)
	if err != nil {
		http.Error(w, "Couldn't get list of locations", http.StatusInternalServerError)
		// don't leak the error to the user, log instead
		log.Printf("ERROR - %v", err)
		return true
	}

	// get items from database and return in json format
	locations := database.GetAllByType(lt, dbPlace.Name)
	rawJson, err := json.Marshal(locations)
	if err != nil {
		http.Error(w, "Couldn't get list of locations", http.StatusInternalServerError)
		// don't leak the error to the user, log instead
		log.Printf("ERROR - %v", err)
		return true
	}
	if rawJson == nil {
		rawJson = []byte("[]")
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(rawJson)
	return false
}
// Given request check if contains a valid locationType
// if not valid returns a Bad Request status code
// and helper message
func validateLocationTypeQuery(r *http.Request, w http.ResponseWriter, place string) (db.LocationType, bool) {
	locationType := r.URL.Query().Get("type")
	if len(locationType) < 1 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "url must contain type query parameters query/?type=restaurant; supported types : ['bicycle_store', 'restaurant']"}`))
		return 0, true
	}
	log.Printf("query for: %s %s", locationType, place)
	var lt db.LocationType
	if locationType == strings.ToLower(db.BicycleStore.String()) {
		lt = db.BicycleStore
	} else if locationType == strings.ToLower(db.Restaurant.String()) {
		lt = db.Restaurant
	} else if locationType == strings.ToLower(db.Store.String()) {
		lt = db.Store
	} else {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "url must contain type query parameters: e.g. query/?type=restaurant; supported types : ['bicycle_store', 'restaurant', 'store']" 2 2}`))
		return 0, true
	}
	return lt, false
}

// Given request check if contains a valid place
// if not valid returns a Bad Request status code
// and helper message
func validatePlaceQuery(r *http.Request, w http.ResponseWriter) (string, bool) {
	place := r.URL.Query().Get("place")
	if len(place) < 1 {
		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "url must contain place query parameters query/?place=a nice place"}`))
		return "", true
	}
	return place, false
}

// Return db.Place from database if exist otherwise
// query to the google places api
func getPlace(placeDatabase db.PlaceDatabase, place string, apiKey string, w http.ResponseWriter) (db.Place, bool) {
	// check if place exists
	dbPlace, err := placeDatabase.Get(place)
	// if doesn't exist try to fetch from google places api
	if err != nil {
		geometryUrl := client.GenerateFindGeometryUrl(place, apiKey)
		log.Printf("sending request to get api; url: %s\n", geometryUrl)
		coordinates, err := client.GetCoordinates(geometryUrl, place)
		if err != nil {
			log.Printf("error from request %s: %v",geometryUrl, err)
			w.Header().Add("content-type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf(`{"error": "couldn't locate %s'"}`, place)))
			// return true to signalize that response was sent
			return db.Place{}, true
		}
		dbPlace = coordinates
		// add place for future reference
		placeDatabase.Add(dbPlace)
	}
	// return false to signalize that response was not sent
	return dbPlace, false
}