package neurolog

import (
	"github.com/deckarep/golang-set"
)

func _validateInput(network NeuralNetwork, input map[string]int64) {
	_validate(network._opts.Inputs, input, "Input does not have the required keys")
}

func _validateOutput(network NeuralNetwork, output map[string]int64) {
	_validate(network._opts.Outputs, output, "Output does not have the required keys")
}

func _validate(array []string, dataToValidate map[string]int64, message string) {
	keys := make([]interface{}, 0, len(dataToValidate))
	for k := range dataToValidate {
		keys = append(keys, k)
	}

	set1 := mapset.NewSetFromSlice(keys)
	set2 := mapset.NewSet()

	for _, k := range array {
		set2.Add(k)
	}

	if set1.Equal(set2) == false {
		panic(message)
	}
}
