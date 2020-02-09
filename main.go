package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/felipeluna/google-maps-places-go/api"
	"github.com/felipeluna/google-maps-places-go/db"
)

func main() {
	apiKey := os.Getenv("API_KEY")
	port := "8080"
	if apiKey == "" {
		log.Fatal("Could not read API_KEY from system env, make sure you have an API_KEY set")
	}
	database := db.NewInMemory()
	placeDatabase := db.NewPlaceDatabase()

	//http.HandleFunc("/", api.HandleGet(db.BicycleStore, db.SergelTorget, apiKey, &database))
	http.HandleFunc("/", api.HandleGetLocationGetLocationType(placeDatabase, apiKey, &database))
	http.HandleFunc("/all", func(writer http.ResponseWriter, request *http.Request) {
		marshal, err := json.Marshal(database.GetAll())
		if err != nil {
			http.Error(writer,"could not serialize objects", http.StatusInternalServerError )
		}
		writer.Header().Add("content-type", "application/json")
		writer.WriteHeader(http.StatusOK)
		writer.Write(marshal)
	})

	log.Println("Starting server on port " + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
