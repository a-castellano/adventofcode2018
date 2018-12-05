// Álvaro Castellano Vela 2018/12/05

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func readLineFromFile(filename string) string {

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	scanner.Scan()
	return scanner.Text()
}

func react(polymer string) (string, int) {
	var removed int = 0
	var indexToRemove int = -1
	var i int
	var polymerLenght = len(polymer)
	var newPolymer string

	var reacted []string

	for i < polymerLenght-1 {
		if abs(int(polymer[i])-int(polymer[i+1])) == 32 {
			indexToRemove = i
			removed = 1
			break
		}
		i++
	}
	if removed == 1 {
		reacted = append(reacted, polymer[:indexToRemove])
		reacted = append(reacted, polymer[indexToRemove+2:])
		newPolymer = strings.Join(reacted, "")

	} else {
		newPolymer = polymer
	}
	return newPolymer, removed
}

func removeByUnit(polymer string, unit int) (string, int) {
	var remove, removed int = 1, 0
	var indexToRemove int = -1
	var i int
	var newPolymer string = polymer

	var reacted []string

	for remove == 1 {
		var polymerLenght = len(newPolymer)
		remove = 0
		for i = 0; i < polymerLenght; i++ {
			if int(newPolymer[i]) == unit || int(newPolymer[i]) == unit+32 {
				indexToRemove = i
				removed = 1
				remove = 1
				break
			}
		}
		if remove == 1 {
			reacted = append(reacted, newPolymer[:indexToRemove])
			reacted = append(reacted, newPolymer[indexToRemove+1:])
			newPolymer = strings.Join(reacted, "")
			reacted = reacted[:0]
		}
	}

	return newPolymer, removed
}

func main() {

	var polymer string
	var candidate string
	var unit int = 65
	var minLen = 10000000000000000
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	polymer = readLineFromFile(filename)

	for unit < 91 {
		polymerCandidate, removed := removeByUnit(polymer, unit)
		if removed == 1 {
			var reacted = 1
			candidate = polymerCandidate
			for reacted == 1 {
				candidate, reacted = react(candidate)
			}
			var candidateLenght = len(candidate)
			if candidateLenght < minLen {
				minLen = candidateLenght
			}
		}
		unit++
	}

	fmt.Printf("Final Polymer len: %d\n", minLen)
}
