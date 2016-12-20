package main

import (
	"os"
)

func main() {
	pwd, _ := os.Getwd()

	addition(pwd + "/sample/datasets/addition.csv")

	//titanic(pwd + pwd + "/sample/datasets/titanic.csv")
}
