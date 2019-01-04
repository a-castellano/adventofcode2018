// Álvaro Castellano Vela 2019/01/02

package main

import (
	"bufio"
	"fmt"
	"github.com/beefsack/go-astar"
	"log"
	"os"
	"strings"
)

const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota
	// KindWall (#) is a tile which blocks movement.
	KindWall
	// KindVerticalDoor (|).
	KindVerticalDoor
	// KindHorizontalDoor (-).
	KindHorizontalDoor
	// KindFrom (F) is a tile which marks where the path should be calculated
	// from.
	KindFrom
	// KindTo (T) is a tile which marks the goal of the path.
	KindTo
	// KindPath (●) is a tile to represent where the path is in the output.
	KindPath
)

// KindRunes map tile kinds to output runes.
var KindRunes = map[int]rune{
	KindPlain:          '.',
	KindWall:           '#',
	KindVerticalDoor:   '|',
	KindHorizontalDoor: '-',
	KindTo:             'T',
	KindPath:           '●',
}

// RuneKinds map input runes to tile kinds.
var RuneKinds = map[rune]int{
	'.': KindPlain,
	'#': KindWall,
	'|': KindVerticalDoor,
	'-': KindHorizontalDoor,
	'F': KindFrom,
	'T': KindTo,
}

// KindCosts map tile kinds to movement costs.
var KindCosts = map[int]float64{
	KindPlain:          1.0,
	KindVerticalDoor:   1.0,
	KindHorizontalDoor: 1.0,
}

type Tile struct {
	// Kind is the kind of tile, potentially affecting movement.
	Kind int
	// X and Y are the coordinates of the tile.
	Point Point
	// W is a reference to the World that the tile is a part of.
	W World
}

// PathNeighbors returns the neighbors of the tile, excluding blockers and
// tiles off the edge of the board.
func (t *Tile) PathNeighbors() []astar.Pather {
	neighbors := []astar.Pather{}
	for _, offset := range [][]int{
		{-1, 0},
		{0, -1},
		{0, 1},
		{1, 0},
	} {
		if n := t.W.Tile(t.Point.X+offset[0], t.Point.Y+offset[1]); n != nil &&
			n.Kind != KindWall { //Try checking only if KindPlain
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

// PathNeighborCost returns the movement cost of the directly neighboring tile.
func (t *Tile) PathNeighborCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	return KindCosts[toT.Kind]
}

// PathEstimatedCost uses Manhattan distance to estimate orthogonal distance
// between non-adjacent nodes.
func (t *Tile) PathEstimatedCost(to astar.Pather) float64 {
	toT := to.(*Tile)
	if toT == nil {
		return float64(0)
	}
	absX := toT.Point.X - t.Point.X
	if absX < 0 {
		absX = -absX
	}
	absY := toT.Point.Y - t.Point.Y
	if absY < 0 {
		absY = -absY
	}
	return float64(absX + absY)
}

// World is a two dimensional map of Tiles.
type World map[int]map[int]*Tile

// Tile gets the tile at the given coordinates in the world.
func (w World) Tile(x, y int) *Tile {
	if w[x] == nil {
		return nil
	}
	return w[x][y]
}

// SetTile sets a tile at the given coordinates in the world.
func (w World) SetTile(t *Tile, x, y int) {
	if w[x] == nil {
		w[x] = map[int]*Tile{}
	}
	w[x][y] = t
	t.Point.X = x
	t.Point.Y = y
	t.W = w
}

// FirstOfKind gets the first tile on the board of a kind, used to get the from
// and to tiles as there should only be one of each.
func (w World) FirstOfKind(kind int) *Tile {
	for _, row := range w {
		for _, t := range row {
			if t.Kind == kind {
				return t
			}
		}
	}
	return nil
}

// From gets the from tile from the world.
func (w World) From() *Tile {
	return w.FirstOfKind(KindFrom)
}

// To gets the to tile from the world.
func (w World) To() *Tile {
	return w.FirstOfKind(KindTo)
}

type Point struct {
	X int
	Y int
}

func fillMaze(regex string, maze *[500][500]rune, startPoint Point) {

	var point Point = startPoint
	var i int = 0
	var bracket bool = false

	for regex[i] != '$' && bracket == false {
		switch regex[i] {
		case 'N':
			point.X--
			(*maze)[point.X][point.Y] = '-'
			point.X--
			(*maze)[point.X][point.Y] = '.'
		case 'S':
			point.X++
			(*maze)[point.X][point.Y] = '-'
			point.X++
			(*maze)[point.X][point.Y] = '.'
		case 'W':
			point.Y--
			(*maze)[point.X][point.Y] = '|'
			point.Y--
			(*maze)[point.X][point.Y] = '.'
		case 'E':
			point.Y++
			(*maze)[point.X][point.Y] = '|'
			point.Y++
			(*maze)[point.X][point.Y] = '.'
		case '(':
			var openBrackets int = 1
			var startOfSustring int = i + 1
			var endOfSustring int = i + 1
			var pipes []int
			for openBrackets != 0 {
				switch regex[endOfSustring] {
				case '(':
					openBrackets++
				case ')':
					openBrackets--
				case '|':
					if openBrackets == 1 {
						pipes = append(pipes, endOfSustring)
					}
				}
				endOfSustring++
			}
			for _, offset := range pipes {
				var subString strings.Builder
				subString.WriteString(regex[startOfSustring:offset])
				subString.WriteString(regex[endOfSustring:])

				fillMaze(subString.String(), maze, point)
				startOfSustring = offset + 1
			}
			var subString strings.Builder
			subString.WriteString(regex[startOfSustring : endOfSustring-1])
			subString.WriteString(regex[endOfSustring:])
			fillMaze(subString.String(), maze, point)
			bracket = true
		}
		i++
	}
}

func createMap(filename string) int {

	var maze [500][500]rune
	world := make(World)

	var pathsWithThousandDoors int = 0

	startPoint := Point{X: 249, Y: 249}

	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			maze[i][j] = '#'
		}
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	regex := scanner.Text()

	regex = strings.Replace(regex, "|)", ")", -1)

	maze[startPoint.X][startPoint.Y] = '.'
	fillMaze(regex[1:], &maze, startPoint)

	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			kind, ok := RuneKinds[maze[i][j]]
			if !ok {
				kind = KindWall
			}
			world.SetTile(&Tile{Kind: kind}, i, j)
		}
	}

	for i := 0; i < 500; i++ {
		for j := 0; j < 500; j++ {
			if maze[i][j] == '.' {
				path, _, found := astar.Path(world.Tile(startPoint.X, startPoint.Y), world.Tile(i, j))
				if found {
					var doors int
					for _, point := range path {
						content := maze[point.(*Tile).Point.X][point.(*Tile).Point.Y]
						if content == '|' || content == '-' {
							doors++
						}
					}
					if doors >= 1000 {
						pathsWithThousandDoors++
					}
				}
			}
		}
	}

	return pathsWithThousandDoors
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	paths := createMap(filename)
	fmt.Printf("Paths with at least thousand Doors: %d\n", paths)
}
