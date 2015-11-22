package main

import (

	"github.com/drone/routes"
    "net/http"
    "controllers"
    "gopkg.in/mgo.v2"
)

func main() {
    mux := routes.New()

    lc := controllers.NewLocationController(getSession())
    uc := controllers.NewUberController(getSession())

    mux.Get("/locations/:location_id", lc.GetLocation)
    mux.Post("/locations", lc.AddLocation)
    mux.Del("/locations/:location_id", lc.DeleteLocation)
    mux.Put("/locations/:location_id", lc.UpdateLocation)

    mux.Post("/trips", uc.AddTrip)
    mux.Get("/trips/:trip_id", uc.GetTrip)
    mux.Put("/trips/:trip_id/request", uc.RequestTrip)
    http.Handle("/", mux)
    http.ListenAndServe(":8080", nil)
}

func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://user:user@ds041144.mongolab.com:41144/test_mongo_db")

	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}

	// Deliver session
	return s
}
