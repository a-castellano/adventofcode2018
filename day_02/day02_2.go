// Álvaro Castellano Vela 2018/12/02

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type CloserIndex struct {
	Target       int
	Closer       int
	EqualLetters int
}

func calculeLetters(target string, candidate string) int {
	var equalLetters int
	for index, _ := range target {
		if target[index] == candidate[index] {
			equalLetters++
		}
	}
	return equalLetters
}

func findCloser(boxes *[]string, boxIndex int) CloserIndex {
	var closerIndex int
	var equalLetters int
	closerIndex = -1
	equalLetters = -1

	//Ugly
	for index, _ := range (*boxes)[:boxIndex] {
		//fmt.Printf("\t\t - %d\n", index)
		var calcualedLetters = calculeLetters((*boxes)[boxIndex], (*boxes)[index])
		if calcualedLetters > equalLetters {
			closerIndex = index
			equalLetters = calcualedLetters
		}
	}
	for index, _ := range (*boxes)[boxIndex+1:] {
		//fmt.Printf("\t\t - %d\n", index+boxIndex+1)
		var calcualedLetters = calculeLetters((*boxes)[boxIndex], (*boxes)[index+boxIndex+1])
		if calcualedLetters > equalLetters {
			closerIndex = index + boxIndex + 1
			equalLetters = calcualedLetters
		}
	}
	return CloserIndex{boxIndex, closerIndex, equalLetters}
}

func findClosers(boxes *[]string) CloserIndex {
	var candidate CloserIndex
	var closest CloserIndex = CloserIndex{-1, -1, -1}
	for index, _ := range *boxes {
		candidate = findCloser(boxes, index)
		if candidate.EqualLetters > closest.EqualLetters {
			closest = candidate
		}
	}

	return closest
}

func commonLetters(boxes *[]string, closerIndex CloserIndex) {
	target, closer := (*boxes)[closerIndex.Target], (*boxes)[closerIndex.Closer]
	for index, _ := range target {
		if target[index] == closer[index] {
			fmt.Printf("%c", target[index])
		}
	}
	fmt.Printf("\n")
}

func processFile(filename string, boxes *[]string) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		*boxes = append(*boxes, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var boxes []string
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]
	processFile(filename, &boxes)

	commonLetters(&boxes, findClosers(&boxes))
}
