// Álvaro Castellano Vela 2018/12/27

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Point struct {
	X int
	Y int
}

type Acre struct {
	Point    Point
	Value    rune
	Adjacent []Point
}

type LumberCollection struct {
	Acres  [][]Acre
	Height int
	Width  int
}

func processFile(filename string) LumberCollection {

	var height, width int
	var lumberCollection LumberCollection

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		var acresLine []Acre
		line := scanner.Text()
		width = len(line)

		for y, acreValue := range line {
			acre := Acre{Point: Point{X: height, Y: y}, Value: acreValue}
			acresLine = append(acresLine, acre)
		}
		lumberCollection.Acres = append(lumberCollection.Acres, acresLine)
		height++
	}

	// Calculate adjacent

	lumberCollection.Width = width
	lumberCollection.Height = height
	for x, line := range lumberCollection.Acres {
		for y, acre := range line {
			var xGreater0, xSmallerWidth, yGreater0, ySmallerHeight bool
			xGreater0 = x-1 >= 0
			yGreater0 = y-1 >= 0
			xSmallerWidth = x+1 < lumberCollection.Width
			ySmallerHeight = y+1 < lumberCollection.Height
			if xGreater0 {
				if yGreater0 {
					acre.Adjacent = append(acre.Adjacent, Point{X: x - 1, Y: y - 1})
				}
				acre.Adjacent = append(acre.Adjacent, Point{X: x - 1, Y: y})
				if ySmallerHeight {
					acre.Adjacent = append(acre.Adjacent, Point{X: x - 1, Y: y + 1})
				}
			}
			if yGreater0 {
				acre.Adjacent = append(acre.Adjacent, Point{X: x, Y: y - 1})
			}
			if ySmallerHeight {
				acre.Adjacent = append(acre.Adjacent, Point{X: x, Y: y + 1})
			}
			if xSmallerWidth {
				if yGreater0 {
					acre.Adjacent = append(acre.Adjacent, Point{X: x + 1, Y: y - 1})
				}
				acre.Adjacent = append(acre.Adjacent, Point{X: x + 1, Y: y})
				if ySmallerHeight {
					acre.Adjacent = append(acre.Adjacent, Point{X: x + 1, Y: y + 1})
				}
			}
			lumberCollection.Acres[x][y] = acre
		}
	}

	return lumberCollection
}

func calculateLumber(lumberCollection LumberCollection) LumberCollection {

	var newLumber LumberCollection

	newLumber.Height = lumberCollection.Height
	newLumber.Width = lumberCollection.Width
	newLumber.Acres = make([][]Acre, lumberCollection.Width)
	for x := 0; x < lumberCollection.Width; x++ {
		newLumber.Acres[x] = make([]Acre, lumberCollection.Height)
	}

	for x := 0; x < lumberCollection.Width; x++ {
		for y := 0; y < lumberCollection.Height; y++ {
			var trees int
			var lumberyards int
			var opens int
			newLumber.Acres[x][y].Adjacent = make([]Point, len(lumberCollection.Acres[x][y].Adjacent))
			for i, point := range lumberCollection.Acres[x][y].Adjacent {
				newLumber.Acres[x][y].Adjacent[i] = lumberCollection.Acres[x][y].Adjacent[i]
				if lumberCollection.Acres[point.X][point.Y].Value == '|' {
					trees++
				} else {
					if lumberCollection.Acres[point.X][point.Y].Value == '#' {
						lumberyards++
					} else {
						opens++
					}
				}
			}
			newLumber.Acres[x][y].Value = lumberCollection.Acres[x][y].Value
			if lumberCollection.Acres[x][y].Value == '.' {
				if trees >= 3 {
					newLumber.Acres[x][y].Value = '|'
				}
			}

			if lumberCollection.Acres[x][y].Value == '|' && lumberyards >= 3 {
				if lumberyards >= 3 {
					newLumber.Acres[x][y].Value = '#'
				}
			}
			if lumberCollection.Acres[x][y].Value == '#' {
				if lumberyards >= 1 && trees >= 1 {
					newLumber.Acres[x][y].Value = '#'
				} else {
					newLumber.Acres[x][y].Value = '.'
				}
			}

			newLumber.Acres[x][y].Point = lumberCollection.Acres[x][y].Point
		}
	}

	return newLumber
}

func main() {

	resources := make(map[int]bool)
	var resourcesArray []int

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	lumberCollection := processFile(filename)
	for minute := 1; minute <= 1000; minute++ {
		lumberCollection = calculateLumber(lumberCollection)
	}
	for minute := 1000; minute <= 2000; minute++ {
		var woodedOcres, lumberyards int
		lumberCollection = calculateLumber(lumberCollection)
		for x := 0; x < lumberCollection.Width; x++ {
			for y := 0; y < lumberCollection.Height; y++ {
				if lumberCollection.Acres[x][y].Value == '|' {
					woodedOcres++
				} else {
					if lumberCollection.Acres[x][y].Value == '#' {
						lumberyards++
					}
				}
			}
		}
		resourceValue := woodedOcres * lumberyards
		if _, ok := resources[resourceValue]; !ok {
			resources[resourceValue] = true
			resourcesArray = append(resourcesArray, resourceValue)
		} else {
			break
		}
	}
	offset := (1000000000 - 1001) % len(resourcesArray)
	fmt.Printf("Resource value: %d\n", resourcesArray[offset])
}
