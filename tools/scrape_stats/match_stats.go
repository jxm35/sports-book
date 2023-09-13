package scrape_stats

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type matchStatsResponse struct {
	Statistics []struct {
		Period string `json:"period"`
		Groups []struct {
			GroupName       string `json:"groupName"`
			StatisticsItems []struct {
				Name        string `json:"name"`
				Home        string `json:"home"`
				Away        string `json:"away"`
				CompareCode int    `json:"compareCode"`
			} `json:"statisticsItems"`
		} `json:"groups"`
	} `json:"statistics"`
}

func RequestMatchStats(matchId int) (matchStatsResponse, error) {
	var stats matchStatsResponse

	req := fmt.Sprintf("https://api.sofascore.com/api/v1/event/%d/statistics", matchId)

	resp, err := requestGet(req)
	if err != nil {
		return stats, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusForbidden {
		return stats, ErrRequestForbidden
	}
	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return stats, err
	}
	err = json.Unmarshal(responseData, &stats)
	if err != nil {
		print(string(responseData))
		return stats, err
	}
	return stats, err
}
