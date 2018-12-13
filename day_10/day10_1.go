// Álvaro Castellano Vela 2018/12/13

package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	X  int
	Y  int
	VX int
	VY int
}

type Sky struct {
	Points []Point
	MaxX   int
	MaxY   int
	MinX   int
	MinY   int
}

func processLine(line string) Point {
	var point Point

	re := regexp.MustCompile("position=<[ ]*(-?[[:digit:]]+),[ ]*(-?[[:digit:]]+)> velocity=<[ ]*(-?[[:digit:]]+),[ ]*(-?[[:digit:]]+)>")
	match := re.FindAllStringSubmatch(line, -1)

	point.X, _ = strconv.Atoi(match[0][1])
	point.Y, _ = strconv.Atoi(match[0][2])
	point.VX, _ = strconv.Atoi(match[0][3])
	point.VY, _ = strconv.Atoi(match[0][4])

	return point
}

func (sky *Sky) RecalculateLimits() {

	(*sky).MaxX = -1000000
	(*sky).MinX = 1000000
	(*sky).MinY = 1000000
	(*sky).MaxY = -1000000

	for _, point := range (*sky).Points {
		if point.X > (*sky).MaxX {
			(*sky).MaxX = point.X
		}
		if point.Y > (*sky).MaxY {
			(*sky).MaxY = point.Y
		}
		if point.X < (*sky).MinX {
			(*sky).MinX = point.X
		}
		if point.Y < (*sky).MinY {
			(*sky).MinY = point.Y
		}
	}
}

func (sky *Sky) CreateImage(iterationNumber int) {

	var columns int = (abs((*sky).MinX) + abs((*sky).MaxX)) * 2
	var rows int = (abs((*sky).MinY) + abs((*sky).MaxY)) * 2
	var imageName string

	myImage := image.NewRGBA(image.Rect(0, 0, columns, rows))

	for _, point := range (*sky).Points {
		myImage.Set(point.X+columns/2, point.Y+rows/2, color.RGBA{255, 0, 0, 255})
	}

	imageName = fmt.Sprintf("images/image_%d.png", iterationNumber)
	outputFile, _ := os.Create(imageName)
	png.Encode(outputFile, myImage)
	outputFile.Close()

}

func (sky *Sky) maybeItsReadeable() bool {
	rows := make(map[int]int)
	columns := make(map[int]int)

	var columnsNumber, columnsWithNumberPoints int = 5, 0
	var rowsNumber, rowsWithNumberPoints int = 6, 0

	for _, point := range (*sky).Points {
		rows[point.Y]++
		columns[point.X]++
	}

	for _, numberOfRows := range rows {
		if numberOfRows > rowsNumber {
			rowsWithNumberPoints++
		}
	}
	for _, numberOfColumns := range columns {
		if numberOfColumns > columnsNumber {
			columnsWithNumberPoints++
		}
	}
	return (columnsWithNumberPoints + rowsWithNumberPoints) > 15
}

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func processFile(filename string) Sky {

	var sky Sky

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {

		var point Point = processLine(scanner.Text())
		sky.Points = append(sky.Points, point)
	}
	sky.RecalculateLimits()
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return sky
}

func processSky(sky *Sky, iterations int) {
	for i := 0; i < iterations; i++ {
		for index, _ := range (*sky).Points {
			(*sky).Points[index].X += (*sky).Points[index].VX
			(*sky).Points[index].Y += (*sky).Points[index].VY
		}
		sky.RecalculateLimits()
		if sky.maybeItsReadeable() {

			fmt.Printf("This seems to be a text.\n")
			sky.CreateImage(i)
		}
	}
}

func main() {

	var sky Sky

	args := os.Args[1:]
	if len(args) != 2 {
		log.Fatal("You must supply a file to process and max iterations.")
	}

	filename := args[0]
	iterations, _ := strconv.Atoi(args[1])
	sky = processFile(filename)
	processSky(&sky, iterations)
	fmt.Printf("_\n")
}
