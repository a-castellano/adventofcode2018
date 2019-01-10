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

func (p Point) distance(t Point) int {
	return abs(p.X-t.X) + abs(p.Y-t.Y) + abs(p.Z-t.Z) + abs(p.T-t.T)
}

type Constelation struct {
	ID     int
	Points []Point
}

type Constelations []Constelation

func processFile(filename string) int {

	var constelationIndex int = 0
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var point Point
		re := regexp.MustCompile("^(-?[[:digit:]]+),(-?[[:digit:]]+),(-?[[:digit:]]+),(-?[[:digit:]]+)$")
		match := re.FindAllStringSubmatch(scanner.Text(), -1)

		point.X, _ = strconv.Atoi(match[0][1])
		point.Y, _ = strconv.Atoi(match[0][2])
		point.Z, _ = strconv.Atoi(match[0][3])
		point.T, _ = strconv.Atoi(match[0][4])

		fmt.Println(point)

	}
	return 0
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
