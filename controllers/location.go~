package controllers

import (
 
    "fmt"
    "strings"
    "encoding/json"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "github.com/DhruvKalaria/cmpe273-Assignment2/tree/master/models"
    "net/http"
    "time"
)

type (

	LocationController struct	{

			session *mgo.Session
	}

)

const(
    timeout = time.Duration(time.Second*100)
)

func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}


func GetCoordinates(address string) (float64, float64)  {

    googloc := models.GoogleLocation{}
    url := "http://maps.google.com/maps/api/geocode/json?address=" + address
    url = getFormattedURL(url)
    fmt.Println(url)
    client := http.Client{Timeout: timeout}
    res, err := client.Get(url)
    if err != nil {
        fmt.Errorf("Cannot read Google API: %v", err)
    }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&googloc)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }
    return googloc.Results[0].Geometry.Location.Lat, googloc.Results[0].Geometry.Location.Lng

}

func (lc LocationController) AddLocation(w http.ResponseWriter, r *http.Request) {

    loc := models.Location{}   
    
    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&loc)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Input JSON: %v", err)
    }

    var address string
    address = loc.Address+loc.City+loc.State+loc.Zip
   	lat, lng := GetCoordinates(address)

   	loc.ID = bson.NewObjectId()
 	  loc.Coordinate.Lat=lat
 	  loc.Coordinate.Lng=lng
    lc.session.DB("test_mongo_db").C("test").Insert(loc)    
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(loc)
}


func getFormattedURL(url string) string {
    
    if url != ""   {
        url = strings.Replace(url," ","%20",-1)
    }
    return url
}


func (lc LocationController) GetLocation(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    uid := params.Get(":location_id")
	  oid := bson.ObjectIdHex(uid)
	  loc := models.Location{}
    err := lc.session.DB("test_mongo_db").C("test").FindId(oid).One(&loc)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    
  	}
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    json.NewEncoder(w).Encode(loc)
}

func (lc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request) {
    
    params := r.URL.Query()
    uid := params.Get(":location_id")
    oid := bson.ObjectIdHex(uid)
    loc := models.Location{}
    oldloc := models.Location{}

    decoder := json.NewDecoder(r.Body)
    err := decoder.Decode(&loc)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Input JSON: %v", err)
    }

    err = lc.session.DB("test_mongo_db").C("test").FindId(oid).One(&oldloc)
    if err != nil {
      fmt.Printf("got an error finding a doc %v\n")
    
    }

    var address string
    address = loc.Address+loc.City+loc.State+loc.Zip
    lat, lng := GetCoordinates(address)
    loc.Coordinate.Lat=lat
    loc.Coordinate.Lng=lng

    loc.ID = oid

    if loc.Name == "" {
      loc.Name = oldloc.Name
    }


    err = lc.session.DB("test_mongo_db").C("test").UpdateId(oid,loc)
    
    if err != nil {
      fmt.Printf("got an error updating the doc %v\n")
    
    }

    err = lc.session.DB("test_mongo_db").C("test").FindId(oid).One(&loc)
    if err != nil {
      fmt.Printf("got an error finding a doc %v\n")
    
    }
    
    w.Header().Set("Content-Type", "application/json; charset=UTF-8")
    json.NewEncoder(w).Encode(loc)   

}

func (lc LocationController) DeleteLocation(w http.ResponseWriter, r *http.Request) {
    params := r.URL.Query()
    uid := params.Get(":location_id")
	  oid := bson.ObjectIdHex(uid)
	
	  err := lc.session.DB("test_mongo_db").C("test").RemoveId(oid)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    
  	}  	
  	w.WriteHeader(200)
}
