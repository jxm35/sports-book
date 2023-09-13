package scrape_stats

type Match struct {
	Id        int
	Date      int64
	HomeTeam  int
	AwayTeam  int
	HomexG    string
	AwayxG    string
	HomeGoals int
	AwayGoals int
	HomexA    string
	AwayxA    string
}

type Appearance struct {
	Player  int
	MatchId int
	Team    int
	XG      float32
	XA      float32
	Goals   int
	Assists int
	Minutes float32
}

type Team struct {
	Id   int
	Name string
}

type Player struct {
	Id       int
	Name     string
	Position string
}
