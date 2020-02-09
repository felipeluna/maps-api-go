package db

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

// silly test to check if i implemented interfaces correctly
func InterfaceTesterDatabase(database Database) {
	fmt.Println("this compiles")
}

func TestInterfaceCompiles(t *testing.T) {
	db := new(InMemoryDatabase)
	InterfaceTesterDatabase(db)
}

func TestInMemoryDatabase_Add(t *testing.T) {
	l, db := createLocationAndDb()
	db.Add(l)
	if len(db.places) < 1 {
		t.Errorf("location not added to the db")
	}
}

func TestInMemoryDatabase_AddDuplicates(t *testing.T) {
	l, db := createLocationAndDb()
	// add the same entry twice
	db.Add(l)
	db.Add(l)
	if len(db.places) > 1 {
		t.Errorf("should contain only one entry")
	}
}

func TestDatabase_Get(t *testing.T) {
	shopName := "Second Hand Bike"
	shopAddress := "Riddargatan 29, 114 57 Stockholm"
	l := Location{
		Name:          shopName,
		Address:       shopAddress,
		LocationTypes: []string{BicycleStore.String()},
		Rating:        1.5,
	}
	db := NewInMemory()
	db.Add(l)
	location, err := db.Get(shopName)

	if err != nil {
		// no point in continue if is no location found
		t.Fatalf(shopName + " should exists in the database")
	}

	// check if returned location is not empty
	if reflect.DeepEqual(Location{}, location) {
		t.Errorf("returned empty location while should have been %v", l)
	}

	if location.Name != shopName && location.Address != shopAddress {
		t.Errorf("shop returned was different than expected, expected: %v, got %v", l, location)
	}
}

func TestInMemoryDatabase_GetError(t *testing.T) {
	l, db := createLocationAndDb()
	db.Add(l)
	location, err := db.Get("nonexistent location")
	if err == nil {
		t.Errorf("should not return any location")
	}
	if !reflect.DeepEqual(Location{}, location) {
		t.Errorf("should return empty location")
	}
}

func TestInMemoryDatabase_GetAll(t *testing.T) {
	l, db := createLocationAndDb()
	l2 := Location{
		Name:    "Elcykelbutik - EcoRide Stockholm City",
		Address: "Kungsgatan 72, 111 22 Stockholm",
	}

	db.Add(l)
	db.Add(l2)
	locations := db.GetAll()
	// check if locations length is 2
	if len(locations) != 2 {
		t.Errorf("locations slice should have 2 locations")
	}
	//check if locations are not empty struct
	for _, l := range locations {
		if reflect.DeepEqual(Location{}, l) {
			t.Errorf("should not be and empty location")
		}
	}
}

func TestLocationGeometry_toString(t *testing.T) {
	expected := "59.3323119,18.0638270"
	got := SergelTorget.Coordinates.String()
	if got != expected {
		t.Errorf("location geometry to String failed, expected: %v, got %v", expected, got)
	}
}

func TestInMemoryDatabase_Cache(t *testing.T) {
	_, database := createLocationAndDb()
	t.Run("cache true", func(t *testing.T) {
		cache := database.Cache(database.cacheTime.Add(time.Second * 15))
		if !cache {
			t.Errorf("cache should be true")
		}
	})
	t.Run("cache false", func(t *testing.T) {
		cache := database.Cache(database.cacheTime.Add(time.Second * 31))
		if cache {
			t.Errorf("cache should be false")
		}
	})

}

func TestInMemoryDatabase_Get(t *testing.T) {
	l, db := createLocationAndDb()
	db.Add(l)
	location, err := db.Get(l.Name)
	if err != nil {
		t.Errorf("unexpected error : %v", err)
	}
	if !reflect.DeepEqual(location, l) {
		t.Errorf("should have returned %v but got %v", location, l)
	}
}


func TestInMemoryDatabase_GetAllByType(t *testing.T) {
	l, db := createLocationAndDb()
	db.Add(l)
	locations := db.GetAllByType(BicycleStore,"sergel torget" )
	if len(locations) < 1 {
		t.Errorf("should have an item")
	}
	locations = db.GetAllByType(Restaurant,"sergel torget")
	if len(locations) < 1 {
		t.Errorf("should have an item")
	}
	locations = db.GetAllByType(Airport,"sergel torget")
	if len(locations) > 0 {
		t.Errorf("should not have an item")
	}
	locations = db.GetAllByType(Restaurant,"mall of scandinavia")
	if len(locations) > 0 {
		t.Errorf("should not have an item")
	}
}

// helper function to create a mock db and a mock location
func createLocationAndDb() (Location, InMemoryDatabase) {
	shopName := "Second Hand Bike"
	shopAddress := "Riddargatan 29, 114 57 Stockholm"
	l := Location{
		Name:    shopName,
		Address: shopAddress,
		LocationTypes: []string{"restaurant", "bicycle_store"},
		Place: "sergel torget",
	}
	db := NewInMemory()
	return l, db
}
