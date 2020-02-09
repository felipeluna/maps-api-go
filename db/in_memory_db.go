package db

import (
	"errors"
	"log"
	"time"
)

type InMemoryDatabase struct {
	places    map[string]Location
	cacheTime time.Time
}

func NewInMemory() InMemoryDatabase {
	return InMemoryDatabase{
		places:    map[string]Location{},
		cacheTime: time.Now(),
	}
}

// Return list of locations
func (db *InMemoryDatabase) GetAll() []Location {
	var locations []Location
	for _, v := range db.places {
		locations = append(locations, v)
	}
	return locations
}

// Get a location based on the name, if can't find return error
func (db *InMemoryDatabase) Get(name string) (Location, error) {
	if l, ok := db.places[name]; ok {
		return l, nil
	}
	return Location{}, errors.New("can't find location")
}

// If location with same name already exists do nothing.
// Otherwise, add it
func (db *InMemoryDatabase) Add(l Location) {
	if _, ok := db.places[l.Name]; !ok {
		log.Printf("inserting %s to database", l.Name)
		db.places[l.Name] = l
	}
}

// if the current request is 30 seconds or more after the last one
// don't use cache, otherwise do.
func (db *InMemoryDatabase) Cache(t time.Time) bool {
	if db.cacheTime.Add(30 * time.Second).After(t) {
		return true
	} else {
		db.cacheTime = t
		return false
	}
}

func (db *InMemoryDatabase) GetAllByType(locationType LocationType, name string) []Location {
	var filtered []Location

	for _, l := range db.GetAll() {
		if l.Place == name {
			if contains(l.LocationTypes, locationType.String()) {
				filtered = append(filtered, l)
			}
		}
	}
	return filtered
}

// check if array of string contains string
func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
