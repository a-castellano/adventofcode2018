// Álvaro Castellano Vela 2018/12/16

package main

import (
	"bufio"
	"fmt"
	"github.com/beefsack/go-astar"
	"log"
	"os"
	"sort"
	"strings"
)

const (
	// KindPlain (.) is a plain tile with a movement cost of 1.
	KindPlain = iota
	// KindWall (#) is a tile which blocks movement.
	KindWall
	// KindElf (E) is a elf tile which blocks movement.
	KindElf
	// KindGoblin (G) is a goblin tile which blocks movement.
	KindGoblin
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
	KindPlain:  '.',
	KindWall:   '#',
	KindElf:    'E',
	KindGoblin: 'G',
	KindFrom:   'F',
	KindTo:     'T',
	KindPath:   '●',
}

// RuneKinds map input runes to tile kinds.
var RuneKinds = map[rune]int{
	'.': KindPlain,
	'#': KindWall,
	'E': KindElf,
	'G': KindGoblin,
	'X': KindWall,
	'F': KindFrom,
	'T': KindTo,
}

// KindCosts map tile kinds to movement costs.
var KindCosts = map[int]float64{
	KindPlain: 1.0,
	KindFrom:  1.0,
	KindTo:    1.0,
}

// A Tile is a tile in a grid which implements Pather.
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
			n.Kind == KindPlain { //Try checking only if KindPlain
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
		return float64(1000000)
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

// RenderPath renders a path on top of a world.
func (w World) RenderPath(path []astar.Pather) string {
	height := len(w)
	if height == 0 {
		return ""
	}
	width := len(w[0])
	pathLocs := map[string]bool{}
	for _, p := range path {
		pT := p.(*Tile)
		pathLocs[fmt.Sprintf("%d,%d", pT.Point.X, pT.Point.Y)] = true
	}
	rows := make([]string, width)
	for x := 0; x < height; x++ {
		for y := 0; y < width; y++ {
			t := w.Tile(x, y)
			r := ' '
			if pathLocs[fmt.Sprintf("%d,%d", x, y)] {
				r = KindRunes[KindPath]
			} else if t != nil {
				r = KindRunes[t.Kind]
			}
			rows[x] += string(r)
		}
	}
	return strings.Join(rows, "\n")
}

type Point struct {
	X int
	Y int
}

type Player struct {
	Point Point
	Type  rune
	HP    int
}

type Players []Player

func (x Players) Len() int      { return len(x) }
func (x Players) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Players) Less(i, j int) bool {
	var result bool

	if x[i].Point.X < x[j].Point.X {
		result = true
	} else {
		if x[i].Point.X > x[j].Point.X {
			result = false
		} else {
			if x[i].Point.Y < x[j].Point.Y {
				result = true
			} else {
				result = false
			}
		}
	}
	return result
}

type Points []Point

func (x Points) Len() int      { return len(x) }
func (x Points) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Points) Less(i, j int) bool {
	var result bool

	if x[i].X < x[j].X {
		result = true
	} else {
		if x[i].X > x[j].X {
			result = false
		} else {
			if x[i].Y < x[j].Y {
				result = true
			} else {
				result = false
			}
		}
	}
	return result
}

type Weakers []Player

func (x Weakers) Len() int           { return len(x) }
func (x Weakers) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
func (x Weakers) Less(i, j int) bool { return x[i].HP < x[j].HP } // ->>> two players can have teh same HP!!!

type Game struct {
	Map          [][]rune
	ElvesAlive   int
	GoblinsAlive int
	PlayedTurns  int
	Players      []Player
	EndGame      bool
	Rows         int
	Columns      int
	World        World
}

func (game *Game) ParseWorld() {
	game.World = World{}
	for x, row := range game.Map {
		for y, raw := range row {
			kind, ok := RuneKinds[raw]
			if !ok {
				kind = KindWall
			}
			game.World.SetTile(&Tile{
				Kind: kind,
			}, x, y)
		}
	}
}

func generateGame(filename string) Game {

	var game Game
	var currentRow int = 0

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		var cells []rune
		for column, symbol := range line {

			if symbol == 'E' || symbol == 'G' {
				var player Player
				player.HP = 200
				player.Type = symbol
				player.Point.X = currentRow
				player.Point.Y = column

				if symbol == 'E' {
					game.ElvesAlive++
				} else {
					game.GoblinsAlive++
				}
				game.Players = append(game.Players, player)
			}
			cells = append(cells, symbol)
		}
		game.Map = append(game.Map, cells)
		currentRow++
	}

	game.Rows, game.Columns = len(game.Map), len(game.Map[0])
	game.ParseWorld()
	return game
}

func (game *Game) getAdjacent(point Point) []Point {
	// According with this exercise all maps are surrounded by walls so this function never checks an index out of map's range

	var adjacent []Point

	if game.Map[point.X-1][point.Y] == '.' {
		adjacent = append(adjacent, Point{X: point.X - 1, Y: point.Y})
	}
	if game.Map[point.X][point.Y-1] == '.' {
		adjacent = append(adjacent, Point{X: point.X, Y: point.Y - 1})
	}
	if game.Map[point.X][point.Y+1] == '.' {
		adjacent = append(adjacent, Point{X: point.X, Y: point.Y + 1})
	}
	if game.Map[point.X+1][point.Y] == '.' {
		adjacent = append(adjacent, Point{X: point.X + 1, Y: point.Y})
	}

	return adjacent
}

func (game *Game) getBesideTarget(player Player) []Player {

	var targetType rune

	var targets []Player

	switch player.Type {
	case 'E':
		targetType = 'G'
	case 'G':
		targetType = 'E'
	}

	if game.Map[player.Point.X-1][player.Point.Y] == targetType {
		for _, targetPlayer := range game.Players {
			if targetPlayer.Point.X == player.Point.X-1 && targetPlayer.Point.Y == player.Point.Y && targetPlayer.HP > 0 {
				targets = append(targets, targetPlayer)
				break
			}
		}
	}
	if game.Map[player.Point.X][player.Point.Y-1] == targetType {
		for _, targetPlayer := range game.Players {
			if targetPlayer.Point.X == player.Point.X && targetPlayer.Point.Y == player.Point.Y-1 && targetPlayer.HP > 0 {
				targets = append(targets, targetPlayer)
				break
			}
		}
	}
	if game.Map[player.Point.X][player.Point.Y+1] == targetType {
		for _, targetPlayer := range game.Players {
			if targetPlayer.Point.X == player.Point.X && targetPlayer.Point.Y == player.Point.Y+1 && targetPlayer.HP > 0 {
				targets = append(targets, targetPlayer)
				break
			}
		}
	}
	if game.Map[player.Point.X+1][player.Point.Y] == targetType {
		for _, targetPlayer := range game.Players {
			if targetPlayer.Point.X == player.Point.X+1 && targetPlayer.Point.Y == player.Point.Y && targetPlayer.HP > 0 {
				targets = append(targets, targetPlayer)
				break
			}
		}
	}

	sort.Sort(Weakers(targets))

	return targets
}

func (game *Game) findTargetsAdjacentPoints(player Player) []Point {

	var points []Point

	for _, target := range game.Players {
		if target.Type != player.Type && target.HP > 0 {
			points = append(points, game.getAdjacent(target.Point)...)
		}
	}

	return points
}

func (game *Game) play() int {
	fmt.Println("PLAY")
	var rounds int

	for game.EndGame == false {
		//for i := 1; i <= 48; i++ {
		fmt.Printf("ROUND %d\n\n", rounds)
		for _, line := range game.Map {
			fmt.Println(string(line))
		}

		// For each player find its targets
		for playerID, player := range game.Players {
			//var move bool
			var attack bool
			if player.HP > 0 {
				nearTargets := game.getBesideTarget(player)
				if len(nearTargets) > 0 {
					attack = true
					//fmt.Printf("Player in %d,%d Attacks\n", player.Point.X, player.Point.Y)
				}
				if attack == false {
					//fmt.Printf("Player in %d,%d:\n\t", player.Point.X, player.Point.Y)
					//fmt.Println(game.findTargetsAdjacentPoints(player))
					nearPoints := make(map[int][]Point)
					var minDistance = 100000000000000
					var foundNearPoint bool
					for _, point := range game.findTargetsAdjacentPoints(player) {
						path, distance, found := astar.Path(game.World.Tile(player.Point.X, player.Point.Y), game.World.Tile(point.X, point.Y))
						if found {
							foundNearPoint = true
							//fmt.Printf("Distance to target %d,%d: %f\n", point.X, point.Y, distance)
							var intDistance int = int(distance)
							if intDistance <= minDistance {
								nearPoints[intDistance] = append(nearPoints[intDistance], Point{X: path[len(path)-2].(*Tile).Point.X, Y: path[len(path)-2].(*Tile).Point.Y})
								minDistance = intDistance
							}
						}
					}
					if foundNearPoint {
						//fmt.Println(nearPoints[minDistance])
						sort.Sort(Points(nearPoints[minDistance]))
						//fmt.Printf("\t\tNearest point: \n")
						//fmt.Println(nearPoints[minDistance][0])
						game.Map[player.Point.X][player.Point.Y] = 46
						game.Players[playerID].Point.X = nearPoints[minDistance][0].X
						game.Players[playerID].Point.Y = nearPoints[minDistance][0].Y
						player.Point.X = nearPoints[minDistance][0].X
						player.Point.Y = nearPoints[minDistance][0].Y
						game.Map[player.Point.X][player.Point.Y] = player.Type
						// After any change regenerate board
						game.ParseWorld()
						//var newPoint Point = sort.Sort(Points(nearPoints[minDistance]))[0]
						nearTargets = game.getBesideTarget(player)
						if len(nearTargets) > 0 {
							attack = true
							//fmt.Printf("After moving Player in %d,%d Attacks\n", player.Point.X, player.Point.Y)
						}
					}
				}
				if attack == true { //Attack
					//fmt.Println("_____________ATTACK INCOMING_________________")
					//fmt.Printf("I'm %s\n", string(player.Type))
					//fmt.Println(nearTargets)
					var target = nearTargets[0]
					/*

						Otherwise, the adjacent target with the fewest hit points is selected; in a tie, the adjacent target with the fewest hit points which is first in reading order is selected.
					*/
					//fmt.Printf("Target POINT: %d,%d\n", target.Point.X, target.Point.Y)
					var targetPlayerID int
					for targetPlayerIDCandidate, targetPlayerCandidate := range game.Players {
						if targetPlayerCandidate.Point.X == target.Point.X && targetPlayerCandidate.Point.Y == target.Point.Y {
							targetPlayerID = targetPlayerIDCandidate
							//fmt.Printf("%d,%d is going to be attacked, current HP: %d\n", targetPlayerCandidate.Point.X, targetPlayerCandidate.Point.Y, game.Players[targetPlayerID].HP)
							break
						}
					}
					//fmt.Println(target)
					//fmt.Println(targetPlayerID)
					game.Players[targetPlayerID].HP -= 3
					if game.Players[targetPlayerID].HP <= 0 {
						game.Map[game.Players[targetPlayerID].Point.X][game.Players[targetPlayerID].Point.Y] = 46 //.
						game.Players[targetPlayerID].Point.X = -1
						game.Players[targetPlayerID].Point.Y = -1
						game.ParseWorld()
						if game.Players[targetPlayerID].Type == 69 {
							game.ElvesAlive--
						} else {
							game.GoblinsAlive--
						}
						//fmt.Println("DEAD PLAYER")
					}
					//fmt.Println(game.Players[targetPlayerID].HP)
					//fmt.Println("_____________ATTACK END_________________")
				}
			}
			if game.ElvesAlive == 0 || game.GoblinsAlive == 0 {
				game.EndGame = true
				fmt.Println("_____________END GAME_________________")
				break
			}
		}

		sort.Sort(Players(game.Players))
		for _, player := range game.Players {
			fmt.Println(player)
		}
		if game.EndGame == true {
			break
		}
		rounds++
	}
	fmt.Printf("\n\nRounds: %d\n", rounds)
	var totalHP int
	for _, player := range game.Players {
		if player.HP > 0 {
			totalHP += player.HP
		}
	}
	fmt.Printf("TotalHP: %d\n", totalHP)
	return rounds * totalHP
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	game := generateGame(filename)
	fmt.Printf("Elves Alive: %d\n", game.ElvesAlive)
	fmt.Printf("Goblins Alive: %d\n", game.GoblinsAlive)
	fmt.Println(string(game.Map[0]))
	outcome := game.play()
	fmt.Println(game.EndGame)
	fmt.Printf("Outcome: %d\n", outcome)
	//	path, distance, found := astar.Path(p1, p2)
	//	if !found {
	//		log.Println("Could not find path")
	//	}
	//	fmt.Println(world.RenderPath(path))
	//	for _, step := range path {
	//		pT := step.(*Tile)
	//		fmt.Println(pT.Point)
	//	}
	//	fmt.Println(distance)
	//
	//	_, _, found = astar.Path(p1, p3)
	//	if !found {
	//		log.Println("Could not find path")
	//	}
	//
}
