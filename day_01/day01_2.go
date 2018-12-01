// Álvaro Castellano Vela 2018/12/01

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func processFile(filename string, frequency *int, seenMap map[int]bool) (bool, int) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		operation, _ := strconv.Atoi(scanner.Text())
		*frequency += operation
		if _, ok := seenMap[*frequency]; !ok {
			seenMap[*frequency] = true
		} else {
			return true, *frequency
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return false, 0
}

func main() {
	var frequency int
	var seen int
	var finished bool = false
	seenMap := make(map[int]bool)
	seenMap[0] = true

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	for finished == false {
		finished, seen = processFile(filename, &frequency, seenMap)
	}
	fmt.Printf("Repeated frequency: %d\n", seen)
}
