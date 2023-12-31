// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameOdds1x2 = "odds_1x2"

// Odds1x2 mapped from table <odds_1x2>
type Odds1x2 struct {
	ID        int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Bookmaker string  `gorm:"column:bookmaker;not null" json:"bookmaker"`
	Match     int32   `gorm:"column:match;not null" json:"match"`
	HomeWin   float64 `gorm:"column:home_win;not null" json:"home_win"`
	Draw      float64 `gorm:"column:draw;not null" json:"draw"`
	AwayWin   float64 `gorm:"column:away_win;not null" json:"away_win"`
}

// TableName Odds1x2's table name
func (*Odds1x2) TableName() string {
	return TableNameOdds1x2
}
