// Álvaro Castellano Vela 2018/12/14

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type FuelCell struct {
	X         int
	Y         int
	ID        int
	PoweLevel int
}

type FuelGrid struct {
	Cells        [301][301]FuelCell
	SerialNumber int
	MaxCell      int
	MaxCellX     int
	MaxCellY     int
}

func (grid *FuelGrid) fillGrid() {

	for i := 1; i < 301; i++ {
		for j := 1; j < 301; j++ {
			(*grid).Cells[i][j].X = i
			(*grid).Cells[i][j].Y = j
			(*grid).Cells[i][j].ID = i + 10
			(*grid).Cells[i][j].PoweLevel = ((*grid).Cells[i][j].ID*j + (*grid).SerialNumber) * (*grid).Cells[i][j].ID
			(*grid).Cells[i][j].PoweLevel = ((*grid).Cells[i][j].PoweLevel/100)%10 - 5
		}
	}

}

func (grid *FuelGrid) getSquareValue(x int, y int) int {

	var result int
	for i := x; i < x+3; i++ {
		for j := y; j < y+3; j++ {
			result += (*grid).Cells[i][j].PoweLevel
		}
	}

	return result
}

func (grid *FuelGrid) getLargest() {

	for i := 1; i < 299; i++ {
		for j := 1; j < 299; j++ {
			var squareValue = grid.getSquareValue(i, j)
			if squareValue > (*grid).MaxCell {
				(*grid).MaxCell = squareValue
				(*grid).MaxCellX = i
				(*grid).MaxCellY = j
			}
		}
	}

}

func main() {

	var grid FuelGrid

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a serial number.")
	}
	grid.SerialNumber, _ = strconv.Atoi(args[0])
	grid.fillGrid()
	grid.getLargest()
	fmt.Printf("Largest total power: %d\n", grid.MaxCell)
	fmt.Printf("Largest total power coordinates: %d,%d\n", grid.MaxCellX, grid.MaxCellY)

}
