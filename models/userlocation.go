package models

import (
	
	"gopkg.in/mgo.v2/bson"
)

type(
	Location struct {
    ID    bson.ObjectId `json:"id" bson:"_id,omitempty"`
    Address    string `json:"address"`
    City       string `json:"city"`
    Name  string `json:"name"`
    State string `json:"state"`
    Zip   string `json:"zip"`
    Coordinate struct {
        	Lat float64 `json:"lat"`
        	Lng float64 `json:"lng"`
    } `json:"coordinate"`
	} 

)
