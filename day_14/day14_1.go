// Álvaro Castellano Vela 2018/12/16

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
)

func main() {

	var recipes int
	var elf1, elf2 int
	scores := make([]int, 2)

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a number of recipes.")
	}

	recipes, _ = strconv.Atoi(args[0])

	scores[0] = 3
	scores[1] = 7

	elf1 = 0
	elf2 = 1

	for len(scores) < recipes+10 {

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

		elf1 = (elf1 + scores[elf1] + 1) % newLen
		elf2 = (elf2 + scores[elf2] + 1) % newLen
	}

	fmt.Printf("Scoreboard: ")
	for i := recipes; i < recipes+10; i++ {
		fmt.Printf("%d", scores[i])
	}
	fmt.Printf("\n")
}
