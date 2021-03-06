// Álvaro Castellano Vela 2018/12/24

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
	//KindFrom:  1.0,
	//KindTo:    1.0,
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
			n.Kind != KindWall && n.Kind != KindElf && n.Kind != KindGoblin { //Try checking only if KindPlain
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

type TargetNearPoint struct {
	X           int
	Y           int
	TargetPoint Point
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

type TargetNearPoints []TargetNearPoint

func (x TargetNearPoints) Len() int      { return len(x) }
func (x TargetNearPoints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x TargetNearPoints) Less(i, j int) bool {
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

type SameDistanceNearPoints []TargetNearPoint

func (x SameDistanceNearPoints) Len() int      { return len(x) }
func (x SameDistanceNearPoints) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x SameDistanceNearPoints) Less(i, j int) bool {
	var result bool

	if x[i].TargetPoint.X < x[j].TargetPoint.X {
		result = true
	} else {
		if x[i].TargetPoint.X > x[j].TargetPoint.X {
			result = false
		} else {
			if x[i].TargetPoint.Y < x[j].TargetPoint.Y {
				result = true
			} else {
				if x[i].TargetPoint.Y > x[j].TargetPoint.Y {
					result = false
				} else { //Same Point
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
				}
			}
		}
	}

	return result
}

type Weakers []Player

func (x Weakers) Len() int      { return len(x) }
func (x Weakers) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Weakers) Less(i, j int) bool {
	var result bool
	if x[i].HP == x[j].HP {
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

	} else {
		result = x[i].HP < x[j].HP
	}
	return result
}

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

func (game *Game) findTargetsAdjacentPoints(player Player) []TargetNearPoint {

	var points []TargetNearPoint

	for _, target := range game.Players {
		if target.Type != player.Type && target.HP > 0 {
			for _, point := range game.getAdjacent(target.Point) {
				points = append(points, TargetNearPoint{X: point.X, Y: point.Y, TargetPoint: target.Point})
			}
		}
	}

	sort.Sort(TargetNearPoints(points))
	return points
}

func (game *Game) play(elfAttackPower int) (int, bool) {
	var rounds int
	var initialElves int = game.ElvesAlive

	for game.EndGame == false {
		// For each player find its targets
		for playerID, player := range game.Players {
			var attack bool
			if player.HP > 0 {
				nearTargets := game.getBesideTarget(player)
				if len(nearTargets) > 0 {
					attack = true
				}
				if attack == false {
					nearPoints := make(map[int][]TargetNearPoint)
					var minDistance = 100000000000000
					var foundNearPoint bool
					for _, point := range game.findTargetsAdjacentPoints(player) {
						for _, offset := range [][]int{
							{-1, 0},
							{0, -1},
							{0, 1},
							{1, 0},
						} {
							newPoint := Point{X: player.Point.X + offset[0], Y: player.Point.Y + offset[1]}
							if game.Map[newPoint.X][newPoint.Y] == '.' {
								path, distance, found := astar.Path(game.World.Tile(newPoint.X, newPoint.Y), game.World.Tile(point.X, point.Y))
								if found {
									foundNearPoint = true
									var intDistance int = int(distance)
									if intDistance <= minDistance {
										nearPoints[intDistance] = append(nearPoints[intDistance], TargetNearPoint{X: path[len(path)-1].(*Tile).Point.X, Y: path[len(path)-1].(*Tile).Point.Y, TargetPoint: point.TargetPoint})
										minDistance = intDistance
									}
								}
							}
						}

					}
					if foundNearPoint {
						sort.Sort(SameDistanceNearPoints(nearPoints[minDistance]))
						game.Map[player.Point.X][player.Point.Y] = 46
						game.Players[playerID].Point.X = nearPoints[minDistance][0].X
						game.Players[playerID].Point.Y = nearPoints[minDistance][0].Y
						player.Point.X = nearPoints[minDistance][0].X
						player.Point.Y = nearPoints[minDistance][0].Y
						game.Map[player.Point.X][player.Point.Y] = player.Type
						game.ParseWorld()
						nearTargets = game.getBesideTarget(player)
						if len(nearTargets) > 0 {
							attack = true
						}
					}
				}
				if attack == true { //Attack
					var target = nearTargets[0]
					var targetPlayerID int
					for targetPlayerIDCandidate, targetPlayerCandidate := range game.Players {
						if targetPlayerCandidate.Point.X == target.Point.X && targetPlayerCandidate.Point.Y == target.Point.Y {
							targetPlayerID = targetPlayerIDCandidate
							break
						}
					}
					if player.Type == 'E' {
						game.Players[targetPlayerID].HP -= elfAttackPower
					} else {
						game.Players[targetPlayerID].HP -= 3
					}
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
					}
				}
			}
			if game.ElvesAlive == 0 || game.GoblinsAlive == 0 {
				game.EndGame = true
				if playerID == len(game.Players)-1 {
					rounds++
				}
				break
			}
		}

		sort.Sort(Players(game.Players))
		if game.EndGame == true {
			break
		}

		rounds++
	}
	var totalHP int
	for _, player := range game.Players {
		if player.HP > 0 {
			totalHP += player.HP
		}
	}
	return rounds * totalHP, initialElves == game.ElvesAlive
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	game := generateGame(filename)
	var outcome int
	var allElvesAlive bool = false
	var elfAttackPower int = 3
	for allElvesAlive == false {
		outcome, allElvesAlive = game.play(elfAttackPower)
		game = generateGame(filename)
		elfAttackPower++
	}
	fmt.Printf("Outcome: %d\n", outcome)
}
