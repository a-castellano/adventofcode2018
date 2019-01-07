// Álvaro Castellano Vela 2019/01/04

package main

import (
	"fmt"
	"github.com/a-castellano/dijkstra"
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
	graph := dijkstra.NewGraph()
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
			graph.AddVertex(id1)
			graph.AddVertex(id2)
			graph.AddArc(id1, id2, 7)
			graph.AddArc(id2, id1, 7)
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
			//			if x < height-1 {
			//				candidateIDs[(x+1)*height*width+y*width] = (*cave)[x+1][y]
			//			}
			if y > 0 {
				candidateIDs[x*height*width+(y-1)*width] = (*cave)[x][y-1]
			}
			//			if y < width-1 {
			//				candidateIDs[x*height*width+(y+1)*width] = (*cave)[x][y+1]
			//			}

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
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
						graph.AddArc(id3, id4, 1)
						graph.AddArc(id4, id3, 1)
					case wet:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					case narrow:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					}
				case wet:
					switch caveSymbolValue[candidateValue] {
					case rocky:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					case wet:
						id1 = caveID + objectValue["climbing gear"]
						id2 = candidateID + objectValue["climbing gear"]
						id3 = caveID + objectValue["neither"]
						id4 = candidateID + objectValue["neither"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
						graph.AddArc(id3, id4, 1)
						graph.AddArc(id4, id3, 1)
					case narrow:
						id1 = caveID + objectValue["neither"]
						id2 = candidateID + objectValue["neither"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					}
				case narrow:
					switch caveSymbolValue[candidateValue] {
					case rocky:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					case wet:
						id1 = caveID + objectValue["neither"]
						id2 = candidateID + objectValue["neither"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
					case narrow:
						id1 = caveID + objectValue["torch"]
						id2 = candidateID + objectValue["torch"]
						id3 = caveID + objectValue["neither"]
						id4 = candidateID + objectValue["neither"]
						graph.AddArc(id1, id2, 1)
						graph.AddArc(id2, id1, 1)
						graph.AddArc(id3, id4, 1)
						graph.AddArc(id4, id3, 1)
					}
				}
			}

			//			var nodeId int = x*height + y
			//			nodeSymbol := (*cave)[x][y]
			//			graph.AddArc(nodeSymbol, nodeSymbol, 7)
			//			switch nodeSymbol {
			//			case caveSymbol[rocky]:
			//			case caveSymbol[wet]:
			//			case caveSymbol[narrow]:
			//			}
		}
	}
	best, err := graph.Shortest(0+objectValue["torch"], targetX*height*width+targetY*width+objectValue["torch"], int64(1000))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Shotest time inverted ", best.Distance)

	//
	//	//Add arcs
	//	graph.AddArc(0, 1, 1)
	//	graph.AddArc(0, 2, 1)
	//	graph.AddArc(1, 0, 1)
	//	graph.AddArc(1, 2, 2)
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

	//	for i := 0; i <= targetX; i++ {
	//		for j := 0; j <= targetY; j++ {
	//			totalRisk += caveSymbolValue[cave[i][j]]
	//		}
	//	}
	//
	//	fmt.Printf("Total Risk: %d\n", totalRisk)
}
