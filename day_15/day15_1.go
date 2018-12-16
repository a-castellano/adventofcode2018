// Álvaro Castellano Vela 2018/12/16

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

type Player struct {
	Point Point
	Type  byte
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

type Game struct {
	Map          [][]byte
	ElvesAlive   int
	GoblinsAlive int
	PlayedTurns  int
	Players      []Player
	EndGame      bool
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
		line := []byte(scanner.Text())
		game.Map = append(game.Map, line)
		for column, symbol := range line {
			if symbol == 69 || symbol == 71 {
				var player Player

				player.Point.X = currentRow
				player.Point.Y = column
				player.HP = 200
				player.Type = symbol
				if symbol == 69 {
					game.ElvesAlive++
				} else {
					game.GoblinsAlive++
				}
				game.Players = append(game.Players, player)
			}
		}
		currentRow++
	}

	return game
}

func (player *Player) findEnemiesAndAdjacent(gameMap [][]byte) ([]Point, []Point) {
	var targets []Point
	var adjacentPoints []Point

	switch square := gameMap[player.Point.X-1][player.Point.Y]; {
	case square == 46:
		adjacentPoints = append(adjacentPoints, Point{X: player.Point.X - 1, Y: player.Point.Y})
	case square == 69 || square == 71:
		if square != player.Type {
			targets = append(targets, Point{X: player.Point.X - 1, Y: player.Point.Y})
		}
	}

	switch square := gameMap[player.Point.X][player.Point.Y-1]; {
	case square == 46:
		adjacentPoints = append(adjacentPoints, Point{X: player.Point.X, Y: player.Point.Y - 1})
	case square == 69 || square == 71:
		if square != player.Type {
			targets = append(targets, Point{X: player.Point.X, Y: player.Point.Y - 1})
		}
	}

	switch square := gameMap[player.Point.X][player.Point.Y+1]; {
	case square == 46:
		adjacentPoints = append(adjacentPoints, Point{X: player.Point.X, Y: player.Point.Y + 1})
	case square == 69 || square == 71:
		if square != player.Type {
			targets = append(targets, Point{X: player.Point.X, Y: player.Point.Y + 1})
		}
	}

	switch square := gameMap[player.Point.X+1][player.Point.Y]; {
	case square == 46:
		adjacentPoints = append(adjacentPoints, Point{X: player.Point.X + 1, Y: player.Point.Y})
	case square == 69 || square == 71:
		if square != player.Type {
			targets = append(targets, Point{X: player.Point.X + 1, Y: player.Point.Y})
		}
	}

	return adjacentPoints, targets
}

func (player *Player) findNearEnemies(gameMap [][]byte) []Point {

	var nearTargets []Point

	for i, row := range gameMap {
		for j, symbol := range row {
			if (symbol == 69 || symbol == 71) && symbol != player.Type {
				nearTargets = append(nearTargets, Point{X: i, Y: j})
			}
		}
	}
	return nearTargets
}

func (game *Game) play() {

	for game.EndGame == false {

		for _, player := range game.Players {
			if player.HP > 0 {

				fmt.Printf("Player in %d,%d\n", player.Point.X, player.Point.Y)

				adjacentPoints, targets := player.findEnemiesAndAdjacent(game.Map)

				fmt.Printf("\tAdjacent points:\n")
				for _, point := range adjacentPoints {
					fmt.Printf("\t\t%d,%d\n", point.X, point.Y)
				}
				fmt.Printf("\tTarget points:\n")
				for _, point := range targets {
					fmt.Printf("\t\t%d,%d\n", point.X, point.Y)
				}

				var attack bool = len(targets) != 0
				var move bool = len(adjacentPoints) != 0

				if attack == true {
					fmt.Printf("\n\t\t\tATTACK\n")
				} else {
					if move == true {
						fmt.Printf("\n\t\t\tMOVE\n")
						nearTargets := player.findNearEnemies(game.Map)
						var nearPoints []Point
						fmt.Printf("\n\t\t\tNear targets:\n")
						for _, point := range nearTargets {
							var enemy Player
							enemy.Point.X = point.X
							enemy.Point.Y = point.Y
							enemyAdjacentPoints, _ := enemy.findEnemiesAndAdjacent(game.Map)
							nearPoints = append(nearPoints, enemyAdjacentPoints...)
						}
						fmt.Printf("\n\t\t\tNear points:\n")
						for _, point := range nearPoints {
							fmt.Printf("\t\t\t\t%d,%d\n", point.X, point.Y)
						}

					}
				}

			}
		}

		game.EndGame = true
	}

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
	fmt.Printf("_%d_\n", game.Map[1][2])
	game.play()
	fmt.Println(game.EndGame)
	fmt.Printf("_\n")
}
