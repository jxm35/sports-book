package scrape_stats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type matchResponse struct {
	Events []matchEvent `json:"events"`
}
type matchEvent struct {
	Id             int           `json:"id"`
	AwayScore      scoreResponse `json:"awayScore"`
	AwayTeam       teamResponse  `json:"awayteam"`
	HomeScore      scoreResponse `json:"homeScore"`
	HomeTeam       teamResponse  `json:"hometeam"`
	StartTimestamp int64         `json:"startTimestamp"`
}

type scoreResponse struct {
	Current    int `json:"current"`
	Display    int `json:"display"`
	NormalTime int `json:"normaltime"`
	Period1    int `json:"period1"`
	Period2    int `json:"period2"`
}

type teamResponse struct {
	Name      string `json:"name"`
	Slug      string `json:"slug"`
	ShortName string `json:"shortName"`
	NameCode  string `json:"nameCode"`
	ID        int    `json:"id"`
}

func RequestMatches(tournamentId string, gameweek int) (matchResponse, error) {
	var matches matchResponse

	req := fmt.Sprintf("https://api.sofascore.com/api/v1/unique-tournament/17/season/%s/events/round/%d", tournamentId, gameweek)

	resp, err := requestGet(req)
	if err != nil {
		return matches, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return matches, ErrRequestForbidden
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return matches, err
	}
	err = json.Unmarshal(responseData, &matches)
	if err != nil {
		print(string(responseData))
		return matches, err
	}
	return matches, err
}
