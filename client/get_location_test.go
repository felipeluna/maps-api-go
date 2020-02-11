package client

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/felipeluna/google-maps-places-go/db"
)

func Test_UnmarshalJson(t *testing.T) {
	rawJson, err := ioutil.ReadFile("./example.json")
	if err != nil {
		t.Errorf("could not open example json. %v", err)
	}
	maps, err := unmarshalJson(string(rawJson))
	if err != nil {
		t.Errorf("error unmarshalling json %v", err)
	}
	if len(maps.Results) < 1 {
		t.Errorf("could not unrmashall properly, should have 20 items has 0")
	}
	// not sure how to test this
	// TODO
	for _, result := range maps.Results {
		if reflect.DeepEqual(result, Result{}) {
			t.Errorf("returned empty location while should have been filled")
		}
		fmt.Println(result.Name, result.Types, result.Vicinity, result.NumberOfRatings, result.Rating)
	}
}

func Test_UnmarshalJsonWrong(t *testing.T) {
	rawJson, err := ioutil.ReadFile("./example-bad.json")
	if err != nil {
		t.Errorf("could not open example json. %v", err)
	}
	_, err = unmarshalJson(string(rawJson))
	if err == nil {
		t.Errorf("should not unmarshal")
	}
}

func Test_generateGetPlacesNearbyURL(t *testing.T) {
	url := generateGetPlacesNearbyURL(db.SergelTorget.Coordinates, db.BicycleStore, "FakeApiString")
	expected := "https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=59.3323119,18.0638270&radius=2000&type=bicycle_store&key=FakeApiString"
	if url != expected {
		t.Errorf("wrong url, expected: %v, got %v", expected, url)
	}

}

func Test_getDetailedInfo(t *testing.T) {
	repo := db.NewInMemory()
	result := Result{
		Name:            "A store",
		Rating:          1.0,
		Types:           []string{"bicycle_store"},
		NumberOfRatings: 2,
		Vicinity:        "Street One",
	}
	results := LocationMaps{Results: []Result{result}}
	updateDatabase(results.Results, &repo, "")
	if len(repo.GetAll()) < 1 {
		t.Errorf("no location was added to the repo")
	}
}

func TestGetListOfLocations(t *testing.T) {

}
