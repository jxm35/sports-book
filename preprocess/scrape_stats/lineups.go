package scrape_stats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type lineupResponse struct {
	Confirmed bool               `json:"confirmed"`
	Home      teamLineupResponse `json:"home"`
	Away      teamLineupResponse `json:"away"`
}

type teamLineupResponse struct {
	Players []playerResponse `json:"players"`
}

type playerResponse struct {
	Player     playerInnerResponse          `json:"player"`
	Statistics appearanceStatisticsResponse `json:"statistics"`
}

type playerInnerResponse struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Position string `json:"position"`
}

type appearanceStatisticsResponse struct {
	MinutesPlayed   float32 `json:"minutesPlayed"`
	ExpectedGoals   float32 `json:"expectedGoals"`
	ExpectedAssists float32 `json:"expectedAssists"`
	Goals           int     `json:"goals"`
	Asists          int     `json:"goalAssist"`
}

func RequestLineups(matchId int) (lineupResponse, error) {
	var lineups lineupResponse

	req := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%d/lineups", matchId)

	resp, err := requestGet(req)
	if err != nil {
		return lineups, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return lineups, ErrRequestForbidden
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return lineups, err
	}
	err = json.Unmarshal(responseData, &lineups)
	if err != nil {
		print(string(responseData))
		return lineups, err
	}
	return lineups, err
}
