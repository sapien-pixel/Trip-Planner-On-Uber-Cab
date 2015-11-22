package models

import (	
	"gopkg.in/mgo.v2/bson"
)

type(
	
	Trip struct {

		ID                     bson.ObjectId   	`json:"id" bson:"_id,omitempty"`
		BestRouteLocationIds   []string 		`json:"best_route_location_ids,omitempty"`
		LocationIds            []string 		`json:"location_ids"`
		StartingFromLocationID string   		`json:"starting_from_location_id"`
		Status                 string   		`json:"status,omitempty"`
		TotalDistance          float64  		`json:"total_distance,omitempty"`
		TotalUberCosts         int      		`json:"total_uber_costs,omitempty"`
		TotalUberDuration      int      		`json:"total_uber_duration,omitempty"`
	}

	TripResponse struct {

		ID                     bson.ObjectId   	`json:"id" bson:"_id,omitempty"`
		BestRouteLocationIds   []string 		`json:"best_route_location_ids,omitempty"`
		StartingFromLocationID string   		`json:"starting_from_location_id"`
		Status                 string   		`json:"status,omitempty"`
		TotalDistance          float64  		`json:"total_distance,omitempty"`
		TotalUberCosts         int      		`json:"total_uber_costs,omitempty"`
		TotalUberDuration      int      		`json:"total_uber_duration,omitempty"`
	}

	TripRequest struct	{

		ID                     		bson.ObjectId   `json:"id" bson:"_id,omitempty"`
		Status                 		string   		`json:"status,omitempty"`
		StartingFromLocationID 		string   		`json:"starting_from_location_id"`
		NextDestinationLocationID 	string 			`json:"next_destination_location_id"`
		BestRouteLocationIds   		[]string 		`json:"best_route_location_ids,omitempty"`
		TotalUberCosts         		int      		`json:"total_uber_costs,omitempty"`
		TotalUberDuration      		int      		`json:"total_uber_duration,omitempty"`
		TotalDistance          		float64  		`json:"total_distance,omitempty"`
		UberWaitTimeEta 			int 			`json:"uber_wait_time_eta"`
		
	}

	RideRequest struct {
		ProductID      string `json:"product_id"`
		StartLatitude  string `json:"start_latitude"`
		StartLongitude string `json:"start_longitude"`
		EndLatitude    string `json:"end_latitude"`
		EndLongitude   string `json:"end_longitude"`
	}

	RideResponse struct {
		Driver          interface{} `json:"driver"`
		Eta             int         `json:"eta"`
		Location        interface{} `json:"location"`
		RequestID       string      `json:"request_id"`
		Status          string      `json:"status"`
		SurgeMultiplier int         `json:"surge_multiplier"`
		Vehicle         interface{} `json:"vehicle"`
	}

	UberProducts struct {
	Products []struct {
		Capacity     int    `json:"capacity"`
		Description  string `json:"description"`
		DisplayName  string `json:"display_name"`
		Image        string `json:"image"`
		PriceDetails struct {
			Base            float64 `json:"base"`
			CancellationFee int     `json:"cancellation_fee"`
			CostPerDistance float64 `json:"cost_per_distance"`
			CostPerMinute   float64 `json:"cost_per_minute"`
			CurrencyCode    string  `json:"currency_code"`
			DistanceUnit    string  `json:"distance_unit"`
			Minimum         float64 `json:"minimum"`
			ServiceFees     []struct {
				Fee  float64 `json:"fee"`
				Name string  `json:"name"`
			} `json:"service_fees"`
		} `json:"price_details"`
		ProductID string `json:"product_id"`
	} `json:"products"`
}
)