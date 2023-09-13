package save_stats

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	model2 "sports-book.com/pkg/db_model"
	"sports-book.com/pkg/db_query"
)

func saveAppearances(db *gorm.DB, apps []*model2.Appearance) {
	db_query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := db_query.Appearance.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).CreateInBatches(apps, 200)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func saveMatches(db *gorm.DB, matches []*model2.Match) {
	db_query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := db_query.Match.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(matches...)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func savePlayers(db *gorm.DB, players []*model2.Player) {
	db_query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := db_query.Player.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(players...)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func saveTeams(db *gorm.DB, teams []*model2.Team) {
	db_query.SetDefault(db)

	// Update all columns, except primary keys, to new value on conflict
	err := db_query.Team.WithContext(context.Background()).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(teams...)
	if err != nil {
		fmt.Println(err)
		return
	}
}
