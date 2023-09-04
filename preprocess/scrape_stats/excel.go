package scrape_stats

import (
	"fmt"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func CreateFile() error {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()
	// Create a new sheet.
	_, err := f.NewSheet("Matches")
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.SetCellValue("Matches", "A1", "Id")
	f.SetCellValue("Matches", "B1", "Date")
	f.SetCellValue("Matches", "C1", "HomeTeam")
	f.SetCellValue("Matches", "D1", "AwayTeam")
	f.SetCellValue("Matches", "E1", "HomexG")
	f.SetCellValue("Matches", "F1", "AwayxG")
	f.SetCellValue("Matches", "G1", "HomeGoals")
	f.SetCellValue("Matches", "H1", "AwayGoals")
	f.SetCellValue("Matches", "I1", "HomexA")
	f.SetCellValue("Matches", "J1", "AwayxA")

	_, err = f.NewSheet("Players")
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.SetCellValue("Players", "A1", "Id")
	f.SetCellValue("Players", "B1", "Name")
	f.SetCellValue("Players", "C1", "Position")

	_, err = f.NewSheet("Teams")
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.SetCellValue("Teams", "A1", "Id")
	f.SetCellValue("Teams", "B1", "Name")

	_, err = f.NewSheet("Appearances")
	if err != nil {
		fmt.Println(err)
		return err
	}
	f.SetCellValue("Appearances", "A1", "MatchId")
	f.SetCellValue("Appearances", "B1", "Team")
	f.SetCellValue("Appearances", "C1", "XG")
	f.SetCellValue("Appearances", "D1", "XA")
	f.SetCellValue("Appearances", "E1", "Goals")
	f.SetCellValue("Appearances", "F1", "Assists")
	f.SetCellValue("Appearances", "G1", "Minutes")

	if err := f.SaveAs("data.xlsx"); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func SaveData(
	matchesToSave []Match,
	teamsToSave []Team,
	playersToSave []Player,
	appearancesToSave []Appearance,
) {
	f, err := excelize.OpenFile("data.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range matchesToSave {
		f.SetCellValue("Matches", "A"+strconv.Itoa(i+2), matchesToSave[i].Id)
		f.SetCellValue("Matches", "B"+strconv.Itoa(i+2), matchesToSave[i].Date)
		f.SetCellValue("Matches", "C"+strconv.Itoa(i+2), matchesToSave[i].HomeTeam)
		f.SetCellValue("Matches", "D"+strconv.Itoa(i+2), matchesToSave[i].AwayTeam)
		f.SetCellValue("Matches", "E"+strconv.Itoa(i+2), matchesToSave[i].HomexG)
		f.SetCellValue("Matches", "F"+strconv.Itoa(i+2), matchesToSave[i].AwayxG)
		f.SetCellValue("Matches", "G"+strconv.Itoa(i+2), matchesToSave[i].HomeGoals)
		f.SetCellValue("Matches", "H"+strconv.Itoa(i+2), matchesToSave[i].AwayGoals)
		f.SetCellValue("Matches", "I"+strconv.Itoa(i+2), matchesToSave[i].HomexA)
		f.SetCellValue("Matches", "J"+strconv.Itoa(i+2), matchesToSave[i].AwayxA)
	}
	for i := range teamsToSave {
		f.SetCellValue("Teams", "A"+strconv.Itoa(i+2), teamsToSave[i].Id)
		f.SetCellValue("Teams", "B"+strconv.Itoa(i+2), teamsToSave[i].Name)
	}
	for i := range playersToSave {
		f.SetCellValue("Players", "A"+strconv.Itoa(i+2), playersToSave[i].Id)
		f.SetCellValue("Players", "B"+strconv.Itoa(i+2), playersToSave[i].Name)
		f.SetCellValue("Players", "C"+strconv.Itoa(i+2), playersToSave[i].Position)
	}
	for i := range appearancesToSave {
		f.SetCellValue("Appearances", "A"+strconv.Itoa(i+2), appearancesToSave[i].MatchId)
		f.SetCellValue("Appearances", "B"+strconv.Itoa(i+2), appearancesToSave[i].Team)
		f.SetCellValue("Appearances", "C"+strconv.Itoa(i+2), appearancesToSave[i].XG)
		f.SetCellValue("Appearances", "D"+strconv.Itoa(i+2), appearancesToSave[i].XA)
		f.SetCellValue("Appearances", "E"+strconv.Itoa(i+2), appearancesToSave[i].Goals)
		f.SetCellValue("Appearances", "F"+strconv.Itoa(i+2), appearancesToSave[i].Assists)
		f.SetCellValue("Appearances", "G"+strconv.Itoa(i+2), appearancesToSave[i].Minutes)
	}

	if err := f.SaveAs("data.xlsx"); err != nil {
		fmt.Println(err)
		return
	}
	return
}

func SaveNewAppearances(appearancesToSave []Appearance, file string) {
	f, err := excelize.OpenFile(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = f.NewSheet("Appearances2")
	if err != nil {
		fmt.Println(err)
		return
	}
	f.SetCellValue("Appearances2", "A1", "MatchId")
	f.SetCellValue("Appearances2", "B1", "Team")
	f.SetCellValue("Appearances2", "C1", "XG")
	f.SetCellValue("Appearances2", "D1", "XA")
	f.SetCellValue("Appearances2", "E1", "Goals")
	f.SetCellValue("Appearances2", "F1", "Assists")
	f.SetCellValue("Appearances2", "G1", "Minutes")
	f.SetCellValue("Appearances2", "H1", "Player")

	for i := range appearancesToSave {
		f.SetCellValue("Appearances2", "A"+strconv.Itoa(i+2), appearancesToSave[i].MatchId)
		f.SetCellValue("Appearances2", "B"+strconv.Itoa(i+2), appearancesToSave[i].Team)
		f.SetCellValue("Appearances2", "C"+strconv.Itoa(i+2), appearancesToSave[i].XG)
		f.SetCellValue("Appearances2", "D"+strconv.Itoa(i+2), appearancesToSave[i].XA)
		f.SetCellValue("Appearances2", "E"+strconv.Itoa(i+2), appearancesToSave[i].Goals)
		f.SetCellValue("Appearances2", "F"+strconv.Itoa(i+2), appearancesToSave[i].Assists)
		f.SetCellValue("Appearances2", "G"+strconv.Itoa(i+2), appearancesToSave[i].Minutes)
		f.SetCellValue("Appearances2", "H"+strconv.Itoa(i+2), appearancesToSave[i].Player)
	}
	if err := f.SaveAs(file); err != nil {
		fmt.Println(err)
		return
	}
	return
}
