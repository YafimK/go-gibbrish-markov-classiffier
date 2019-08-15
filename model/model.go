package model

import (
	"github.com/mb-14/gomarkov"
)

type Model struct {
	Mean   float64         `json:"mean"`
	StdDev float64         `json:"std_dev"`
	Chain  *gomarkov.Chain `json:"chain"`
	MinimumProbabilityForTraining float64 `json:"mean"`
	MinimumProbabilityForPrediction float64 `json:"mean"`
}


