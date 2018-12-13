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
	Value    int
	Next     *Marble
	Previous *Marble
}

type MarbleCircle struct {
	CurrentMarble *Marble
	FirstMarble   *Marble
}

func (marbles *MarbleCircle) Show() {
	var current *Marble

	current = marbles.FirstMarble
	fmt.Printf("[ ")
	fmt.Printf("%d ", (*current).Value)
	current = (*current).Next
	for current != marbles.FirstMarble {
		fmt.Printf("%d ", (*current).Value)
		current = (*current).Next
	}
	fmt.Printf("]\n")
}

func (marbles *MarbleCircle) AddMarble(value int) {
	var newMarble Marble

	newMarble.Value = value

	// Next Position
	marbles.CurrentMarble = marbles.CurrentMarble.Next
	marbles.CurrentMarble.Next.Previous = &newMarble
	newMarble.Next = marbles.CurrentMarble.Next
	newMarble.Previous = marbles.CurrentMarble
	marbles.CurrentMarble.Next = &newMarble
	marbles.CurrentMarble = marbles.CurrentMarble.Next

}

func (marbles *MarbleCircle) RemoveMarble(offsetToRemove int) int {

	var value int
	var auxMarble *Marble

	for i := 0; i < offsetToRemove; i++ {
		marbles.CurrentMarble = marbles.CurrentMarble.Previous
	}
	value = marbles.CurrentMarble.Value
	auxMarble = marbles.CurrentMarble
	marbles.CurrentMarble = marbles.CurrentMarble.Next
	auxMarble.Previous.Next = marbles.CurrentMarble
	marbles.CurrentMarble.Previous = auxMarble.Previous

	auxMarble.Next = nil
	auxMarble.Previous = nil
	auxMarble.Value = -1

	return value
}

func play(circle MarbleCircle, players Players, lastValue int) int {
	var newValue int = 1
	players.Current = 0
	var highestScore = 0

	for newValue <= lastValue {
		if (newValue % 23) != 0 {
			circle.AddMarble(newValue)
		} else {
			players.Points[players.Current] += newValue
			players.Points[players.Current] += circle.RemoveMarble(7)
		}
		newValue++
		players.Current = (players.Current + 1) % players.NumberOfPlayers
	}
	for _, score := range players.Points {
		if highestScore < score {
			highestScore = score
		}
	}

	return highestScore
}

func main() {

	var numberOfPlayers, lastValue int
	var firstMarble Marble
	var marbles MarbleCircle

	firstMarble.Value = 0
	firstMarble.Next = &firstMarble
	firstMarble.Previous = &firstMarble
	marbles.CurrentMarble = &firstMarble
	marbles.FirstMarble = &firstMarble

	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatal("You must supply a how many elves are going to play and how many points last marble worths.")
	}

	numberOfPlayers, _ = strconv.Atoi(args[0])
	lastValue, _ = strconv.Atoi(args[1])

	players := Players{Current: 1, NumberOfPlayers: numberOfPlayers, Points: map[int]int{}}

	fmt.Printf("Score :%d\n", play(marbles, players, lastValue))

}
