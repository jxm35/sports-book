package pipeline

import (
	"errors"

	"sports-book.com/pkg/bet_selector"
	"sports-book.com/pkg/probability_generator"
	"sports-book.com/pkg/score_predictor"
)

var ErrNotEnoughComponents = errors.New("not enough parts to build a pipeline")

type PipelineBuilder interface {
	Build() (Pipeline, error)
	SetPredictor(predictor score_predictor.ScorePredictor) PipelineBuilder
	SetProbabilityGenerator(probabilityGenerator probability_generator.ProbabilityGenerator) PipelineBuilder
	SetBetPlacer(betPlacer bet_selector.BetSelector) PipelineBuilder
}

type pipelineBuilderImpl struct {
	predictor            score_predictor.ScorePredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
	betPlacer            bet_selector.BetSelector
}

func NewPipelineFromConfig() (Pipeline, error) {
	predictor, err := score_predictor.NewScorePredictorFromConfig()
	if err != nil {
		return nil, err
	}
	betSelector, err := bet_selector.NewBetSelectorFromConfig()
	if err != nil {
		return nil, err
	}
	return &pipelineImpl{
		predictor:            predictor,
		probabilityGenerator: &probability_generator.WeibullOddsGenerator{},
		betPlacer:            betSelector,
	}, nil
}

func NewPipelineBuilder() PipelineBuilder {
	return &pipelineBuilderImpl{}
}

func (p *pipelineBuilderImpl) SetPredictor(predictor score_predictor.ScorePredictor) PipelineBuilder {
	p.predictor = predictor
	return p
}

func (p *pipelineBuilderImpl) SetProbabilityGenerator(probabilityGenerator probability_generator.ProbabilityGenerator) PipelineBuilder {
	p.probabilityGenerator = probabilityGenerator
	return p
}

func (p *pipelineBuilderImpl) SetBetPlacer(betPlacer bet_selector.BetSelector) PipelineBuilder {
	p.betPlacer = betPlacer
	return p
}

func (p *pipelineBuilderImpl) Build() (Pipeline, error) {
	if p.predictor == nil || p.probabilityGenerator == nil || p.betPlacer == nil {
		return nil, ErrNotEnoughComponents
	}

	return &pipelineImpl{
		predictor:            p.predictor,
		probabilityGenerator: p.probabilityGenerator,
		betPlacer:            p.betPlacer,
	}, nil
}
