// Álvaro Castellano Vela 2018/12/06

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	X    int
	Y    int
	id   int
	area int
}

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func manhattanDistance(a Point, b Point) int {
	return abs(b.X-a.X) + abs(b.Y-a.Y)
}

func processLine(line string) (int, int) {
	var x, y int

	re := regexp.MustCompile("([[:digit:]]+), ([[:digit:]]+)$")
	match := re.FindAllStringSubmatch(line, -1)

	x, _ = strconv.Atoi(match[0][1])
	y, _ = strconv.Atoi(match[0][2])

	return x, y
}

func processFile(filename string, offset int) []Point {

	var points []Point
	var counter = 0

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var p Point
		p.X, p.Y = processLine(scanner.Text())
		p.X += offset
		p.Y += offset
		p.id = counter
		points = append(points, p)
		counter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return points
}

func countRegion(points *[]Point, limit int) int {

	var count = 0

	for x := 0; x < 4000; x++ {
		for y := 0; y < 4000; y++ {

			var distance = 0
			var currentPoint Point

			currentPoint.X = x
			currentPoint.Y = y

			for _, point := range *points {
				distance += manhattanDistance(currentPoint, point)
				if distance >= limit {
					break
				}
			}
			if distance < limit {
				count++
			}
		}
	}

	return count

}

func main() {

	var points []Point
	var offset = 2000

	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatal("You must supply a file to process and limit.")
	}

	filename := args[0]
	limit, _ := strconv.Atoi(args[1])

	points = processFile(filename, offset)

	fmt.Printf("Size of the region => %d\n", countRegion(&points, limit))
}
