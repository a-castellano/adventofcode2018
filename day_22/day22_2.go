// Álvaro Castellano Vela 2019/01/04

package main

import (
	"fmt"
	"github.com/yourbasic/graph"
	"log"
	"os"
	"strconv"
)

const (
	rocky = iota
	wet
	narrow
	modulo         = 20183
	geologicIndexX = 16807
	geologicIndexY = 48271
)

var caveSymbol = map[int]rune{
	rocky:  '.',
	wet:    '=',
	narrow: '|',
}

var caveSymbolValue = map[rune]int{
	'.': rocky,
	'=': wet,
	'|': narrow,
}

var objectValue = map[string]int{
	"neither":       0,
	"climbing gear": 1,
	"torch":         2,
}

func fillCave(cave *[][]rune, depth int, height int, width int) {

	var erosionLevels [][]int

	for i := 0; i < height; i++ {
		row := make([]int, width)
		erosionLevels = append(erosionLevels, row)
	}

	for x := 0; x < height; x++ {
		var geologicIndex int = x * geologicIndexX
		var erosionLevel = (geologicIndex + depth) % modulo
		erosionLevels[x][0] = erosionLevel
		(*cave)[x][0] = caveSymbol[erosionLevel%3]
	}

	for y := 0; y < width; y++ {
		var geologicIndex int = y * geologicIndexY
		var erosionLevel = (geologicIndex + depth) % modulo
		erosionLevels[0][y] = erosionLevel
		(*cave)[0][y] = caveSymbol[erosionLevel%3]
	}

	for x := 1; x < height; x++ {
		for y := 1; y < width; y++ {
			var geologicIndex int = erosionLevels[x-1][y] * erosionLevels[x][y-1]
			var erosionLevel = (geologicIndex + depth) % modulo
			erosionLevels[x][y] = erosionLevel
			(*cave)[x][y] = caveSymbol[erosionLevel%3]
		}
	}
}

func calculatePath(cave *[][]rune, height int, width int, targetX int, targetY int) {

	caveGraph := graph.New((height-1)*height*width + (width-1)*width + 3)
	//Add nodes and transitions inside the same node
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			var id1, id2 int
			switch caveSymbolValue[(*cave)[x][y]] {
			case rocky:
				/*Allowed objects: climbing gear and torch*/
				id1 = x*height*width + y*width + objectValue["climbing gear"]
				id2 = x*height*width + y*width + objectValue["torch"]
			case wet:
				/*Allowed objects: climbing gear and neither*/
				id1 = x*height*width + y*width + objectValue["climbing gear"]
				id2 = x*height*width + y*width + objectValue["neither"]

			case narrow:
				/*Allowed objects: torch and neither*/
				id1 = x*height*width + y*width + objectValue["neither"]
				id2 = x*height*width + y*width + objectValue["torch"]
			}
			caveGraph.AddBothCost(id1, id2, 7)
			caveGraph.AddBothCost(id2, id1, 7)
		}
	}

	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {

			candidateIDs := make(map[int]rune)

			//calculate adjacent
			caveSymbol := (*cave)[x][y]
			caveID := x*height*width + y*width
			if x > 0 {
				candidateIDs[(x-1)*height*width+y*width] = (*cave)[x-1][y]
			}
			if y > 0 {
				candidateIDs[x*height*width+(y-1)*width] = (*cave)[x][y-1]
			}

			for candidateID, candidateValue := range candidateIDs {
				var id1, id2, id3, id4 int
				switch caveSymbolValue[caveSymbol] {
				case rocky:
					switch caveSymbolValue[candidateValue] {
					case rocky:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						id3 = caveID + objectValue["torch"]
						id4 = candidateID + objectValue["torch"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
						caveGraph.AddBothCost(id3, id4, 1)
						caveGraph.AddBothCost(id4, id3, 1)
					case wet:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					case narrow:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					}
				case wet:
					switch caveSymbolValue[candidateValue] {
					case rocky:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					case wet:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						id3 = caveID + objectValue["neither"]
						id4 = candidateID + objectValue["neither"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
						caveGraph.AddBothCost(id3, id4, 1)
						caveGraph.AddBothCost(id4, id3, 1)
					case narrow:
						id1 = caveID + objectValue["neither"]
						id2 = candidateID + objectValue["neither"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					}
				case narrow:
					switch caveSymbolValue[candidateValue] {
					case rocky:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					case wet:
						id1 = caveID + objectValue["neither"]
						id2 = candidateID + objectValue["neither"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
					case narrow:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						id3 = caveID + objectValue["neither"]
						id4 = candidateID + objectValue["neither"]
						caveGraph.AddBothCost(id1, id2, 1)
						caveGraph.AddBothCost(id2, id1, 1)
						caveGraph.AddBothCost(id3, id4, 1)
						caveGraph.AddBothCost(id4, id3, 1)
					}
				}
			}

		}
	}
	_, dist := graph.ShortestPath(caveGraph, 0+objectValue["torch"], targetX*height*width+targetY*width+objectValue["torch"])
	fmt.Println("Shotest time inverted ", dist)
}

func main() {

	var depth, targetX, targetY, caveDimension int

	var cave [][]rune

	args := os.Args[1:]
	if len(args) != 4 {
		log.Fatal("You must supply:\n\tCave depth.\n\tTarget's X coordinate.\n\tTarget's Y coordinate.\n\tCave dimension.\n")
	}

	depth, _ = strconv.Atoi(args[0])
	targetX, _ = strconv.Atoi(args[1])
	targetY, _ = strconv.Atoi(args[2])
	caveDimension, _ = strconv.Atoi(args[3])

	if targetX >= caveDimension || targetY >= caveDimension {
		log.Fatal("Cave dimmension is too low.")
	}

	for i := 0; i < caveDimension; i++ {
		row := make([]rune, caveDimension)
		cave = append(cave, row)
	}

	fillCave(&cave, depth, caveDimension, caveDimension)

	cave[targetX][targetY] = caveSymbol[(depth%modulo)%3]

	calculatePath(&cave, caveDimension, caveDimension, targetX, targetY)

}
