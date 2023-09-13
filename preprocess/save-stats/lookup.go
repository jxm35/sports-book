package save_stats

import (
	"context"
	"fmt"
	"strconv"

	"gorm.io/gorm"

	"sports-book.com/model"
	"sports-book.com/query"
)

func lookupPlayer(excelId string, db *gorm.DB, players []*model.Player) (int32, error) {
	query.SetDefault(db)
	idInt, err := strconv.Atoi(excelId)
	if err != nil {
		return 0, err
	}
	id32 := int32(idInt)

	for _, player := range players {
		if player.ID == id32 {
			// now get the id from the database
			dbPlayer, err := query.Player.WithContext(context.Background()).Where(query.Player.Name.Eq(player.Name)).First()
			if err != nil {
				return 0, err
			}
			return dbPlayer.ID, nil
		}
	}
	return 0, fmt.Errorf("no player found")
}

func lookupMatch(excelId string, db *gorm.DB, matches []*model.Match) (int32, error) {
	query.SetDefault(db)
	idInt, err := strconv.Atoi(excelId)
	if err != nil {
		return 0, err
	}
	id32 := int32(idInt)

	for _, match := range matches {
		if match.ID == id32 {
			// now get the id from the database
			dbPlayer, err := query.Match.WithContext(context.Background()).
				Where(
					query.Match.Date.Eq(match.Date),
					query.Match.HomeTeam.Eq(match.HomeTeam),
					query.Match.AwayTeam.Eq(match.AwayTeam),
				).First()
			if err != nil {
				return 0, err
			}
			return dbPlayer.ID, nil
		}
	}
	return 0, fmt.Errorf("no match found")
}

func lookupTeam(excelId string, db *gorm.DB, teams []*model.Team) (int32, error) {
	query.SetDefault(db)
	idInt, err := strconv.Atoi(excelId)
	if err != nil {
		return 0, err
	}
	id32 := int32(idInt)

	for _, team := range teams {
		if team.ID == id32 {
			// now get the id from the database
			dbTeam, err := query.Team.WithContext(context.Background()).Where(query.Team.Name.Eq(team.Name)).First()
			if err != nil {
				return 0, err
			}
			return dbTeam.ID, nil
		}
	}
	return 0, fmt.Errorf("no team found")
}
