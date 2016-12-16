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

type Pair struct {
	El1 map[string]int64
	El2 map[string]int64
}

func main() {

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
	csvData := loadCsvData()

	trainDataSet := csvData[:90]
	testDataSet := csvData[90:]

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
		output := network.Run(element.El1)
		fmt.Sprintf("Neural Network calculation %s+%s = %s",
			element.El1["number1"],
			element.El1["number2"],
			output["result"])
		result, _ := strconv.ParseInt(output["result"], 10, 64)
		if result != element.El2["result"] {
			errors += 1
		}
	}

	fmt.Sprintf("%s prediction errors on %s test items", errors, len(testDataSet))
}

func loadCsvData() []*Pair {
	pwd, _ := os.Getwd()
	file, err := os.Open(pwd + "/sample/datasets/addition.csv")
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

		input["number1"], err = strconv.ParseInt(record[0], 10, 64)
		input["number2"], err = strconv.ParseInt(record[1], 10, 64)
		output["result"], err = strconv.ParseInt(record[2], 10, 64)

		p := &Pair{
			El1: input,
			El2: output,
		}
		println(input)
		println(output)
		println(result)
		result = append(result, p)
	}

	return result
}
