package main

import (
	"flag"
	"fmt"
	"github.houston.softwaregrp.net/Performance-Engineering/go-gibbrish-markov-chain-classifier/model"
	"log"
	"os"
)

func main() {
	train := flag.NewFlagSet("train", flag.ExitOnError)
	checkString := flag.NewFlagSet("check", flag.ExitOnError)
	testList := flag.NewFlagSet("checkFileList", flag.ExitOnError)
	flag.Parse()
	args := flag.Args()
	switch args[0] {
	case "train":
		err := train.Parse(nil)
		if err != nil {
			log.Fatalln(err)
		}
	case "check":
		if len(args) < 2 {
			log.Fatalf("Please provide word to manually test ")
		}
		err := checkString.Parse(args[1:])
		if err != nil {
			log.Fatalln(err)
		}
	case "checkFileList":
		if len(args) < 2 {
			log.Fatalf("Please provide file name with words seperated by line delimeter" +
				" to manually test ")
		}
		err := testList.Parse(args[1:])
		if err != nil {
			log.Fatalln(err)
		}
	default:
		fmt.Printf("%q is not valid command.\n", args[1])
		os.Exit(2)
	}

	if train.Parsed() {
		model := model.BuildModel()
		model.SaveModelToJson(model)
	} else if checkString.Parsed() {
		word := checkString.Args()[0]
		if len(word) == 0 {
			fmt.Printf("Please write a word after the command to check a string")
			return
		}
		model, err := loadModel()
		if err != nil {
			fmt.Println(err)
			return
		}
		traceIsWordGibrish(model, word)
	} else if testList.Parsed() {
		inputFileName := testList.Args()[0]
		words, err := getDataset(inputFileName)
		if err != nil {
			log.Fatalln(err)
		}
		model, err := loadModel()
		if err != nil {
			fmt.Println(err)
			return
		}
		var nonGibbrishWords []string
		var gibbrishWords []string
		var underMinimalLengthWords = 9
		for _, word := range words {
			if len(word) < 4 {
				underMinimalLengthWords++
				continue
			}
			_, isGibberish := model.isWordGibrish(word, false)
			if isGibberish {
				gibbrishWords = append(gibbrishWords, word)
			} else {
				nonGibbrishWords = append(nonGibbrishWords, word)
			}
		}
		log.Printf("Found %v words out of %v | %v words weren't scanned since they are under the minimal length \n", len(gibbrishWords), len(words), underMinimalLengthWords)
		if len(gibbrishWords) < 20 {
			log.Printf("gibbrish words: %v\n", gibbrishWords)
		}
		if len(nonGibbrishWords) < 20 {
			log.Printf("Non gibbrish words: %v\n", nonGibbrishWords)
		}
	}
}
