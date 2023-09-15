package save_stats

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	model2 "sports-book.com/pkg/gorm/model"
	"sports-book.com/pkg/gorm/query"
)

func saveAppearances(db *gorm.DB, apps []*model2.Appearance) {
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

func saveMatches(db *gorm.DB, matches []*model2.Match) {
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

func savePlayers(db *gorm.DB, players []*model2.Player) {
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

func saveTeams(db *gorm.DB, teams []*model2.Team) {
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
