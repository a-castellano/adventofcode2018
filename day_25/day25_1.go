// Álvaro Castellano Vela 2019/01/07

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type Point struct {
	X int
	Y int
	Z int
	T int
}

func (p Point) Distance(t Point) int {
	return abs(p.X-t.X) + abs(p.Y-t.Y) + abs(p.Z-t.Z) + abs(p.T-t.T)
}

type Constelation struct {
	ID     int
	Points []Point
}

type Constelations []Constelation

func processFile(filename string) int {

	var constelations Constelations
	var constelationIndex int = 0
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		constelationsFound := make([]int, 0)
		var point Point
		re := regexp.MustCompile("^(-?[[:digit:]]+),(-?[[:digit:]]+),(-?[[:digit:]]+),(-?[[:digit:]]+)$")
		match := re.FindAllStringSubmatch(scanner.Text(), -1)

		point.X, _ = strconv.Atoi(match[0][1])
		point.Y, _ = strconv.Atoi(match[0][2])
		point.Z, _ = strconv.Atoi(match[0][3])
		point.T, _ = strconv.Atoi(match[0][4])

		//		fmt.Println(point)

		for pos, constelation := range constelations {
			for _, contalationPoint := range constelation.Points {
				if contalationPoint.Distance(point) <= 3 {
					constelationsFound = append(constelationsFound, pos)
					break
				}
			}
		}

		var constelation Constelation
		constelation.ID = constelationIndex
		constelation.Points = append(constelation.Points, point)

		if len(constelationsFound) > 0 {
			// Merge constelations
			for i := len(constelationsFound) - 1; i >= 0; i-- {
				constelationToMerge := constelationsFound[i]
				constelation.Points = append(constelation.Points, constelations[constelationToMerge].Points...)
				constelations = append(constelations[:constelationToMerge], constelations[constelationToMerge+1:]...)
			}
		}

		constelations = append(constelations, constelation)
		constelationIndex++
	}

	return len(constelations)
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]

	constelations := processFile(filename)

	fmt.Printf("Constelations %d\n", constelations)
}
