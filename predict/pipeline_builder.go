package predict

import (
	"sports-book.com/predict/goals_predictor"
	"sports-book.com/predict/probability_generator"
)

type PipelineBuilder interface {
	Build() Pipeline
	SetPredictor(predictor goals_predictor.GoalsPredictor) PipelineBuilder
	SetProbabilityGenerator(probabilityGenerator probability_generator.ProbabilityGenerator) PipelineBuilder
}

type pipelineBuilderImpl struct {
	predictor            goals_predictor.GoalsPredictor
	probabilityGenerator probability_generator.ProbabilityGenerator
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

func (p *pipelineBuilderImpl) Build() Pipeline {
	return &pipelineImpl{
		predictor:            p.predictor,
		probabilityGenerator: p.probabilityGenerator,
	}
}
