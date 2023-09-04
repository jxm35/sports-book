package backtest

import (
	"fmt"
	"image/color"
	"math"
	"os"

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
