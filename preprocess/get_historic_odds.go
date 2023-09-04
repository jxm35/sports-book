package preprocess

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"sports-book.com/model"
	"sports-book.com/query"
)

func SaveOdds(year int32, league string) {
	gormDb, err := gorm.Open(mysql.Open("root:password@tcp(127.0.0.1:3306)/sports-book?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		fmt.Println(err)
		return
	}
	query.SetDefault(gormDb)

	filename := fmt.Sprintf("./preprocess/oddsFiles/%s/%dodds.csv", league, year)
	allOdds := make([]MatchOdds1x2, 0)
	b365odds, _ := Get1x2Odds(filename, "Bet365", "B365H", "B365D", "B365A")
	ladbrokesOdds, _ := Get1x2Odds(filename, "ladbrokes", "LBH", "LBD", "LBA")
	pinnacleOdds, _ := Get1x2Odds(filename, "pinnacle", "PSH", "PSD", "PSA")
	willHillOdds, _ := Get1x2Odds(filename, "william hill", "WHH", "WHD", "WHA")
	allOdds = append(allOdds, b365odds...)
	allOdds = append(allOdds, ladbrokesOdds...)
	allOdds = append(allOdds, pinnacleOdds...)
	allOdds = append(allOdds, willHillOdds...)
	if len(allOdds) == 0 {
		panic("no odds found")
	}

	Save1x2odds(allOdds, year, league)
}

type MatchOdds1x2 struct {
	Bookie   string
	HomeTeam string
	AwayTeam string
	HomeWin  float64
	Draw     float64
	AwayWin  float64
}

func Get1x2Odds(fileName string, bookies, hw, d, aw string) ([]MatchOdds1x2, error) {
	csvGrid := readCsvFile(fileName)

	allOdds := make([]MatchOdds1x2, 0)

	homeTeamIdx, err := getColumnFromTitle(csvGrid, "HomeTeam")
	if err != nil {
		return nil, err
	}
	awayTeamIdx, err := getColumnFromTitle(csvGrid, "AwayTeam")
	if err != nil {
		return nil, err
	}
	homeWinIdx, err := getColumnFromTitle(csvGrid, hw)
	if err != nil {
		return nil, err
	}
	drawIdx, err := getColumnFromTitle(csvGrid, d)
	if err != nil {
		return nil, err
	}
	awayWinIdx, err := getColumnFromTitle(csvGrid, aw)
	if err != nil {
		return nil, err
	}

	for idx, row := range csvGrid {
		if idx != 0 {
			homeWinString := row[homeWinIdx]
			drawString := row[drawIdx]
			awayWinString := row[awayWinIdx]
			if homeWinString == "" && drawString == "" && awayWinString == "" {
				log.Println("empty row found", idx)
				continue
			}
			homeWin, err := strconv.ParseFloat(homeWinString, 64)
			if err != nil {
				panic(err)
			}
			draw, err := strconv.ParseFloat(drawString, 64)
			if err != nil {
				panic(err)
			}
			awayWin, err := strconv.ParseFloat(awayWinString, 64)
			if err != nil {
				panic(err)
			}
			odds := MatchOdds1x2{
				Bookie:   bookies,
				HomeTeam: row[homeTeamIdx],
				AwayTeam: row[awayTeamIdx],
				HomeWin:  homeWin,
				Draw:     draw,
				AwayWin:  awayWin,
			}
			allOdds = append(allOdds, odds)
		}
	}
	return allOdds, nil
}

func Save1x2odds(odds []MatchOdds1x2, year int32, league string) {
	teamMap := getTeamMap()
	oddsToSave := make([]*model.Odds1x2, 0)
	for _, odd := range odds {
		match := getMatchId(teamMap[odd.HomeTeam], teamMap[odd.AwayTeam], year)
		if match <= 0 {
			panic("could not find match")
		}
		toSave := model.Odds1x2{
			Bookmaker: odd.Bookie,
			Match:     match,
			HomeWin:   odd.HomeWin,
			Draw:      odd.Draw,
			AwayWin:   odd.AwayWin,
		}
		oddsToSave = append(oddsToSave, &toSave)
	}
	log.Println(fmt.Sprintf("saved %d odds", len(oddsToSave)))
	err := query.Odds1x2.WithContext(context.Background()).CreateInBatches(oddsToSave, 200)
	if err != nil {
		fmt.Println(err)
	}
}

func getMatchId(homeTeam, awayTeam, year int32) int32 {
	m := query.Match
	c := query.Competition
	var match model.Match
	err := m.WithContext(context.Background()).
		Select(m.ALL).
		LeftJoin(c, m.Competition.EqCol(c.ID)).
		Where(m.HomeTeam.Eq(homeTeam), m.AwayTeam.Eq(awayTeam), c.Year.Eq(year)).
		Scan(&match)
	if err != nil {
		return -1
	}
	return match.ID
}

func getTeamMap() map[string]int32 {
	resp := make(map[string]int32)
	var teams []model.Team
	t := query.Team
	err := t.WithContext(context.Background()).
		Select(t.ALL).Scan(&teams)
	if err != nil {
		return nil
	}
	for _, team := range teams {
		switch team.Name {
		case "Manchester United":
			team.Name = "Man United"
		case "Newcastle United":
			team.Name = "Newcastle"
		case "Wolverhampton Wanderers":
			team.Name = "Wolves"
		case "Manchester City":
			team.Name = "Man City"
		case "West Bromwich Albion":
			team.Name = "West Brom"
		case "Nottingham Forest":
			team.Name = "Nott'm Forest"
		case "Queens Park Rangers":
			team.Name = "QPR"

			// la liga
		case "Espanyol":
			team.Name = "Espanol"
		case "Deportivo La Coruna":
			team.Name = "La Coruna"
		case "Athletic Club":
			team.Name = "Ath Bilbao"
		case "Celta Vigo":
			team.Name = "Celta"
		case "Real Sociedad":
			team.Name = "Sociedad"
		case "Rayo Vallecano":
			team.Name = "Vallecano"
		case "Atletico Madrid":
			team.Name = "Ath Madrid"
		case "Real Betis":
			team.Name = "Betis"
		case "Sporting Gijon":
			team.Name = "Sp Gijon"
		case "Real Valladolid":
			team.Name = "Valladolid"
		case "SD Huesca":
			team.Name = "Huesca"

		// bundesliga
		case "Borussia Dortmund":
			team.Name = "Dortmund"
		case "Bayer Leverkusen":
			team.Name = "Leverkusen"
		case "Eintracht Frankfurt":
			team.Name = "Ein Frankfurt"
		case "FC Cologne":
			team.Name = "FC Koln"
		case "Hamburger SV":
			team.Name = "Hamburg"
		case "Hannover 96":
			team.Name = "Hannover"
		case "Hertha Berlin":
			team.Name = "Hertha"
		case "Borussia M.Gladbach":
			team.Name = "M'gladbach"
		case "VfB Stuttgart":
			team.Name = "Stuttgart"
		case "Mainz 05":
			team.Name = "Mainz"
		case "RasenBallsport Leipzig":
			team.Name = "RB Leipzig"
		case "Fortuna Duesseldorf":
			team.Name = "Fortuna Dusseldorf"
		case "Nuernberg":
			team.Name = "Nurnberg"
		case "Arminia Bielefeld":
			team.Name = "Bielefeld"
		case "Greuther Fuerth":
			team.Name = "Greuther Furth"

		// serie a
		case "AC Milan":
			team.Name = "Milan"
		case "SPAL 2013":
			team.Name = "Spal"
		case "Parma Calcio 1913":
			team.Name = "Parma"

		}
		resp[team.Name] = team.ID
	}
	return resp
}

func readCsvFile(filePath string) [][]string {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	return records
}

var ErrColumnNotFound = errors.New("could not find col")

func getColumnFromTitle(csvGrid [][]string, title string) (int, error) {
	row := csvGrid[0]
	for idx, header := range row {
		if header == title {
			return idx, nil
		}
		if idx > 200 {
			return -1, ErrColumnNotFound
		}
	}
	return -1, ErrColumnNotFound
}
