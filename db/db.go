package db

import (
	"fmt"
	"time"
)

type Location struct {
	Name          string   `json:"name"`
	Address       string   `json:"address"`
	LocationTypes []string `json:"location_type"`
	Rating        float32  `json:"rating"`
	Place         string   `json:"place"`
}

// Interface for creating a db
type Database interface {
	Get(string) (Location, error)
	GetAll() []Location
	GetAllByType(locationType LocationType, name string) []Location
	Add(Location)
	Cache(time.Time) bool
}

type LocationGeometry struct {
	Lat float64
	Lng float64
}

func (lg *LocationGeometry) String() string {
	return fmt.Sprintf("%.07f,%.07f", lg.Lat, lg.Lng)
}

var SergelTorget = Place{
	Name: "sergel torget",
	Coordinates: LocationGeometry{
		Lat: 59.3323119,
		Lng: 18.0638270,
	},
}

type LocationType int

func (lt LocationType) String() string {
	return [...]string{
		"accounting",
		"airport",
		"amusement_park",
		"aquarium",
		"art_gallery",
		"atm",
		"bakery",
		"bank",
		"bar",
		"beauty_salon",
		"bicycle_store",
		"book_store",
		"bowling_alley",
		"bus_station",
		"cafe",
		"campground",
		"car_dealer",
		"car_rental",
		"car_repair",
		"car_wash",
		"casino",
		"cemetery",
		"church",
		"city_hall",
		"clothing_store",
		"convenience_store",
		"courthouse",
		"dentist",
		"department_store",
		"doctor",
		"drugstore",
		"electrician",
		"electronics_store",
		"embassy",
		"fire_station",
		"florist",
		"funeral_home",
		"furniture_store",
		"gas_station",
		"grocery_or_supermarket",
		"gym",
		"hair_care",
		"hardware_store",
		"hindu_temple",
		"home_goods_store",
		"hospital",
		"insurance_agency",
		"jewelry_store",
		"laundry",
		"lawyer",
		"library",
		"light_rail_station",
		"liquor_store",
		"local_government_office",
		"locksmith",
		"lodging",
		"meal_delivery",
		"meal_takeaway",
		"mosque",
		"movie_rental",
		"movie_theater",
		"moving_company",
		"museum",
		"night_club",
		"painter",
		"park",
		"parking",
		"pet_store",
		"pharmacy",
		"physiotherapist",
		"plumber",
		"police",
		"post_office",
		"primary_school",
		"real_estate_agency",
		"restaurant",
		"roofing_contractor",
		"rv_park",
		"school",
		"secondary_school",
		"shoe_store",
		"shopping_mall",
		"spa",
		"stadium",
		"storage",
		"store",
		"subway_station",
		"supermarket",
		"synagogue",
		"taxi_stand",
		"tourist_attraction",
		"train_station",
		"transit_station",
		"travel_agency",
		"university",
		"veterinary_care",
		"zoo",
	}[lt]
}

const (
	Accounting LocationType = iota
	Airport
	AmusementPark
	Aquarium
	ArtGallery
	Atm
	Bakery
	Bank
	Bar
	BeautySalon
	BicycleStore
	BookStore
	BowlingAlley
	BusStation
	Cafe
	Campground
	CarDealer
	CarRental
	CarRepair
	CarWash
	Casino
	Cemetery
	Church
	CityHall
	ClothingStore
	ConvenienceStore
	Courthouse
	Dentist
	DepartmentStore
	Doctor
	Drugstore
	Electrician
	ElectronicsStore
	Embassy
	FireStation
	Florist
	Funeral_home
	FurnitureStore
	GasStation
	GroceryOrSupermarket
	Gym
	HairCare
	HardwareStore
	HinduTemple
	HomeGoods_store
	Hospital
	InsuranceAgency
	JewelryStore
	Laundry
	Lawyer
	Library
	LightRail_station
	LiquorStore
	LocalGovernmentOffice
	Locksmith
	Lodging
	MealDelivery
	MealTakeaway
	Mosque
	MovieRental
	MovieTheater
	MovingCompany
	Museum
	NightClub
	Painter
	Park
	Parking
	PetStore
	Pharmacy
	Physiotherapist
	Plumber
	Police
	PostOffice
	PrimarySchool
	RealEstateAgency
	Restaurant
	RoofingContractor
	RvPark
	School
	SecondarySchool
	ShoeStore
	ShoppingMall
	Spa
	Stadium
	Storage
	Store
	SubwayStation
	Supermarket
	Synagogue
	TaxiStand
	TouristAttraction
	TrainStation
	TransitStation
	TravelAgency
	University
	VeterinaryCare
	Zoo
)
