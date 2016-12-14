package main

import (
	"fmt"

	"github.com/ziyasal/neurolog/neurolog"
)

func main() {
	options := neurolog.Options{
		Name:         "titanic",
		Type:         "classifier",
		Inputs:       []string{"class1", "class2", "class3", "female", "male", "age", "sibsp", "parch", "fare"},
		Outputs:      []string{"dead", "alive"},
		HiddenLayers: []int{15},
		DatasetSize:  1000,
		TestsetSize:  500,
		RedisHost:    "localhost:6379",
	}

	nn := neurolog.New(options)
	info := nn.Info()

	fmt.Println(info)
}
