package backtest

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/samber/lo"

	plot2 "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"sports-book.com/model"
	"sports-book.com/predict/domain"
)

func plotDistribution(caliMap map[model.Match]domain.MatchProbability, yearRange string) {
	bucketSize := 0.1
	numBuckets := int(1 / bucketSize)
	buckets := make([]float64, numBuckets)

	matchesHW := make(map[int][]model.Match, numBuckets)
	matchesD := make(map[int][]model.Match, numBuckets)
	matchesAW := make(map[int][]model.Match, numBuckets)

	histoPtsHomeWin := make(plotter.Values, numBuckets)
	histoPtsDraw := make(plotter.Values, numBuckets)
	histoPtsAwayWin := make(plotter.Values, numBuckets)

	// get wins dist
	for match, prediction := range caliMap {
		bucketNum := int(math.Floor(prediction.HomeWin / bucketSize))
		matchesHW[bucketNum] = append(matchesHW[bucketNum], match)

		bucketNum = int(math.Floor(prediction.Draw / bucketSize))
		matchesD[bucketNum] = append(matchesD[bucketNum], match)

		bucketNum = int(math.Floor(prediction.AwayWin / bucketSize))
		matchesAW[bucketNum] = append(matchesAW[bucketNum], match)
	}

	for idx := range buckets {
		histoPtsHomeWin[idx] = float64(len(matchesHW[idx]))
		histoPtsDraw[idx] = float64(len(matchesD[idx]))
		histoPtsAwayWin[idx] = float64(len(matchesAW[idx]))
	}
	blueC := color.RGBA{B: 255, A: 255}
	redC := color.RGBA{R: 255, A: 255}
	greenC := color.RGBA{G: 255, A: 255}

	histoPlot := plot2.New()
	np := vg.Points(10)
	hp, err := plotter.NewBarChart(histoPtsHomeWin, np)
	hp.Color = blueC
	hp.Offset = -np
	dp, err := plotter.NewBarChart(histoPtsDraw, np)
	dp.Color = greenC
	ap, err := plotter.NewBarChart(histoPtsAwayWin, np)
	ap.Color = redC
	ap.Offset = np
	histoPlot.Add(hp, dp, ap)
	histoPlot.Legend.Add("Home Win", hp)
	histoPlot.Legend.Add("Draw", dp)
	histoPlot.Legend.Add("Away Win", ap)
	histoPlot.Legend.Top = true
	histoPlot.Y.Label.Text = "Times Predicted"
	histoPlot.X.Label.Text = "Predicted Likelihood"

	bucketNames := []string{"0-0.1", "0.1-0.2", "0.2-0.3", "0.3-0.4", "0.4-0.5", "0.5-0.6", "0.6-0.7", "0.7-0.8", "0.8-0.9", "0.9-1"}
	histoPlot.NominalX(bucketNames...)

	hwt, err := histoPlot.WriterTo(512, 512, "png")
	if err != nil {
		panic(err)
	}
	hf, err := os.Create(fmt.Sprintf("prediction_distribution_%s.png", yearRange))
	if err != nil {
		panic(err)
	}
	defer hf.Close()
	_, err = hwt.WriteTo(hf)
	if err != nil {
		panic(err)
	}
}

func plotBetDistribution(bets []betResult, yearRange string) {
	bucketSize := 0.05
	numBuckets := 6

	points := make(plotter.Values, numBuckets)

	// get wins dist
	buckets := lo.GroupBy(bets, func(bet betResult) int {
		bookieOdds := 1 / bet.OddsTaken
		var diff float64
		switch bet.Backing {
		case domain.BackHomeWin:
			diff = bet.MatchProbability.HomeWin - bookieOdds
		case domain.BackDraw:
			diff = bet.MatchProbability.Draw - bookieOdds
		case domain.BackAwayWin:
			diff = bet.MatchProbability.AwayWin - bookieOdds
		}
		return int(math.Floor(diff / bucketSize))
	})
	betBuckets := make([]float64, numBuckets) // map from diff to percent won
	for diff, betsAtDiff := range buckets {
		winLoss := lo.GroupBy(betsAtDiff, func(bet betResult) bool {
			return bet.Won
		})
		percentWon := float64(len(winLoss[true])) / float64(len(winLoss[true])+len(winLoss[false]))
		betBuckets[diff] = percentWon
	}

	for i, percent := range betBuckets {
		points[i] = percent
	}

	blueC := color.RGBA{B: 255, A: 255}
	// redC := color.RGBA{R: 255, A: 255}
	// greenC := color.RGBA{G: 255, A: 255}

	histoPlot := plot2.New()
	np := vg.Points(10)
	chart, err := plotter.NewBarChart(points, np)
	chart.Color = blueC
	chart.Offset = -np
	histoPlot.Add(chart)
	histoPlot.Legend.Add("Successful Bets (%)", chart)
	histoPlot.Legend.Top = true
	histoPlot.Y.Label.Text = "Percentage Success"
	histoPlot.X.Label.Text = "Difference in predicted likelihood (my % - bookmaker %)"

	bucketNames := []string{"0-0.05", "0.05-0.1", "0.1-0.15", "0.15-0.2", "0.2-0.25", "0.25-0.3"}
	histoPlot.NominalX(bucketNames...)

	hwt, err := histoPlot.WriterTo(512, 512, "png")
	if err != nil {
		panic(err)
	}
	hf, err := os.Create(fmt.Sprintf("bet_success_distribution_%s.png", yearRange))
	if err != nil {
		panic(err)
	}
	defer hf.Close()
	_, err = hwt.WriteTo(hf)
	if err != nil {
		panic(err)
	}
}

func plotBetWinningsDistribution(bets []betResult, yearRange string) {
	bucketSize := 0.05
	numBuckets := 6

	points := make(plotter.Values, numBuckets)

	// get wins dist
	buckets := lo.GroupBy(bets, func(bet betResult) int {
		bookieOdds := 1 / bet.OddsTaken
		var diff float64
		switch bet.Backing {
		case domain.BackHomeWin:
			diff = bet.MatchProbability.HomeWin - bookieOdds
		case domain.BackDraw:
			diff = bet.MatchProbability.Draw - bookieOdds
		case domain.BackAwayWin:
			diff = bet.MatchProbability.AwayWin - bookieOdds
		}
		return int(math.Floor(diff / bucketSize))
	})
	betBuckets := make([]float64, numBuckets) // map from diff to percent won
	for diff, betsAtDiff := range buckets {
		amountWon := lo.Reduce(betsAtDiff, func(agg float64, bet betResult, index int) float64 {
			if bet.Won {
				return agg + (bet.Amount * bet.OddsTaken) - bet.Amount
			} else {
				return agg - bet.Amount
			}
		}, 0.0)

		betBuckets[diff] = amountWon
	}

	for i, percent := range betBuckets {
		points[i] = percent
	}

	blueC := color.RGBA{B: 255, A: 255}
	// redC := color.RGBA{R: 255, A: 255}
	// greenC := color.RGBA{G: 255, A: 255}

	histoPlot := plot2.New()
	np := vg.Points(10)
	chart, err := plotter.NewBarChart(points, np)
	chart.Color = blueC
	chart.Offset = -np
	histoPlot.Add(chart)
	histoPlot.Legend.Add("Amount Won/Lost", chart)
	histoPlot.Legend.Top = true
	histoPlot.Y.Label.Text = "Amount Won (GBP)"
	histoPlot.X.Label.Text = "Difference in predicted likelihood (my % - bookmaker %)"

	bucketNames := []string{"0-0.05", "0.05-0.1", "0.1-0.15", "0.15-0.2", "0.2-0.25", "0.25-0.3"}
	histoPlot.NominalX(bucketNames...)

	hwt, err := histoPlot.WriterTo(512, 512, "png")
	if err != nil {
		panic(err)
	}
	hf, err := os.Create(fmt.Sprintf("bet_win_distribution_%s.png", yearRange))
	if err != nil {
		panic(err)
	}
	defer hf.Close()
	_, err = hwt.WriteTo(hf)
	if err != nil {
		panic(err)
	}
}

func plotBetsPlacedDistribution(bets []betResult, yearRange string) {
	bucketSize := 0.05
	numBuckets := 6

	points := make(plotter.Values, numBuckets)

	// get wins dist
	buckets := lo.GroupBy(bets, func(bet betResult) int {
		bookieOdds := 1 / bet.OddsTaken
		var diff float64
		switch bet.Backing {
		case domain.BackHomeWin:
			diff = bet.MatchProbability.HomeWin - bookieOdds
		case domain.BackDraw:
			diff = bet.MatchProbability.Draw - bookieOdds
		case domain.BackAwayWin:
			diff = bet.MatchProbability.AwayWin - bookieOdds
		}
		return int(math.Floor(diff / bucketSize))
	})
	betBuckets := make([]float64, numBuckets) // map from diff to percent won
	for diff, betsAtDiff := range buckets {
		betBuckets[diff] = float64(len(betsAtDiff))
	}

	for i, percent := range betBuckets {
		points[i] = percent
	}

	blueC := color.RGBA{B: 255, A: 255}
	// redC := color.RGBA{R: 255, A: 255}
	// greenC := color.RGBA{G: 255, A: 255}

	histoPlot := plot2.New()
	np := vg.Points(10)
	chart, err := plotter.NewBarChart(points, np)
	chart.Color = blueC
	chart.Offset = -np
	histoPlot.Add(chart)
	histoPlot.Legend.Add("Bets Placed", chart)
	histoPlot.Legend.Top = true
	histoPlot.Y.Label.Text = "Bets Placed"
	histoPlot.X.Label.Text = "Difference in predicted likelihood (my % - bookmaker %)"

	bucketNames := []string{"0-0.05", "0.05-0.1", "0.1-0.15", "0.15-0.2", "0.2-0.25", "0.25-0.3"}
	histoPlot.NominalX(bucketNames...)

	hwt, err := histoPlot.WriterTo(512, 512, "png")
	if err != nil {
		panic(err)
	}
	hf, err := os.Create(fmt.Sprintf("bets_placed_distribution_%s.png", yearRange))
	if err != nil {
		panic(err)
	}
	defer hf.Close()
	_, err = hwt.WriteTo(hf)
	if err != nil {
		panic(err)
	}
}
