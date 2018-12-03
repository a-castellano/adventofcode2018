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

func findNoOverlap(overlapMap map[int]bool) {
	for index, value := range overlapMap {
		if value == false {
			fmt.Printf("%d does not overlap.\n", index)
		}
	}
}

func processLine(line string, fabric *[1000][1000]int, overlapMap map[int]bool) {
	var claimID, xPos, yPos, width, height int
	var overlap bool = false

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
				overlapMap[fabric[i][j]] = true
				fabric[i][j] = -1
				overlap = true
			}
		}
	}
	if overlap {
		overlapMap[claimID] = true
	} else {
		overlapMap[claimID] = false
	}
}

func processFile(filename string, fabric *[1000][1000]int, overlap map[int]bool) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		processLine(scanner.Text(), fabric, overlap)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	var fabric [1000][1000]int
	overlap := make(map[int]bool)

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]
	processFile(filename, &fabric, overlap)
	findNoOverlap(overlap)
}
