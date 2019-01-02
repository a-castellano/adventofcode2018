// Álvaro Castellano Vela 2019/01/01

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type Point struct {
	X int
	Y int
}

func fillMaze(regex string, maze *[12][12]rune, startPoint Point) {

	var point Point = startPoint

	var i int = 0
	var bracket bool = false
	for regex[i] != '$' && bracket == false {
		fmt.Printf("%s\n", string(regex[i]))
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
				fmt.Println(subString.String())

				fillMaze(subString.String(), maze, point)
				startOfSustring = offset + 1
			}
			//fillMaze(regex, maze, startOfSustring, endOfSustring-1)
			var subString strings.Builder
			subString.WriteString(regex[startOfSustring : endOfSustring-1])
			subString.WriteString(regex[endOfSustring:])
			fmt.Println(subString.String())
			fillMaze(subString.String(), maze, point)
			bracket = true
		}
		i++
	}
}

func createMap(filename string) {

	var maze [12][12]rune

	startPoint := Point{X: 6, Y: 6}

	for i := 0; i < 12; i++ {
		for j := 0; j < 12; j++ {
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

	maze[startPoint.X][startPoint.Y] = '.'
	fillMaze(regex[1:], &maze, startPoint)

	for i := 0; i < 12; i++ {
		for j := 0; j < 12; j++ {
			fmt.Printf("%s", string(maze[i][j]))
		}
		fmt.Printf("\n")
	}
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	createMap(filename)
}
