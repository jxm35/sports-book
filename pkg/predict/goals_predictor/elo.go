package goals_predictor

import (
	"fmt"
	"time"

	"sports-book.com/pkg/model"
	"sports-book.com/pkg/util"
)

type eloGoalsPredictor struct {
	predictions               map[int32]prediction
	FormTimespan              int
	chanceConversionWeighting int
}

func NewEloGoalsPredictor(formTimespan, chanceConversionRating int) GoalsPredictor {
	return &eloGoalsPredictor{
		predictions:               make(map[int32]prediction),
		FormTimespan:              formTimespan,
		chanceConversionWeighting: chanceConversionRating,
	}
}

func (e *eloGoalsPredictor) PredictScore(homeTeam, awayTeam, season int32, league string, date time.Time, matchID int32) (float64, float64, error) {
	seasonStats, err := util.GetSeasonDetails(season-1, league)
	if err != nil {
		return -1, -1, err
	}
	if seasonStats.TotalHG == 0 && seasonStats.TotalAG == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if seasonStats.MatchCount == 0 || seasonStats.MatchCount%2 != 0 {
		return -1, -1, ErrInvalidSeason
	}
	defer func() {
		if _, exists := e.predictions[matchID]; !exists {
			e.predictions[matchID] = prediction{
				homexG: -1,
				awayxG: -1,
			}
		}
	}()
	avgHomeXg := seasonStats.TotalHomexG / float64(seasonStats.MatchCount)
	avgAwayXg := seasonStats.TotalAwayxG / float64(seasonStats.MatchCount)
	avgHomeGoalsConceded := avgAwayXg
	avgAwayGoalsConceded := avgHomeXg

	// calculate home team's strengths
	homeSeason, err := util.GetTeamSeasonDetails(season-1, homeTeam)
	if err != nil {
		return -1, -1, err
	}
	if homeSeason.XGScoredAtHome == 0 && homeSeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if homeSeason.HomeCount == 0 || homeSeason.AwayCount == 0 || homeSeason.HomeCount != homeSeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}

	lastGames, _ := util.GetHomeLastXGames(homeTeam, season, date, e.FormTimespan)
	aForm, dForm := e.calcFormHome(homeTeam, lastGames)
	homeChanceConversion := e.calcGoalConversion(homeSeason, homeSeason.HomeCount+homeSeason.AwayCount)

	homeAvgxG := homeSeason.XGScoredAtHome / float64(homeSeason.AwayCount)
	homeAvgxGConceded := homeSeason.XGConcededAtHome / float64(homeSeason.AwayCount)
	homeAvgxG += aForm
	if e.chanceConversionWeighting != 0 {
		homeAvgxG += homeChanceConversion / float64(e.chanceConversionWeighting)
	}
	homeAvgxGConceded += dForm
	homeAttackStrength := homeAvgxG / avgHomeXg
	homeDefenseStrength := homeAvgxGConceded / avgHomeGoalsConceded

	// calculate away team's strengths
	awaySeason, err := util.GetTeamSeasonDetails(season-1, awayTeam)
	if err != nil {
		return -1, -1, err
	}
	if awaySeason.XGScoredAtHome == 0 && awaySeason.XGConcededAtHome == 0 {
		return -1, -1, ErrNoPreviousData
	}
	if awaySeason.HomeCount == 0 || awaySeason.AwayCount == 0 || awaySeason.HomeCount != awaySeason.AwayCount {
		return -1, -1, ErrInvalidSeason
	}

	lastGames, _ = util.GetAwayLastXGames(awayTeam, season, date, e.FormTimespan)
	aForm, dForm = e.calcFormAway(awayTeam, lastGames)
	awayChanceConversion := e.calcGoalConversion(awaySeason, awaySeason.HomeCount+awaySeason.AwayCount)

	awayAvgxG := awaySeason.XGScoredAway / float64(awaySeason.AwayCount)
	awayAvgxGConceded := awaySeason.XGConcededAway / float64(awaySeason.AwayCount)
	awayAvgxG += aForm
	if e.chanceConversionWeighting != 0 {
		awayAvgxG += awayChanceConversion / float64(e.chanceConversionWeighting)
	}
	awayAvgxGConceded += dForm
	awayDefenseStrength := awayAvgxGConceded / avgAwayGoalsConceded
	awayAttackStrength := awayAvgxG / avgAwayXg

	// use strengths to project home and away goals
	projectedHomeGoals := homeAttackStrength * awayDefenseStrength * avgHomeXg
	projectedAwayGoals := awayAttackStrength * homeDefenseStrength * avgAwayXg

	fmt.Printf("%d: %f | %d: %f", homeTeam, projectedHomeGoals, awayTeam, projectedAwayGoals)
	e.predictions[matchID] = prediction{
		homexG: projectedHomeGoals,
		awayxG: projectedAwayGoals,
	}
	return projectedHomeGoals, projectedAwayGoals, nil
}

type prediction struct {
	homexG float64
	awayxG float64
}

// additive or multiplicative?
// calculate home and away separately or together?

func (e *eloGoalsPredictor) calcFormHome(team int32, matches []model.Match) (float64, float64) {
	attackForm := 0.0  // how many more goals they score than expected
	defenseForm := 0.0 // how many more goals they concede than expected
	count := 0
	for i := 0; i < len(matches); i++ {
		pred, ok := e.predictions[matches[i].ID]
		if !ok {
			panic("could not find prediction")
		}
		if pred.awayxG == -1 && pred.homexG == -1 {
			continue
		}
		attackForm = attackForm + (matches[i].HomeExpectedGoals - pred.homexG)
		defenseForm = defenseForm + (matches[i].AwayExpectedGoals - pred.awayxG)
		count++
	}
	if count < 4 {
		return 0, 0
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)

	return attackForm, defenseForm
}

func (e *eloGoalsPredictor) calcFormAway(team int32, matches []model.Match) (float64, float64) {
	attackForm := 0.0
	defenseForm := 0.0
	count := 0
	for i := 0; i < len(matches); i++ {
		pred, ok := e.predictions[matches[i].ID]
		if !ok {
			panic("could not find prediction")
		}
		if pred.awayxG == -1 && pred.homexG == -1 {
			continue
		}
		attackForm = attackForm + (matches[i].AwayExpectedGoals - pred.awayxG)
		defenseForm = defenseForm + (matches[i].HomeExpectedGoals - pred.homexG)
		count++
	}
	if count < 4 {
		return 0, 0
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)
	return attackForm, defenseForm
}

func (e *eloGoalsPredictor) calcForm(team int32, matches []model.Match) (float64, float64) {
	attackForm := 0.0
	defenseForm := 0.0
	count := 0
	for _, match := range matches {
		pred, ok := e.predictions[match.ID]
		if !ok {
			panic("could not find prediction")
		}
		if pred.awayxG == -1 && pred.homexG == -1 {
			continue
		}
		if match.HomeTeam == team {
			attackForm = attackForm + (match.HomeExpectedGoals - pred.homexG)
			defenseForm = defenseForm + (match.AwayExpectedGoals - pred.awayxG)
		} else if match.AwayTeam == team {
			attackForm = attackForm + (match.AwayExpectedGoals - pred.awayxG)
			defenseForm = defenseForm + (match.HomeExpectedGoals - pred.homexG)
		} else {
			panic("team not in match")
		}
		count++
	}
	if count < 4 {
		return 0, 0
	}
	attackForm = attackForm / float64(count)
	defenseForm = defenseForm / float64(count)
	return attackForm, defenseForm
}

func (e *eloGoalsPredictor) calcGoalConversion(lastSeason util.TeamSeasonDetails, matchesPlayed int) float64 {
	diff := float64(lastSeason.GoalsScoredAtHome+lastSeason.GoalsScoredAway) - (lastSeason.XGScoredAtHome + lastSeason.XGScoredAway)
	return diff / float64(matchesPlayed)
}
