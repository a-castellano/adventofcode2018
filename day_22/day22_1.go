// Álvaro Castellano Vela 2019/01/04

package main

import (
	"fmt"
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

func main() {

	var depth, targetX, targetY int

	var totalRisk int

	var cave [][]rune

	args := os.Args[1:]
	if len(args) != 3 {
		log.Fatal("You must supply:\n\tCave depth.\n\tTarget's X coordinate.\n\tTarget's Y coordinate.\n")
	}

	depth, _ = strconv.Atoi(args[0])
	targetX, _ = strconv.Atoi(args[1])
	targetY, _ = strconv.Atoi(args[2])

	for i := 0; i < targetX+1; i++ {
		row := make([]rune, targetY+1)
		cave = append(cave, row)
	}

	fillCave(&cave, depth, targetX+1, targetY+1)

	cave[targetX][targetY] = caveSymbol[(depth%modulo)%3]

	for i := 0; i <= targetX; i++ {
		for j := 0; j <= targetY; j++ {
			totalRisk += caveSymbolValue[cave[i][j]]
		}
	}

	fmt.Printf("Total Risk: %d\n", totalRisk)
}
