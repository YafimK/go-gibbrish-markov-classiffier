package model

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/mb-14/gomarkov"
	"github.com/montanaflynn/stats"
	"io/ioutil"
	"log"
	"math"
	"os"
	"regexp"
	"strings"
)

const defaultMinimumProbabilityForTraining = 0.05
const defaultMinimumProbabilityForPrediction = 0.00005


func buildChain(sourceFiles []string) *gomarkov.Chain {
	chain := gomarkov.NewChain(2)
	for _, sourceFile := range sourceFiles {
		//TODO: add check that path exists!
		words, err := getDataset(sourceFile)
		if err != nil {
			log.Fatalln(err)
		}
		for _, data := range words {
			data = strings.ToLower(data)
			parsedWord := sanitizeString(data)
			if len(parsedWord) == 0 {
				continue
			}
			chain.Add(split(data))
		}
	}
	return chain
}

func getModelScoresForFile(chain *gomarkov.Chain, filepath string, minimumProbability float64) []float64 {
	scores := make([]float64, 0)
	words, err := getDataset(filepath)
	if err != nil {
		panic(err)
	}
	for _, word := range words {
		word = strings.ToLower(word)

		if len(word) < 4 {
			continue
		}
		score := sequenceProbability(chain, word, false, minimumProbability)
		if math.IsNaN(score) {
			score = minimumProbability
			continue
		}
		scores = append(scores, score)
	}
	return scores
}

func calculateStatsForSource(wordScores []float64) (mean float64, stdDev float64, err error) {
	stdDev, err = stats.StandardDeviation(wordScores)
	if err != nil {
		return mean, stdDev, err
	}
	mean, err = stats.Mean(wordScores)
	if err != nil {
		return mean, stdDev, err
	}
	return mean, stdDev, err
}

func BuildModel() Model {
	var model Model
	sourceDatasetFiles := []string{"words.txt", "glove.6B.50d-words.txt", "bigGoodEnglishWords.txt"}
	model.Chain = buildChain(sourceDatasetFiles)

	scores := getModelScoresForFile(model.Chain, "basic2.txt", defaultMinimumProbabilityForTraining)
	var err error
	correlationsMean, correlationsStdDev, err := calculateStatsForSource(scores)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("stats on correlations candidate[%v] examples %v scores | Mean: %v | stdDev: %v \n","basic2.txt", len(scores),
		correlationsMean, correlationsStdDev)



	scores = getModelScoresForFile(model.Chain, "bigGoodEnglishWords.txt", defaultMinimumProbabilityForTraining)
	bigGoodEnglishMean, bigGoodEnglishStdDev, err := calculateStatsForSource(scores)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("stats on correlations candidate[%v] examples %v scores | Mean: %v | stdDev: %v \n","bigGoodEnglishWords.txt", len(scores),
		bigGoodEnglishMean, bigGoodEnglishStdDev)

	scores = getModelScoresForFile(model.Chain, "falsepositive.txt", defaultMinimumProbabilityForTraining)
	falsepositiveMean, falsepositiveStdDev, err := calculateStatsForSource(scores)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("stats on correlations candidate[%v] examples %v scores | Mean: %v | stdDev: %v \n","falsepositive.txt", len(scores),
		falsepositiveMean, falsepositiveStdDev)

	//We currently define the desired threshold using correlations... we can change this to any other calculated stats
	model.StdDev = correlationsStdDev
	model.Mean = correlationsMean

	return model
}

func SaveModelToJson(model Model) {
	jsonObj, _ := json.Marshal(model)
	err := ioutil.WriteFile("model.json", jsonObj, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func LoadModelFromJson() (Model, error) {
	data, err := ioutil.ReadFile("model.json")
	if err != nil {
		return Model{}, err
	}
	var m Model
	err = json.Unmarshal(data, &m)
	if err != nil {
		return Model{}, err
	}
	return m, nil
}

func sanitizeString(value string) string {
	reg, err := regexp.Compile("[^a-zA-Z]+")
	if err != nil {
		log.Fatal(err)
	}
	return reg.ReplaceAllString(value, "")
}

func getDataset(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	var list []string
	for scanner.Scan() {
		splitString := strings.Fields(scanner.Text())
		list = append(list, splitString...)
	}
	return list, nil
}

func split(str string) []string {
	return strings.Split(str, "")
}

func sequenceProbability(chain *gomarkov.Chain, input string, trace bool, unseenPairProbability float64) float64 {
	// sanitized := santizeString(input)
	input = strings.ToLower(input)
	tokens := split(input)
	logProb := float64(0)
	pairs := gomarkov.MakePairs(tokens, chain.Order)
	for _, pair := range pairs {
		prob, _ := chain.TransitionProbability(pair.NextState, pair.CurrentState)
		if trace {
			log.Println(pair, prob)
		}
		if prob > 0 {
			logProb += math.Log10(prob)
		} else {
			logProb += math.Log10(unseenPairProbability)
		}
	}
	return math.Pow(10, logProb/float64(len(pairs)))
}
