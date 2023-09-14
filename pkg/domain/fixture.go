package domain

type Fixture struct {
	AwayTeam struct {
		UsId       string `json:"id"`
		ShortTitle string `json:"short_title"`
		Title      string `json:"title"`
	} `json:"a"`
	Datetime string `json:"datetime"`
	HomeTeam struct {
		UsId       string `json:"id"`
		ShortTitle string `json:"short_title"`
		Title      string `json:"title"`
	} `json:"h"`
	Id       string      `json:"id"`
	IsResult interface{} `json:"isResult"`
}
