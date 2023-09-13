package predict

import (
	"errors"

	"sports-book.com/pkg/predict/bet_placer"
	"sports-book.com/pkg/predict/goals_predictor"
	"sports-book.com/pkg/predict/probability_generator"
)

var ErrNotEnoughComponents = errors.New("not enough parts to build a pipeline")

type PipelineBuilder interface {
	Build() (Pipeline, error)
	SetPredictor(predictor goals_predictor.GoalsPredictor) PipelineBuilder
	SetProbabilityGenerator(probabilityGenerator probability_generator.ProbabilityGenerator) PipelineBuilder
	SetBetPlacer(betPlacer bet_placer.BetPlacer) PipelineBuilder
}

type pipelineBuilderImpl struct {
	predictor            goals_predictor.GoalsPredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
	betPlacer            bet_placer.BetPlacer
}

func NewPipelineBuilder() PipelineBuilder {
	return &pipelineBuilderImpl{}
}

func (p *pipelineBuilderImpl) SetPredictor(predictor goals_predictor.GoalsPredictor) PipelineBuilder {
	p.predictor = predictor
	return p
}

func (p *pipelineBuilderImpl) SetProbabilityGenerator(probabilityGenerator probability_generator.ProbabilityGenerator) PipelineBuilder {
	p.probabilityGenerator = probabilityGenerator
	return p
}

func (p *pipelineBuilderImpl) SetBetPlacer(betPlacer bet_placer.BetPlacer) PipelineBuilder {
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
