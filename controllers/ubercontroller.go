package controllers

import(
	"fmt"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"models"
	"encoding/binary"
    "bytes"
)

type (

	UberController struct	{

			session *mgo.Session
	}

)

func NewUberController(s *mgo.Session) *UberController {
	return &UberController{s}
}

func (uc UberController) GetEstimates(locationid string) {
	oid := bson.ObjectIdHex(locationid)
	loc := models.Location{}
    err := uc.session.DB("test_mongo_db").C("test").FindId(oid).One(&loc)
  	if err != nil {
    	fmt.Printf("got an error finding a doc %v\n")
    
  	}
  	fmt.Print(loc)
}

func (uc UberController) AddTrip(w http.ResponseWriter, r *http.Request) {

	trip := models.Trip{}
	tripResp := models.TripResponse{}
	location := models.Location{}

	loc := make(map[string] models.Location)

	decoder := json.NewDecoder(r.Body) 
    err := decoder.Decode(&trip)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Input JSON: %v", err)
    }
    tripResp.StartingFromLocationID = trip.StartingFromLocationID
    tripResp.ID = bson.NewObjectId()
    tripResp.Status = "planning"

    //Initial Location
    oid := bson.ObjectIdHex(trip.StartingFromLocationID)
	err = uc.session.DB("test_mongo_db").C("test").FindId(oid).One(&location)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	}

	loc[trip.StartingFromLocationID] = location
	//fmt.Print(json.NewEncoder(w).Encode(loc[trip.StartingFromLocationID]))

	for _, each := range trip.LocationIds {
		
		oid = bson.ObjectIdHex(each)
		err = uc.session.DB("test_mongo_db").C("test").FindId(oid).One(&location)
		if err != nil {
			fmt.Printf("got an error finding a doc %v\n")
		}
		loc[each] = location
		//fmt.Print(json.NewEncoder(w).Encode(loc[each]))		
	}

	/*for _, each := range loc {
		fmt.Print(json.NewEncoder(w).Encode(each))
	}*/

	var price int
	var duration int
	var distance float64

	startId := trip.StartingFromLocationID
	startLat := loc[startId].Coordinate.Lat
	startLng := loc[startId].Coordinate.Lng
	nextId := startId
	originLat := loc[startId].Coordinate.Lat
	originLng := loc[startId].Coordinate.Lng

	minPrice := 99999
	minDuration := 0
	minDistance := 0.0
	
	totalCost := 0
	totalDuration := 0
	totalDistance := 0.0
	pos := 0
	for len(loc) > 1 {
		for index, each := range loc 	{
			if index != startId	{
				
				price,duration,distance = GetEstimates(startLat,startLng,each.Coordinate.Lat,each.Coordinate.Lng)
				
				if price < minPrice	{
					minPrice = price
					minDuration = duration
					minDistance = distance
					nextId = index
				}
			}
			
		}
		trip.LocationIds[pos] = nextId
		//fmt.Println("The min price is: %d with startid: %s and nextId: %s", minPrice, startId, nextId)
		totalCost+=minPrice
		totalDuration+=minDuration
		totalDistance+=minDistance
		delete(loc,startId)
		startId = nextId
		startLat = loc[startId].Coordinate.Lat
		startLng = loc[startId].Coordinate.Lng
		minPrice = 99999
		minDuration = 0
		minDistance = 0.00
		pos++
	}

   	price,duration,distance = GetEstimates(startLat,startLng,originLat,originLng)
   	totalCost+=price
	totalDuration+=duration
	totalDistance+=distance

	tripResp.BestRouteLocationIds = trip.LocationIds
	tripResp.TotalUberCosts = totalCost
	tripResp.TotalUberDuration = totalDuration
	tripResp.TotalDistance = totalDistance


	err = uc.session.DB("test_mongo_db").C("trip_details").Insert(tripResp) 
	if err != nil {
		fmt.Printf("Can't insert document: %v\n", err)
	}

	err = uc.session.DB("test_mongo_db").C("trip_details").FindId(tripResp.ID).One(&tripResp)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(tripResp)

}

func GetEstimates(startLat float64, startLng float64, endLat float64, endLng float64)	(int,int,float64)	{

	estimate := models.Estimate{}
	url := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%f&start_longitude=%f&end_latitude=%f&end_longitude=%f",startLat,startLng,endLat,endLng)
	client := http.Client{Timeout: timeout}
   	req,_ := http.NewRequest("GET",url,nil)
   	req.Header.Set("Authorization", "Token iD0QDvQe5pG6cMVb4Q23vwlPldl-CvtLkGbfFj65") 
    res, err := client.Do(req)
    if err != nil {
        fmt.Errorf("Cannot read Google API: %v", err)
    }
    defer res.Body.Close()
    decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&estimate)
    return estimate.Prices[0].HighEstimate, estimate.Prices[0].Duration, estimate.Prices[0].Distance
}

func (uc UberController) GetTrip(w http.ResponseWriter, r *http.Request) {

	params := r.URL.Query()
    uid := params.Get(":trip_id")
    oid := bson.ObjectIdHex(uid)
	tripResp := models.TripResponse{}
	err := uc.session.DB("test_mongo_db").C("trip_details").FindId(oid).One(&tripResp)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	}
	
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(tripResp)

}

func (uc UberController) RequestTrip(w http.ResponseWriter, r *http.Request)	{

	params := r.URL.Query()
    uid := params.Get(":trip_id")
    oid := bson.ObjectIdHex(uid)
	tripResp := models.TripResponse{}
	err := uc.session.DB("test_mongo_db").C("trip_details").FindId(oid).One(&tripResp)
	if err != nil {
		fmt.Printf("got an error finding a doc %v\n")
	}
		tripReq := models.TripRequest{}
		tripReq.ID =  tripResp.ID             		
		tripReq.Status = tripResp.Status                		
		tripReq.StartingFromLocationID = tripResp.StartingFromLocationID	
		tripReq.BestRouteLocationIds = tripResp.BestRouteLocationIds  		
		tripReq.TotalUberCosts  = tripResp.TotalUberCosts
		tripReq.TotalUberDuration = tripResp.TotalUberDuration     		
		tripReq.TotalDistance  = tripResp.TotalDistance			

	if tripReq.NextDestinationLocationID == tripReq.StartingFromLocationID	{
		tripReq.Status = "trip over"
		tripReq.NextDestinationLocationID = ""
		tripReq.StartingFromLocationID = ""
		tripReq.UberWaitTimeEta = 0
	}	else	{

			if tripReq.Status == "planning"	{
				tripReq.Status = "requesting"
				tripReq.NextDestinationLocationID = tripReq.BestRouteLocationIds[0]

			}	else if tripResp.Status == "requesting"	{
					if len(tripReq.BestRouteLocationIds)>1	{
						value := tripReq.BestRouteLocationIds[1:len(tripReq.BestRouteLocationIds)]
						tripReq.BestRouteLocationIds = value
						tripReq.NextDestinationLocationID = tripReq.BestRouteLocationIds[0]
					}	else	{
						tripReq.BestRouteLocationIds = nil
						tripReq.NextDestinationLocationID,tripReq.StartingFromLocationID = tripReq.StartingFromLocationID,tripReq.NextDestinationLocationID
					}
				}

	  			oid := bson.ObjectIdHex(tripReq.StartingFromLocationID)
	 			loc := models.Location{}
   				err := uc.session.DB("test_mongo_db").C("test").FindId(oid).One(&loc)
  				if err != nil {
    				fmt.Printf("got an error finding a doc %v\n")
    
  				}
				startLat:=loc.Coordinate.Lat
				startLng:=loc.Coordinate.Lng

				oid = bson.ObjectIdHex(tripReq.StartingFromLocationID)
	 			loc = models.Location{}
   				err = uc.session.DB("test_mongo_db").C("test").FindId(oid).One(&loc)
  				if err != nil {
    				fmt.Printf("got an error finding a doc %v\n")
    
  				}
				nextLat:=loc.Coordinate.Lat
				nextLng:=loc.Coordinate.Lng

				tripReq.UberWaitTimeEta = GetEstimatedTime(startLat,startLng,nextLat,nextLng)

				if len(tripReq.BestRouteLocationIds) == 0	{
					tripReq.StartingFromLocationID = tripReq.NextDestinationLocationID
				}
		}

		tripResp.ID =  tripReq.ID             		
		tripResp.Status = tripReq.Status                		
		tripResp.StartingFromLocationID = tripReq.StartingFromLocationID	
		tripResp.BestRouteLocationIds = tripReq.BestRouteLocationIds  		
		tripResp.TotalUberCosts  = tripReq.TotalUberCosts
		tripResp.TotalUberDuration = tripReq.TotalUberDuration     		
		tripResp.TotalDistance  = tripReq.TotalDistance	

		err = uc.session.DB("test_mongo_db").C("trip_details").UpdateId(oid,tripResp)
  		if err != nil {
    		fmt.Printf("got an error updating a doc %v\n")
    	} 
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(tripReq)
}

func GetEstimatedTime(startLat float64, startLng float64,nextLat float64,nextLng float64) (int)	{

	UBER_ACCESS_TOKEN := "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicHJvZmlsZSIsImhpc3RvcnkiXSwic3ViIjoiOWRlMWU4MTktNjBkMi00MmNjLTlkYTAtNjFkZDNjODkwYjE3IiwiaXNzIjoidWJlci11czEiLCJqdGkiOiJlMzJmZDY3ZC0yNTIwLTQ0ZTMtYWViYy1kNjZlNzM4MjhjZGEiLCJleHAiOjE0NTA3MzU3MDYsImlhdCI6MTQ0ODE0MzcwNSwidWFjdCI6InpXWjFRNlBxbmdYQW5ZYjdDcXNYZExsZHZyenMwZiIsIm5iZiI6MTQ0ODE0MzYxNSwiYXVkIjoibElGN2lUUXJqUFUxdmEteGZCREx5RU4zQlhnR3V4Rk8ifQ.TJcZb3tExZMhcV6xGoiiXJVamQtzOSrDAabNp8ET11mfGCAHVASUIweSqIfVFRw2HtOJz9m2rgM_5JKue1quE8xA_IPN3zBOj3g4cDdiz7AxdV2G-pcb_C2I4tUTk_LdabkeKOEyYxrwQyMrMT5otdCGfRRpwuxqugWI6CUzQdYH6edlZ-OzDvbgRgC9EqcXPo-ATImUggTAllyHrByxgS8bFxN4j3yundXhfce2BhhaErAsdGlAFRkhoc16SaH9XOWmoLpwdUy9jF3lZhm31umYFSIXUhlWgNcqD5KtrGI1OybjB7GJX9gP69101fcuQ0aBqYBd4TVd9zji1n4jXw"
	rideRequest := models.RideRequest{}
	rideResponse := models.RideResponse{}
	rideRequest.ProductID = GetUberProductID(startLat,startLng)
	rideRequest.StartLatitude = fmt.Sprintf("%.6f",startLat)
	rideRequest.StartLongitude = fmt.Sprintf("%.6f",startLng)
	rideRequest.EndLatitude = fmt.Sprintf("%.6f",nextLat)
	rideRequest.EndLongitude = fmt.Sprintf("%.6f",nextLng)

	url := fmt.Sprintf("%s/requests","https://sandbox-api.uber.com/v1")
	client := http.Client{Timeout: timeout}
	b, err := json.Marshal(rideRequest)
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, &b)
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization",UBER_ACCESS_TOKEN)
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
        fmt.Errorf("Error in UBER ride request API: %v", err)
    }

	decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&rideResponse)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the UBER_RideRequest JSON: %v", err)
    }

    return rideResponse.Eta
}

func GetUberProductID(latitude float64, longitude float64) string {

	UBER_SERVER_TOKEN := "Token iD0QDvQe5pG6cMVb4Q23vwlPldl-CvtLkGbfFj65"
	uberDetails := models.UberProducts{}
    url := fmt.Sprintf("%s/products?latitude=%f&longitude=%f","https://sandbox-api.uber.com/v1",latitude,longitude)
	client := http.Client{Timeout: timeout}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization",UBER_SERVER_TOKEN)
	res, err := client.Do(req)
	defer res.Body.Close()
    if err != nil {
        fmt.Errorf("Cannot read UBER API: %v", err)
    }
	decoder := json.NewDecoder(res.Body)
    err = decoder.Decode(&uberDetails)
    if(err!=nil)    {
        fmt.Errorf("Error in decoding the Google JSON: %v", err)
    }
    return uberDetails.Products[0].ProductID
}