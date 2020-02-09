package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/felipeluna/google-maps-places-go/db"
)

type Result struct {
	Name            string   `json:"name"`
	Rating          float32  `json:"rating"`
	Types           []string `json:"types"`
	NumberOfRatings int      `json:"number_of_ratings"`
	Vicinity        string   `json:"vicinity"`
}

type ResultByName struct {
	FormattedAddress string `json:"formatted_address"`
	Name             string `json:"name"`
}

type CandidatesByName struct {
	Candidates []ResultByName `json:"candidates"`
}
type LocationMaps struct {
	Results []Result
}

func unmarshalJson(rawJson string) (LocationMaps, error) {
	bytes := []byte(rawJson)
	var lm LocationMaps
	err := json.Unmarshal(bytes, &lm)
	if err != nil {
		return LocationMaps{}, errors.New("couldn't unmarshal json")
	}
	return lm, nil
}
// Giving a place and a location type query google places api for items, add those into the database
func GetListOfLocations(locationType db.LocationType, place db.Place, apiKey string, database db.Database) error {
	// check when was the last request, if was recently (e.g: 1 min ago) check the
	// database. otherwise send a new request, and repopulate the db. Send a request
	// to google places api
	if len(database.GetAllByType(locationType, strings.ToLower(place.Name))) != 0 {
		cache := database.Cache(time.Now())
		if cache {
			log.Println("cache, not performing request")
			return nil
		}
	}

	url := generateGetPlacesNearbyURL(place.Coordinates, locationType, apiKey)
	log.Println("sending request to :" + url)
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code %d from requests", response.StatusCode)
	}
	// unmarshal the results to the LocationMaps
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	maps, err := unmarshalJson(string(body))
	if err != nil {
		return err
	}


	if len(maps.Results) < 1 {
		var any map[string]interface{}
		log.Println(string(body))
		err := json.Unmarshal(body, &any)
		if err != nil {
			return errors.New("couldn't find any results")
		}
		if v, ok := any["status"]; ok {
			return errors.New(v.(string))
		}

		return errors.New("couldn't find any results")
	}

	// create location objects and add to the database
	updateDatabase(maps.Results, database, place.Name)
	return nil
}

// TODO: make another query to get formatted_url
func updateDatabase(results []Result, database db.Database, name string) {
	for _, r := range results {
		_, err := database.Get(r.Name)
		// don't exist
		if err != nil {
			location := db.Location{
				Name:          r.Name,
				Address:       r.Vicinity,
				LocationTypes: r.Types,
				Rating:        r.Rating,
				Place:         name,
			}
			database.Add(location)
		}
	}
}

// get list of locations based on type and, 2000 radius and geometry location
func generateGetPlacesNearbyURL(geometry db.LocationGeometry, locationType db.LocationType, apiKey string) string {
	return fmt.Sprintf(`https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s&radius=2000&type=%s&key=%s`, geometry.String(), locationType.String(), apiKey)
}

// url to get formatted_address
func generateGetFindPlaceFromText(geometry db.LocationGeometry, name string, apiKey string) string {
	return strings.ReplaceAll(
		fmt.Sprintf("https://maps.googleapis.com/maps/api/place/findplacefromtext/json?input=%s&inputtype=textquery&locationbias=circle:2000@%s&fields=photos,formatted_address,name,rating,opening_hours,geometry&key=%s", name, geometry.String(), apiKey),
		" ", "%20")
}
