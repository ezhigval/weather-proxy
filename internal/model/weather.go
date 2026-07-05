package model

type Weather struct {
	City        string  `json:"city"`
	Temperature float64 `json:"temperature_c"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed_ms"`
	Cached      bool    `json:"cached"`
}

type BatchResult struct {
	City    string   `json:"city"`
	Weather *Weather `json:"weather,omitempty"`
	Error   string   `json:"error,omitempty"`
}

type BatchResponse struct {
	Results []BatchResult `json:"results"`
}
