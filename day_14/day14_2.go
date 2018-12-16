// Álvaro Castellano Vela 2018/12/16

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

func checkFound(input *[]int, inputSize *int, scores *[]int, startIndex int) bool {

	var found bool = true
	for i, _ := range *input {
		if (*input)[*inputSize-1-i] != (*scores)[startIndex-i] {
			found = false
			break
		}
	}

	return found
}

func main() {

	var inputString string
	var input []int
	var inputSize int
	var inputFound bool

	var elf1, elf2 int
	scores := make([]int, 2)

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a sequence to be found.")
	}

	inputString = args[0]
	for _, c := range inputString {
		input = append(input, int(c)-48)
	}
	inputSize = len(input)

	scores[0] = 3
	scores[1] = 7

	elf1 = 0
	elf2 = 1

	for inputFound == false {

		var combine int = scores[elf1] + scores[elf2]
		newRecipes := make([]int, 0)

		var newLen int

		if combine == 0 {
			newRecipes = append(newRecipes, combine)
		} else {
			for combine > 0 {
				newRecipes = append(newRecipes, combine%10)
				combine = combine / 10
			}
			sort.Slice(newRecipes, func(i, j int) bool {
				return true
			})
		}
		scores = append(scores, newRecipes...)
		newLen = len(scores)
		for i, _ := range newRecipes {
			inputFound = checkFound(&input, &inputSize, &scores, newLen-1-i)
			if inputFound == true {
				fmt.Printf("There are %d recipes before sequence.\n", len(scores)-inputSize-i)
				break
			}
		}

		elf1 = (elf1 + scores[elf1] + 1) % newLen
		elf2 = (elf2 + scores[elf2] + 1) % newLen

	}
}
