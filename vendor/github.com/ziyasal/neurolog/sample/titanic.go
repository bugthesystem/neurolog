package main

import (
	"fmt"
	"time"

	"github.com/ziyasal/neurolog/neurolog"
)

func titanic(sampleFilePath string) {
	options := neurolog.Options{
		Name:            "titanic",
		Type:            "classifier",
		Inputs:          []string{"class1", "class2", "class3", "female", "male", "age", "sibsp", "parch", "fare"},
		Outputs:         []string{"dead", "alive"},
		HiddenLayers:    []int{15},
		DatasetSize:     1000,
		TestDatasetSize: 500,
		RedisHost:       "localhost:6379",
	}

	network := neurolog.New(options)

	//Load training data set
	csvData := loadTitanicCsvData(sampleFilePath)

	trainDataSet := csvData[:600]
	testDataSet := csvData[600:]

	//Training data-set
	for _, pair := range trainDataSet {
		network.ObserveTrain(pair.Input, pair.Output)
	}

	//Testing data-set
	for _, pair := range testDataSet {
		network.ObserveTest(pair.Input, pair.Output)
	}
	//Train
	network.Train(0, 0, true, false)
	for network.IsTraining() {
		fmt.Println("Training...")
		time.Sleep(1)
	}

	//Test
	errors := 0
	for _, pair := range testDataSet {
		output := network.Classify(pair.Input)

		//result, _ := strconv.ParseInt(output["result"], 10, 64)
		//if _, ok := pair.Output[result]; !ok {
		//	errors += 1
		//}
		fmt.Println(output)
	}

	fmt.Sprintf("%s prediction errors on %s test items", errors, len(testDataSet))
}
