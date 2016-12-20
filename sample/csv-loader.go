package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
)

type Pair struct {
	Input  map[string]int64
	Output map[string]int64
}

func loadAdditionCsvData(sampleFilePath string) []*Pair {
	file, err := os.Open(sampleFilePath)
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
			Input:  input,
			Output: output,
		}

		result = append(result, p)
	}

	return result
}

func loadTitanicCsvData(sampleFilePath string) []*Pair {
	file, err := os.Open(sampleFilePath)
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
			Input:  input,
			Output: output,
		}

		result = append(result, p)
	}

	return result
}
