// Álvaro Castellano Vela 2019/01/07

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Nanobot struct {
	Range int
	X     int
	Y     int
	Z     int
}

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func processFile(filename string) ([]Nanobot, int) {

	var nanobots []Nanobot
	var maxRangeIndex, maxRange int

	var numberOfNanobots int

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var nanobot Nanobot
		nanobotString := scanner.Text()
		re := regexp.MustCompile("pos=<([-]?[[:digit:]]+),([-]?[[:digit:]]+),([-]?[[:digit:]]+)>, r=([[:digit:]]+)$")
		match := re.FindAllStringSubmatch(nanobotString, -1)

		nanobot.Range, _ = strconv.Atoi(match[0][4])
		nanobot.X, _ = strconv.Atoi(match[0][1])
		nanobot.Y, _ = strconv.Atoi(match[0][2])
		nanobot.Z, _ = strconv.Atoi(match[0][3])

		if nanobot.Range > maxRange {
			maxRange = nanobot.Range
			maxRangeIndex = numberOfNanobots
		}

		nanobots = append(nanobots, nanobot)
		numberOfNanobots++
	}
	return nanobots, maxRangeIndex
}

func getNanobotsInRange(nanobots []Nanobot, maxRangeIndex int) int {

	var nanobotsInRange int
	var maxRangeNanobot Nanobot = nanobots[maxRangeIndex]

	for _, nanobot := range nanobots {
		var distance int
		distance = abs(nanobot.X-maxRangeNanobot.X) + abs(nanobot.Y-maxRangeNanobot.Y) + abs(nanobot.Z-maxRangeNanobot.Z)
		if distance <= maxRangeNanobot.Range {
			nanobotsInRange++
		}
	}

	return nanobotsInRange
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	nanobots, maxRangeIndex := processFile(filename)

	nanobotsInRange := getNanobotsInRange(nanobots, maxRangeIndex)
	fmt.Printf("Nanobots in range of the nanobot with the biggest range: %d\n", nanobotsInRange)
}
