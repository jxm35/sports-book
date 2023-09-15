// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNamePrediction = "predictions"

// Prediction mapped from table <predictions>
type Prediction struct {
	ID      int32   `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	MatchID int32   `gorm:"column:match_id;not null" json:"match_id"`
	HomeXg  float64 `gorm:"column:home_xg;not null" json:"home_xg"`
	AwayXg  float64 `gorm:"column:away_xg;not null" json:"away_xg"`
}

// TableName Prediction's table name
func (*Prediction) TableName() string {
	return TableNamePrediction
}