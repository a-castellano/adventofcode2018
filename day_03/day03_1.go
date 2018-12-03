// Álvaro Castellano Vela 2018/12/03

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func countOverclaimed(fabric *[1000][1000]int) int {
	var count int
	for i := 0; i < 1000; i++ {
		for j := 0; j < 1000; j++ {
			if fabric[i][j] == -1 {
				count++
			}
		}
	}

	return count
}

func processLine(line string, fabric *[1000][1000]int) {
	var claimID, xPos, yPos, width, height int

	re := regexp.MustCompile("#([[:digit:]]+) @ ([[:digit:]]+),([[:digit:]]+): ([[:digit:]]+)x([[:digit:]]+)$")
	match := re.FindAllStringSubmatch(line, -1)

	claimID, _ = strconv.Atoi(match[0][1])
	xPos, _ = strconv.Atoi(match[0][2])
	yPos, _ = strconv.Atoi(match[0][3])
	width, _ = strconv.Atoi(match[0][4])
	height, _ = strconv.Atoi(match[0][5])

	for i := xPos; i < xPos+width; i++ {
		for j := yPos; j < yPos+height; j++ {
			if fabric[i][j] == 0 {
				fabric[i][j] = claimID
			} else {
				fabric[i][j] = -1
			}
		}
	}
}

func processFile(filename string, fabric *[1000][1000]int) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		processLine(scanner.Text(), fabric)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	var fabric [1000][1000]int
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]
	processFile(filename, &fabric)
	fmt.Printf("Overclaimed square inches: %d\n", countOverclaimed(&fabric))
}
