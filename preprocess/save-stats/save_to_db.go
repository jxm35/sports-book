package save_stats

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"sports-book.com/model"
	"sports-book.com/query"
)

func saveAppearances(db *gorm.DB, apps []*model.Appearance) {
	query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := query.Appearance.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(apps, 200)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func saveMatches(db *gorm.DB, matches []*model.Match) {
	query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := query.Match.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(matches...)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func savePlayers(db *gorm.DB, players []*model.Player) {
	query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := query.Player.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(players...)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func saveTeams(db *gorm.DB, teams []*model.Team) {
	query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := query.Team.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(teams...)
	if err != nil {
		fmt.Println(err)
		return
	}
}
