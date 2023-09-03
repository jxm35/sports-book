package goals_predictor

import (
	"fmt"
	"sports-book.com/model"
	"sports-book.com/util"
	"time"
)

type eloGoalsPredictor struct {
	predictions  map[int32]prediction
	FormTimespan int
}

func NewEloGoalsPredictor(formTimespan int) GoalsPredictor {
	return &eloGoalsPredictor{
		predictions:  make(map[int32]prediction),
		FormTimespan: formTimespan,
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
	defer func() {
		if _, exists := e.predictions[matchID]; !exists {
			e.predictions[matchID] = prediction{
				homexG: -1,
				awayxG: -1,
			}
		}
	}()
	avgHomeXg := seasonStats.TotalHomexG / 380
	avgAwayXg := seasonStats.TotalAwayxG / 380
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

	lastGames, _ := util.GetLastXGames(homeTeam, season, date, e.FormTimespan)
	aForm, dForm := e.calcForm(homeTeam, lastGames)

	homeAvgxG := homeSeason.XGScoredAtHome / float64(19)
	homeAvgxGConceded := homeSeason.XGConcededAtHome / float64(19)
	homeAvgxG += aForm
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
	lastGames, _ = util.GetLastXGames(awayTeam, season, date, e.FormTimespan)
	aForm, dForm = e.calcForm(awayTeam, lastGames)

	awayAvgxG := awaySeason.XGScoredAway / float64(19)
	awayAvgxGConceded := awaySeason.XGConcededAway / float64(19)
	awayAvgxG += aForm
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
