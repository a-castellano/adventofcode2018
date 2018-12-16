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
	Square *Square
	Type   byte
	HP     int
}

type Players []Player

func (x Players) Len() int      { return len(x) }
func (x Players) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Players) Less(i, j int) bool {
	var result bool

	if x[i].Square.Point.X < x[j].Square.Point.X {
		result = true
	} else {
		if x[i].Square.Point.X > x[j].Square.Point.X {
			result = false
		} else {
			if x[i].Square.Point.Y < x[j].Square.Point.Y {
				result = true
			} else {
				result = false
			}
		}
	}
	return result
}

type Square struct {
	ID     int
	Player *Player
	Type   byte
	Point  Point
}

type Game struct {
	Map          [][]Square
	ElvesAlive   int
	GoblinsAlive int
	PlayedTurns  int
	Players      []Player
	EndGame      bool
	Rows         int
	Columns      int
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
		var squares []Square
		for column, symbol := range line {

			var square Square

			square.Point.X = currentRow
			square.Point.Y = column
			square.Type = symbol

			if symbol == 69 || symbol == 71 {
				var player Player
				player.HP = 200
				player.Type = symbol
				player.Square = &square
				square.Player = &player

				if symbol == 69 {
					game.ElvesAlive++
				} else {
					game.GoblinsAlive++
				}
				game.Players = append(game.Players, player)
			}
			squares = append(squares, square)
		}
		game.Map = append(game.Map, squares)
		currentRow++
	}

	game.Rows, game.Columns = len(game.Map), len(game.Map[0])

	return game
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
	fmt.Println(game.Map[1][2])
	//game.play()
	fmt.Println(game.EndGame)
	fmt.Printf("_\n")
}
