package main

import (
	"fmt"
	"time"

	"github.com/ziyasal/neurolog/neurolog"
)

func addition(sampleFilePath string) {

	options := neurolog.Options{
		Name:            "additions",
		Type:            "REGRESSOR",
		Inputs:          []string{"number1", "number2"},
		Outputs:         []string{"result"},
		HiddenLayers:    []int{3},
		DatasetSize:     50,
		TestDatasetSize: 10,
		RedisHost:       "localhost:6379",
		AutoCreate:      true,
		Normalize:       true,
		Prefix:          "nn:",
	}

	network := neurolog.New(options)

	//Load training data set
	csvData := loadAdditionCsvData(sampleFilePath)

	trainDataSet := csvData[:90]
	testDataSet := csvData[90:]

	//Training data-set
	for _, pair := range trainDataSet {
		network.ObserveTrain(pair.Input, pair.Output)
	}

	network.Info()

	//Testing data-set
	for _, element := range testDataSet {
		network.ObserveTest(element.Input, element.Output)
	}
	//Train
	network.Train(0, 0, true, true)
	for network.IsTraining() == true {
		fmt.Println("Training...")
		time.Sleep(1)
	}

	//Test
	errors := 0
	for _, pair := range testDataSet {
		input := pair.Input
		output := pair.Output
		networkOutput := network.Run(input)

		s := fmt.Sprintf("Neural Network calculation %d + %d = %d",
			input["number1"],
			input["number2"],
			int64(networkOutput["result"]))
		fmt.Println(s)

		result := int64(networkOutput["result"])
		if result != output["result"] {
			errors += 1
		}
	}

	s := fmt.Sprintf("%s prediction errors on %s test items", errors, len(testDataSet))
	fmt.Println(s)
}
