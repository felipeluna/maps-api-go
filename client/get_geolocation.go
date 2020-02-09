package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/felipeluna/google-maps-places-go/db"
)

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type LocationGeometryAPI struct {
	Location LatLng `json:"location"`
}

type GeometryAPI struct {
	Geometry LocationGeometryAPI `json:"geometry"`

}

type LocationAPICandidates struct {
	Candidates []GeometryAPI `json:"candidates"`
}

// TODO: create test with external dependencies
// Send request to google places api and return a db.Place with coordinates
// if not found return error
func GetCoordinates(findGeometryUrl string, name string) (db.Place, error) {
	log.Printf("sending request to %s\n", findGeometryUrl)
	// send request to return the coordinates of a given place
	response, err := http.Get(findGeometryUrl)
	if err != nil {
		log.Printf("error getting %s. : %v", findGeometryUrl, err)
		return db.Place{}, errors.New("couldn't find this place")
	}
	defer response.Body.Close()
	// unmarshal body to API struct
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("error reading body from %s: %v", findGeometryUrl, err)
		return db.Place{}, errors.New("error retrieving place coordinates")
	}
	lac := new(LocationAPICandidates)
	err = json.Unmarshal(body, &lac)
	if err != nil {
		log.Printf("error unmarshalling body from %s: %v", findGeometryUrl, err)
		return db.Place{}, errors.New("error retrieving place coordinates")
	}
	// if there is an empty candidates array something went wrong
	if len(lac.Candidates) < 1 {
		return db.Place{}, errors.New("error retrieving place coordinates")
	}
	locationFromAPI := lac.Candidates[0].Geometry.Location
	// create db.Place based on API results
	place := db.Place{
		Name: name,
		Coordinates: db.LocationGeometry{
			Lat: locationFromAPI.Lat,
			Lng: locationFromAPI.Lng,
		},
	}
	return place, nil
}

// Generate url for querying the place with coordinates
func GenerateFindGeometryUrl(name string, apiKey string) string {
	escapedString := url.QueryEscape(name)
	return fmt.Sprintf("https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=%s&inputtype=textquery&fields=geometry&key=%s", escapedString, apiKey)
}
