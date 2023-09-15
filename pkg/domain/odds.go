package domain

import "time"

type OddsResponse struct {
	Id           string                  `json:"id"`
	SportKey     string                  `json:"sport_key"`
	SportTitle   string                  `json:"sport_title"`
	CommenceTime time.Time               `json:"commence_time"`
	HomeTeam     string                  `json:"home_team"`
	AwayTeam     string                  `json:"away_team"`
	Bookmakers   []BookmakerOddsResponse `json:"bookmakers"`
}

type BookmakerOddsResponse struct {
	Key        string    `json:"key"`
	Title      string    `json:"title"`
	LastUpdate time.Time `json:"last_update"`
	Markets    []struct {
		Key        string    `json:"key"`
		LastUpdate time.Time `json:"last_update"`
		Outcomes   []struct {
			Name  string  `json:"name"`
			Price float64 `json:"price"`
		} `json:"outcomes"`
	} `json:"markets"`
}
