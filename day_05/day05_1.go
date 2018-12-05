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

func main() {

	var polymer string
	var reacted int = 1

	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	polymer = readLineFromFile(filename)

	for reacted == 1 {
		polymer, reacted = react(polymer)
	}
	fmt.Printf("Final Polymer len: %d\n", len(polymer))
}
