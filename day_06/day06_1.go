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

func processFile(filename string, offset int, grid *[4000][4000]int) []Point {

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
		grid[p.X][p.Y] = p.id
		points = append(points, p)
		counter++
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return points
}

func getClosest(x int, y int, points *[]Point) int {
	var minDistance = 10000000
	var point Point
	var closestPoint int
	distances := make(map[int]int)

	point.X = x
	point.Y = y

	for _, RegisteredPoint := range *points {
		var distance = manhattanDistance(point, RegisteredPoint)
		distances[distance] += 1
		if minDistance > distance {
			minDistance = distance
			closestPoint = RegisteredPoint.id
		}
	}
	if distances[minDistance] > 1 {
		closestPoint = -1
	}

	return closestPoint
}

func fillGrid(grid *[4000][4000]int, points *[]Point) {

	for x := 0; x < 4000; x++ {
		for y := 0; y < 4000; y++ {
			var closest = getClosest(x, y, points)
			grid[x][y] = closest
			if closest != -1 {
				(*points)[closest].area++
			}
		}
	}

}

func main() {

	var points []Point
	var grid [4000][4000]int
	var offset = 2000
	var maxArea, limit = 0, 20000

	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	points = processFile(filename, offset, &grid)

	for _, point := range points {
		if point.area < limit && point.area > maxArea {
			maxArea = point.area
		}
	}

	fmt.Printf("Max area => %d\n", maxArea)
}
