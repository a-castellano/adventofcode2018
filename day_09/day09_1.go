// Álvaro Castellano Vela 2018/12/12

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Players struct {
	Points          map[int]int
	Current         int
	NumberOfPlayers int
}

type Marble struct {
	Value int
}

type MarbleCircle struct {
	Marbles         []Marble
	Size            int
	CurrentPosition int
}

func (circle *MarbleCircle) CalculateNextPosition(increment int) int {
	return (circle.CurrentPosition + increment) % circle.Size
}

func (circle *MarbleCircle) AddMarble(value int) {

	var position1 int = circle.CalculateNextPosition(1)

	circle.Marbles = append(circle.Marbles, Marble{Value: -1})
	circle.Size++
	var position2 int = circle.CalculateNextPosition(2)

	fmt.Printf("Position -> %d\n", circle.CurrentPosition)
	fmt.Printf("Position + 1 -> %d\n", position1)
	fmt.Printf("Position + 2 -> %d\n", position2)
	//
	//	circle.Marbles = append(circle.Marbles, Marble{Value: -1})
	//	circle.Size++
	fmt.Printf("_________after append -> %s\n", circle.Marbles)
	copy(circle.Marbles[position2+1:], circle.Marbles[position2:circle.Size])
	fmt.Printf("aftercopyn -> %s\n", circle.Marbles)
	circle.Marbles[position1+1].Value = value
	fmt.Printf("ADDING -> %s\n", circle.Marbles)
	fmt.Printf("____END -> %d\n", position2)

	circle.CurrentPosition = position1 + 1
}

func play(circle MarbleCircle, players Players, lastValue int) {
	var newValue int = 1

	for newValue <= lastValue {
		circle.AddMarble(newValue)
		newValue++
		fmt.Println(circle.Marbles)
	}

	fmt.Printf("___:::_current player __  %d__\n", players.Points[players.Current])
	fmt.Printf("___:::_new value ___  %d__\n", newValue)
}

func main() {

	var numberOfPlayers, lastValue int
	circle := MarbleCircle{Size: 1, CurrentPosition: 0, Marbles: make([]Marble, 1)}

	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatal("You must supply a how many elves are going to play and how many points last marble worths.")
	}

	numberOfPlayers, _ = strconv.Atoi(args[0])
	lastValue, _ = strconv.Atoi(args[1])

	players := Players{Current: 1, NumberOfPlayers: numberOfPlayers}

	play(circle, players, lastValue)

	fmt.Printf("_%d\n", numberOfPlayers)
	fmt.Printf("_%d\n", lastValue)
	fmt.Printf("_%d\n", circle.Marbles[0].Value)
	fmt.Printf("_\n")
}
