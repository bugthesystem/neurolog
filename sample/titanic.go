package main

import (
	"fmt"
	"github.com/ziyasal/neurolog/neurolog"
)

func titanic() {
	options := neurolog.Options{
		Name:            "additions",
		Type:            "regressor",
		Inputs:          []string{"number1", "number2"},
		Outputs:         []string{"result"},
		HiddenLayers:    []int{3},
		DatasetSize:     50,
		TestDatasetSize: 10,
		RedisHost:       "localhost:6379",
	}

	nn := neurolog.New(options)
	info := nn.Info()
	fmt.Println(info)

	//TODO:

}
