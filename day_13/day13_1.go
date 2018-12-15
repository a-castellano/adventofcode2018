// Álvaro Castellano Vela 2018/12/15

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

/*

  |  ->  124
  -  ->  45
  /  ->  47
  \  ->  92
  +  ->  43
  ^  ->  94
  v  ->  118
  <  ->  60
  >  ->  62

*/

type Cart struct {
	X              int
	Y              int
	Turn           int
	Symbol         byte
	PreviousSymbol byte
}

type Carts []Cart

func (x Carts) Len() int      { return len(x) }
func (x Carts) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Carts) Less(i, j int) bool {

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

func processFile(filename string) ([][]byte, []Cart) {

	var tracks [][]byte
	var carts []Cart

	var currentRow int = 0

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		track := []byte(scanner.Text())
		tracks = append(tracks, track)
		for column, symbol := range track {
			if symbol == 94 || symbol == 118 || symbol == 60 || symbol == 62 {
				var cart Cart
				cart.X = currentRow
				cart.Y = column
				cart.Turn = 0
				cart.Symbol = symbol
				switch {
				case symbol == 94 || symbol == 118:
					cart.PreviousSymbol = 124
				case symbol == 60 || symbol == 62:
					cart.PreviousSymbol = 45
				}
				carts = append(carts, cart)
			}
		}
		currentRow++
	}

	return tracks, carts
}

func processNextMove(cart *Cart) (byte, int, int, int) {

	var offsetX, offsetY int = 0, 0
	var nextSymbol byte
	var nextTurn int = (*cart).Turn

	switch (*cart).PreviousSymbol {
	case 124: // |
		nextSymbol = (*cart).Symbol
		offsetY = 0
		switch (*cart).Symbol {
		case 94: // ^
			offsetX = -1
		case 118: //v
			offsetX = 1
		}
	case 45: // -
		nextSymbol = (*cart).Symbol
		offsetX = 0
		switch (*cart).Symbol {
		case 60: // <
			offsetY = -1
		case 62: // >
			offsetY = 1
		}
	case 47: // /
		switch (*cart).Symbol {
		case 94: // ^
			nextSymbol = 62 // >
			offsetX = 0
			offsetY = 1
		case 118: // v
			nextSymbol = 60 // <
			offsetX = 0
			offsetY = -1
		case 60: // <
			nextSymbol = 118 // v
			offsetX = 1
			offsetY = 0
		case 62: // >
			nextSymbol = 94 // ^
			offsetX = -1
			offsetY = 0
		}
	case 92: // \
		switch (*cart).Symbol {
		case 94: // ^
			nextSymbol = 60 // <
			offsetX = 0
			offsetY = -1
		case 118: // v
			nextSymbol = 62 // >
			offsetX = 0
			offsetY = 1
		case 60: // <
			nextSymbol = 94 // ^
			offsetX = -1
			offsetY = 0
		case 62: // >
			nextSymbol = 118 // v
			offsetX = 1
			offsetY = 0
		}
	case 43: // +
		switch (*cart).Turn {
		case 0: // left
			switch (*cart).Symbol {
			case 94: // ^
				nextSymbol = 60 // <
				offsetX = 0
				offsetY = -1
			case 118: // v
				nextSymbol = 62 // >
				offsetX = 0
				offsetY = 1
			case 60: // <
				nextSymbol = 118 // v
				offsetX = 1
				offsetY = 0
			case 62: // >
				nextSymbol = 94 // ^
				offsetX = -1
				offsetY = -0
			}
		case 1: // straight
			switch (*cart).Symbol {
			case 94: // ^
				nextSymbol = 94 // ^
				offsetX = -1
				offsetY = 0
			case 118: // v
				nextSymbol = 118 // v
				offsetX = 1
				offsetY = 0
			case 60: // <
				nextSymbol = 60 // <
				offsetX = 0
				offsetY = -1
			case 62: // >
				nextSymbol = 62 // >
				offsetX = 0
				offsetY = 1
			}
		case 2: // right
			switch (*cart).Symbol {
			case 94: // ^
				nextSymbol = 62 // >
				offsetX = 0
				offsetY = 1
			case 118: // v
				nextSymbol = 60 // <
				offsetX = 0
				offsetY = -1
			case 60: // <
				nextSymbol = 94 // ^
				offsetX = -1
				offsetY = 0
			case 62: // >
				nextSymbol = 118 // v
				offsetX = 1
				offsetY = 0
			}
		}
		nextTurn = (nextTurn + 1) % 3
	}

	return nextSymbol, nextTurn, offsetX, offsetY
}

func run(tracks *[][]byte, carts *[]Cart) {

	var collision bool = false

	for second := 0; collision == false; second++ {
		for i, _ := range *carts {

			var nextPositionSymbol byte

			nextSymbol, nextTurn, offsetX, offsetY := processNextMove(&(*carts)[i])
			(*tracks)[(*carts)[i].X][(*carts)[i].Y] = (*carts)[i].PreviousSymbol
			(*carts)[i].X += offsetX
			(*carts)[i].Y += offsetY
			(*carts)[i].Turn = nextTurn
			(*carts)[i].PreviousSymbol = (*tracks)[(*carts)[i].X][(*carts)[i].Y]
			(*carts)[i].Symbol = nextSymbol

			nextPositionSymbol = (*tracks)[(*carts)[i].X][(*carts)[i].Y]
			if nextPositionSymbol == 94 || nextPositionSymbol == 118 || nextPositionSymbol == 60 || nextPositionSymbol == 62 {
				collision = true
				fmt.Printf("Collision at %d,%d\n", (*carts)[i].Y, (*carts)[i].X)
				(*tracks)[(*carts)[i].X][(*carts)[i].Y] = 88 // X
				break
			} else {
				(*tracks)[(*carts)[i].X][(*carts)[i].Y] = nextSymbol
			}
		}

		sort.Sort(Carts(*carts))
	}

}

func main() {

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	tracks, carts := processFile(filename)
	run(&tracks, &carts)
}
