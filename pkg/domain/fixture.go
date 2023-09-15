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

type Result struct {
	Id       string      `json:"id"`
	IsResult interface{} `json:"isResult"`
	HomeTeam struct {
		Id         string `json:"id"`
		Title      string `json:"title"`
		ShortTitle string `json:"short_title"`
	} `json:"h"`
	AwayTeam struct {
		Id         string `json:"id"`
		Title      string `json:"title"`
		ShortTitle string `json:"short_title"`
	} `json:"a"`
	Goals struct {
		Home string `json:"h"`
		Away string `json:"a"`
	} `json:"goals"`
	XG struct {
		Home string `json:"h"`
		Away string `json:"a"`
	} `json:"xG"`
	Datetime string `json:"datetime"`
	Forecast struct {
		W string `json:"w"`
		D string `json:"d"`
		L string `json:"l"`
	} `json:"forecast"`
}
