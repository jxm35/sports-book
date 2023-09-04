package scrape_stats

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func ScrapeSeason(id, startWeek, maxWeek int) {
	err := CreateFile()
	var matchesToSave []Match
	var teamsToSave []Team
	var playersToSave []Player
	var appearancesToSave []Appearance

	_, err = excelize.OpenFile("data.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	tournamentId := fmt.Sprintf("%d", id)
	for gameweek := startWeek; gameweek <= maxWeek; gameweek++ {
		matches, err := RequestMatches(tournamentId, gameweek)
		if err != nil {
			fmt.Println(err)
			SaveData(matchesToSave, teamsToSave, playersToSave, appearancesToSave)
			return
		}
		for _, match := range matches.Events {

			matchStats, err := RequestMatchStats(match.Id)
			if err != nil {
				fmt.Println(err)
				SaveData(matchesToSave, teamsToSave, playersToSave, appearancesToSave)
				return
			}
			var homeXg string
			var awayxG string
			var homeXa string
			var awayxa string
			for _, stat := range matchStats.Statistics {
				if stat.Period == "ALL" {
					for _, group := range stat.Groups {
						if group.GroupName == "Expected" {
							for _, item := range group.StatisticsItems {
								if item.Name == "Expected goals" {
									homeXg = item.Home
									awayxG = item.Away
								}
								if item.Name == "Expected assists" {
									homeXa = item.Home
									awayxa = item.Away
								}
							}
						}
					}
				}
			}

			saveMatch := Match{
				Id:        match.Id,
				Date:      match.StartTimestamp,
				HomeTeam:  match.HomeTeam.ID,
				AwayTeam:  match.AwayTeam.ID,
				HomexG:    homeXg,
				AwayxG:    awayxG,
				HomeGoals: match.HomeScore.NormalTime,
				AwayGoals: match.AwayScore.NormalTime,
				HomexA:    homeXa,
				AwayxA:    awayxa,
			}
			matchesToSave = append(matchesToSave, saveMatch)

			homeTeam := Team{
				Id:   match.HomeTeam.ID,
				Name: match.HomeTeam.Name,
			}
			awayTeam := Team{
				Id:   match.HomeTeam.ID,
				Name: match.HomeTeam.Name,
			}
			teamsToSave = append(teamsToSave, homeTeam)
			teamsToSave = append(teamsToSave, awayTeam)

			lineups, err := RequestLineups(match.Id)
			if err != nil {
				fmt.Println(err)
				SaveData(matchesToSave, teamsToSave, playersToSave, appearancesToSave)
				return
			}
			for _, app := range lineups.Home.Players {
				playerToSave := Player{
					Id:       app.Player.Id,
					Name:     app.Player.Name,
					Position: app.Player.Position,
				}
				playersToSave = append(playersToSave, playerToSave)

				appToSave := Appearance{
					MatchId: match.Id,
					Team:    match.HomeTeam.ID,
					XG:      app.Statistics.ExpectedGoals,
					XA:      app.Statistics.ExpectedAssists,
					Goals:   app.Statistics.Goals,
					Assists: app.Statistics.Asists,
					Minutes: app.Statistics.MinutesPlayed,
				}
				appearancesToSave = append(appearancesToSave, appToSave)
			}

			for _, app := range lineups.Away.Players {
				playerToSave := Player{
					Id:       app.Player.Id,
					Name:     app.Player.Name,
					Position: app.Player.Position,
				}
				playersToSave = append(playersToSave, playerToSave)

				appToSave := Appearance{
					Player:  app.Player.Id,
					MatchId: match.Id,
					Team:    match.AwayTeam.ID,
					XG:      app.Statistics.ExpectedGoals,
					XA:      app.Statistics.ExpectedAssists,
					Goals:   app.Statistics.Goals,
					Assists: app.Statistics.Asists,
					Minutes: app.Statistics.MinutesPlayed,
				}
				appearancesToSave = append(appearancesToSave, appToSave)
			}
		}
	}
	SaveData(matchesToSave, teamsToSave, playersToSave, appearancesToSave)
}

func RectifyAppearances(file string) {
	var matchIds []string
	var awayIds []string
	var homeIds []string
	var appearancesToSave []Appearance
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := f.GetRows("Matches")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i, row := range rows {
		if i != 0 {
			matchIds = append(matchIds, row[0])
			homeIds = append(homeIds, row[2])
			awayIds = append(awayIds, row[3])
		}
	}

	for i, id := range matchIds {

		matchInt, err := strconv.Atoi(id)
		if err != nil {
			fmt.Println(err)
			continue
		}
		lineups, err := RequestLineups(matchInt)
		if err != nil {
			fmt.Println(err)
			return
		}
		homeInt, err := strconv.Atoi(homeIds[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		awayInt, err := strconv.Atoi(awayIds[i])
		if err != nil {
			fmt.Println(err)
			continue
		}
		for _, app := range lineups.Home.Players {

			appToSave := Appearance{
				MatchId: matchInt,
				Player:  app.Player.Id,
				Team:    homeInt,
				XG:      app.Statistics.ExpectedGoals,
				XA:      app.Statistics.ExpectedAssists,
				Goals:   app.Statistics.Goals,
				Assists: app.Statistics.Asists,
				Minutes: app.Statistics.MinutesPlayed,
			}
			appearancesToSave = append(appearancesToSave, appToSave)
		}
		for _, app := range lineups.Away.Players {

			appToSave := Appearance{
				Player:  app.Player.Id,
				MatchId: matchInt,
				Team:    awayInt,
				XG:      app.Statistics.ExpectedGoals,
				XA:      app.Statistics.ExpectedAssists,
				Goals:   app.Statistics.Goals,
				Assists: app.Statistics.Asists,
				Minutes: app.Statistics.MinutesPlayed,
			}
			appearancesToSave = append(appearancesToSave, appToSave)
		}
	}
	SaveNewAppearances(appearancesToSave, file)
}
