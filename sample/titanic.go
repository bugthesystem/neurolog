package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/ziyasal/neurolog/neurolog"
)

func titanic() {
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
	csvData := loadAdditionCsvData()

	trainDataSet := csvData[:600]
	testDataSet := csvData[600:]

	//Training data-set
	for _, element := range trainDataSet {
		network.ObserveTrain(element.El1, element.El2)
	}

	//Testing data-set
	for _, element := range testDataSet {
		network.ObserveTest(element.El1, element.El2)
	}
	//Train
	network.Train(0, 0, true, false)
	for network.IsTraining() {
		fmt.Println("Training...")
		time.Sleep(1)
	}

	//Test
	errors := 0
	for _, element := range testDataSet {
		output := network.Classify(element.El1)

		result, _ := strconv.ParseInt(output["result"], 10, 64)
		if _, ok := element.El2[result]; !ok {
			errors += 1
		}
	}

	fmt.Sprintf("%s prediction errors on %s test items", errors, len(testDataSet))
}

func loadTitanicCsvData() []*Pair {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/sample/datasets/titanic.csv")
	if err != nil {
		fmt.Println("Error:", err)
		return []*Pair{}
	}
	defer file.Close()
	reader := csv.NewReader(file)

	reader.Comma = ','
	result := []*Pair{}
	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error:", err)
			return []*Pair{}
		}
		input := make(map[string]int64)
		output := make(map[string]int64)
		//Columns
		//passid,survival,pclass,name,sex,age,sibsp,parch,ticket,fare,cabin,embarked

		//TODO: populate
		fmt.Println(record)
		p := &Pair{
			El1: input,
			El2: output,
		}

		result = append(result, p)
	}

	return result
}
