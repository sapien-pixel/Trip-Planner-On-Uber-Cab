package models

type Estimate struct {
	Prices []struct {
		CurrencyCode         string  `json:"currency_code"`
		DisplayName          string  `json:"display_name"`
		Distance             float64 `json:"distance"`
		Duration             int     `json:"duration"`
		Estimate             string  `json:"estimate"`
		HighEstimate         int     `json:"high_estimate"`
		LocalizedDisplayName string  `json:"localized_display_name"`
		LowEstimate          int     `json:"low_estimate"`
		Minimum              int     `json:"minimum"`
		ProductID            string  `json:"product_id"`
		SurgeMultiplier      int     `json:"surge_multiplier"`
	} `json:"prices"`
}