package db

import "errors"

type Place struct {
	Name string
	Coordinates LocationGeometry
}

type PlaceDatabase struct {
	Places map[string]Place
}

func NewPlaceDatabase() PlaceDatabase {
	return PlaceDatabase{Places: map[string]Place{}}
}

func (pdb *PlaceDatabase) Get (name string) (Place, error) {
	if p, ok := pdb.Places[name]; ok {
		return p, nil
	}
	return Place{}, errors.New("place doesn't exist")
}

func (pdb *PlaceDatabase) Add(place Place) {
	pdb.Places[place.Name] = place
}