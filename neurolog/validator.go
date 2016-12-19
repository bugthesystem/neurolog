package neurolog

import (
	"fmt"
	"github.com/deckarep/golang-set"
)

func _validateInput(network NeuralNetwork, input map[string]int64) {
	keys := make([]interface{}, 0, len(input))
	for k := range input {
		keys = append(keys, k)
	}

	set1 := mapset.NewSetFromSlice(keys)
	set2 := mapset.NewSet()
	for k := range network._opts.Inputs {
		set2.Add(k)
	}

	fmt.Println(set1)
	fmt.Println(set2)

	//TODO: fix set equality
	//if !set1.Equal(set2) {
	//	panic("Input does not have the required keys")
	//}
}

func _validateOutput(network NeuralNetwork, output map[string]int64) {
	keys := make([]interface{}, 0, len(output))
	for k := range output {
		keys = append(keys, k)
	}

	set1 := mapset.NewSetFromSlice(keys)
	set2 := mapset.NewSet()

	for k := range network._opts.Outputs {
		set2.Add(k)
	}

	fmt.Println(set1)
	fmt.Println(set2)

	//TODO: fix set equality
	//if !set1.Equal(set2) {
	//	panic("Output does not have the required keys")
	//}
}
