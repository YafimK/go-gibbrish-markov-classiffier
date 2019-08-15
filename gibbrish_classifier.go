package main

import (
	"fmt"
	"log"
)

const defaultMinimumProbabilityForPrediction = 0.00005

type GibberishClassifier struct {
	model                           Model
	minimumProbabilityForPrediction float64
}

func (gc GibberishClassifier) isWordGibrish(word string, trace bool) (float64, bool) {
	score := SequenceProbability(gc.model.Chain, word, trace, gc.minimumProbabilityForPrediction)
	isGibberish := false
	if score + gc.model.StdDev < gc.model.Mean {
		isGibberish = true
	}
	normalizedScore := (score - gc.model.Mean) / gc.model.StdDev //the previous method used to calculate if a word is under the threshold
	// isGibberish := normalizedScore < 0
	if trace {
		fmt.Printf("Word: %v | before normalization: %v| after %v\n", word, score, normalizedScore)
	}
	return normalizedScore, isGibberish
}

func (gc GibberishClassifier)  traceIsWordGibrish(word string) (float64, bool) {
	if len(word) < 4 {
		fmt.Printf("Word %v is under the %v char minimum", word, 4)
		return 0, false
	}
	normalizedScore, isGibberish := gc.isWordGibrish(word, true)
	log.Printf("Word: %v | Score: %f | Gibberish: %t\n", word, normalizedScore, isGibberish)
	return normalizedScore, isGibberish
}


func InitGibbrishClassiferFromModelFIle(path string) *GibberishClassifier {

}