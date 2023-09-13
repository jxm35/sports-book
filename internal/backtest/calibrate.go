package backtest

import (
	"fmt"
	"image/color"
	"math"
	"os"

	plot2 "gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"

	"sports-book.com/pkg/model"
	"sports-book.com/pkg/predict/domain"
)

func getCalibration(
	caliMap map[model.Match]domain.MatchProbability,
	yearString string,
) {
	bucketSize := 0.1
	numBuckets := int(1 / bucketSize)

	matchesHW := make(map[int][]model.Match, numBuckets)
	bucketsHW := make([]float64, numBuckets)
	HWCount := 0

	matchesD := make(map[int][]model.Match, numBuckets)
	bucketsD := make([]float64, numBuckets)
	DCount := 0

	matchesAW := make(map[int][]model.Match, numBuckets)
	bucketsAW := make([]float64, numBuckets)
	AWCount := 0

	for match, prediction := range caliMap {
		bucketNum := int(math.Floor(prediction.HomeWin / bucketSize))
		matchesHW[bucketNum] = append(matchesHW[bucketNum], match)

		bucketNum = int(math.Floor(prediction.Draw / bucketSize))
		matchesD[bucketNum] = append(matchesD[bucketNum], match)

		bucketNum = int(math.Floor(prediction.AwayWin / bucketSize))
		matchesAW[bucketNum] = append(matchesAW[bucketNum], match)
	}

	perfectPts := make(plotter.XYs, numBuckets+1)
	ptsHW := make(plotter.XYs, numBuckets)
	ptsD := make(plotter.XYs, numBuckets)
	ptsAW := make(plotter.XYs, numBuckets)

	for idx := range bucketsHW {
		var HWCorrect, HWIncorrect, DCorrect, DIncorrect, AWCorrect, AWIncorrect float64

		for _, match := range matchesHW[idx] {
			if match.HomeGoals > match.AwayGoals {
				HWCorrect += 1
				HWCount += 1
			} else {
				HWIncorrect += 1
			}
		}

		for _, match := range matchesD[idx] {
			if match.HomeGoals == match.AwayGoals {
				DCorrect += 1
				DCount += 1
			} else {
				DIncorrect += 1
			}
		}
		for _, match := range matchesAW[idx] {
			if match.HomeGoals < match.AwayGoals {
				AWCorrect += 1
				AWCount += 1
			} else {
				AWIncorrect += 1
			}
		}

		bucketsHW[idx] = HWCorrect / (HWCorrect + HWIncorrect)
		bucketsD[idx] = DCorrect / (DCorrect + DIncorrect)
		bucketsAW[idx] = AWCorrect / (AWCorrect + AWIncorrect)

		perfectPts[idx].X = float64(idx) * bucketSize
		perfectPts[idx].Y = float64(idx) * bucketSize

		if bucketsHW[idx] != 0 {
			ptsHW[idx].X = perfectPts[idx].X
			ptsHW[idx].Y = bucketsHW[idx]
			if math.IsNaN(ptsHW[idx].Y) {
				ptsHW[idx].Y = 0
			}
		}

		if bucketsD[idx] != 0 {
			ptsD[idx].X = perfectPts[idx].X
			ptsD[idx].Y = bucketsD[idx]
			if math.IsNaN(ptsD[idx].Y) {
				ptsD[idx].Y = 0
			}
		}

		if bucketsAW[idx] != 0 {
			ptsAW[idx].X = perfectPts[idx].X
			ptsAW[idx].Y = bucketsAW[idx]
			if math.IsNaN(ptsAW[idx].Y) {
				ptsAW[idx].Y = 0
			}
		}

	}
	perfectPts[10].X = 1
	perfectPts[10].Y = 1
	blueC := color.RGBA{B: 255, A: 255}
	redC := color.RGBA{R: 255, A: 255}
	greenC := color.RGBA{G: 255, A: 255}
	blackC := color.RGBA{G: 0, R: 0, B: 0, A: 255}

	pl, _ := plotter.NewLine(perfectPts)
	pl.Color = blackC

	hw, _ := plotter.NewLine(ptsHW)
	hw.Color = blueC

	d, _ := plotter.NewLine(ptsD)
	d.Color = greenC

	aw, _ := plotter.NewLine(ptsAW)
	aw.Color = redC

	linePlot := plot2.New()
	linePlot.Add(pl, hw, d, aw)

	linePlot.Title.Text = "Model Calibration"
	linePlot.X.Label.Text = "Predicted Likelihood of Outcome"
	linePlot.Y.Label.Text = "Realised Likelihood of Outcome"

	linePlot.Legend.Add("Perfect Model", pl)
	linePlot.Legend.Add("Home Wins", hw)
	linePlot.Legend.Add("Draws", d)
	linePlot.Legend.Add("Away Wins", aw)
	linePlot.Legend.Top = true

	// calculate rms error
	var diff float64
	diff = 0
	count := 0
	for _, pt := range ptsHW {
		if pt.Y != 0 {
			diff += math.Pow(pt.X-pt.Y, 2)
			count++
		}
	}
	diff = diff / float64(count)
	rmseHW := math.Sqrt(diff)

	diff = 0
	count = 0
	for _, pt := range ptsD {
		if pt.Y != 0 {
			diff += math.Pow(pt.X-pt.Y, 2)
			count++
		}
	}
	diff = diff / float64(count)
	rmseD := math.Sqrt(diff)

	diff = 0
	count = 0
	for _, pt := range ptsAW {
		if pt.Y != 0 {
			diff += math.Pow(pt.X-pt.Y, 2)
			count++
		}
	}
	diff = diff / float64(count)
	rmseAW := math.Sqrt(diff)

	labels, err := plotter.NewLabels(plotter.XYLabels{
		XYs: []plotter.XY{
			{X: +0.5, Y: 0.3},
			{X: +0.5, Y: 0.275},
			{X: +0.5, Y: 0.25},
			{X: +0.5, Y: 0.225},
			{X: +0.5, Y: 0.2},
			{X: +0.5, Y: 0.175},
			{X: +0.5, Y: 0.15},
		},
		Labels: []string{
			fmt.Sprintf("Matches Predicted: %d", len(caliMap)),
			fmt.Sprintf("RMSE (home): %f", rmseHW),
			fmt.Sprintf("home win occurences: %d", HWCount),
			fmt.Sprintf("RMSE (draw): %f", rmseD),
			fmt.Sprintf("draw occurences: %d", DCount),
			fmt.Sprintf("RMSE (away): %f", rmseAW),
			fmt.Sprintf("away win occurences: %d", AWCount),
		},
	})

	linePlot.Add(labels)

	wt, err := linePlot.WriterTo(512, 512, "png")
	if err != nil {
		panic(err)
	}
	f, err := os.Create(fmt.Sprintf("analysis/calibrations/calibration_%s.png", yearString))
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = wt.WriteTo(f)
	if err != nil {
		panic(err)
	}
}
