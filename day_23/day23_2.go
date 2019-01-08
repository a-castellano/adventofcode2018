// Álvaro Castellano Vela 2019/01/08

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
)

type Point struct {
	X int
	Y int
	Z int
}

type Nanobot struct {
	Range int
	Point Point
}

func abs(x int) int {
	if x < 0 {
		x = -x
	}
	return x
}

func (p Point) WalkPartway(t Point, d int) Point {
	xDelta := t.X - p.X
	yDelta := t.Y - p.Y
	zDelta := t.Z - p.Z

	sum := 0
	if xDelta < 0 {
		sum -= xDelta
	} else {
		sum += xDelta
	}
	if yDelta < 0 {
		sum -= yDelta
	} else {
		sum += yDelta
	}
	if zDelta < 0 {
		sum -= zDelta
	} else {
		sum += zDelta
	}

	return Point{p.X + d*xDelta/sum, p.Y + d*yDelta/sum, p.Z + d*zDelta/sum}

}

func (p Point) Distance(t Point) int {

	xDelta := p.X - t.X
	yDelta := p.Y - t.Y
	zDelta := p.Z - t.Z
	if xDelta < 0 {
		xDelta = -xDelta
	}
	if yDelta < 0 {
		yDelta = -yDelta
	}
	if zDelta < 0 {
		zDelta = -zDelta
	}
	return xDelta + yDelta + zDelta

}

func (n Nanobot) inRange(t Point) bool {
	return n.Point.Distance(t) <= n.Range
}

func (n Nanobot) distanceToRange(p Point) int {
	return n.Point.Distance(p) - n.Range
}

type Nanobots []Nanobot

func (nanobots Nanobots) inRange(t Point) int {
	inRange := 0
	for _, d := range nanobots {
		if d.inRange(t) {
			inRange++
		}
	}
	return inRange
}

type Candidate struct {
	Point   Point
	Quality int
}

func processFile(filename string) (Nanobots, map[Point]int) {

	var nanobots Nanobots
	seeds := make(map[Point]int)

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
		nanobot.Point.X, _ = strconv.Atoi(match[0][1])
		nanobot.Point.Y, _ = strconv.Atoi(match[0][2])
		nanobot.Point.Z, _ = strconv.Atoi(match[0][3])

		nanobots = append(nanobots, nanobot)
	}

	for _, nanobot := range nanobots {
		rangeToNanobot := 0
		for _, nearNanobot := range nanobots {
			if nearNanobot.inRange(nanobot.Point) {
				rangeToNanobot++
			}
		}
		seeds[nanobot.Point] = rangeToNanobot
	}

	return nanobots, seeds
}

func largestBumberOfNanobots(nanobots Nanobots, seeds map[Point]int) {
	ranges := make(map[Point]int)
	worklist := make([]Candidate, 0)

	for l, n := range seeds {
		worklist = append(worklist, Candidate{l, n})
	}
	highScore := 0
	for len(worklist) > 0 {
		sort.Slice(worklist, func(i, j int) bool {
			return worklist[i].Quality > worklist[j].Quality
		})

		c := worklist[0]
		point := c.Point
		worklist = worklist[1:]

		if ranges[point] != 0 {
			// Don't re-process items.
			continue
		}

		alreadyIn := 0
		candidates := make([]Candidate, 0)

		for _, nanobot := range nanobots {
			if nanobot.inRange(point) {
				// We don't need to get any closer.
				alreadyIn++
				ranges[point]++
				continue
			}
			dist := nanobot.distanceToRange(point)
			mid := point.WalkPartway(nanobot.Point, dist+1)
			candidates = append(candidates, Candidate{mid, dist})
		}

		if ranges[point] < highScore {
			// This isn't worth bothering with.
			continue
		} else if ranges[point] > highScore {
			highScore = ranges[point]
		}

		sort.Slice(candidates, func(i, j int) bool {
			return candidates[i].Quality < candidates[j].Quality
		})

		for _, c := range candidates {
			worklist = append(worklist, Candidate{c.Point, ranges[point] + 1})
		}
	}

	finalists := make([]Point, 0)
	for k, v := range ranges {
		if v != highScore {
			continue
		}
		finalists = append(finalists, k)
	}

	champions := make(map[Point]int)
	for _, s := range finalists {
		localHigh := highScore

		points := []*int{&s.X, &s.Y, &s.Z}
		for i := range points {
			first := points[(i+1)%3]
			second := points[(i+2)%3]

			for _, fdirection := range []int{1, -1} {
				for _, sdirection := range []int{1, -1} {
					oFirst := *first
					oSecond := *second

					for {
						*first += fdirection
						*second += sdirection

						var inRange int
						if ranges[s] == 0 {
							inRange = nanobots.inRange(s)
							ranges[s] = inRange
						} else {
							inRange = ranges[s]
						}

						if inRange < localHigh {
							*first = oFirst
							*second = oSecond
							break
						} else if inRange > localHigh {
							localHigh = inRange
							oFirst = *first
							oSecond = *second
						}
					}
				}
			}
		}
		champions[s] = localHigh
	}

	score := 0
	lowestSum := 9999999999999
	var solution Point
	for k, v := range champions {
		if v < score {
			continue
		} else if v > score {
			lowestSum = 9999999999999
			score = v
		}

		mySum := k.X + k.Y + k.Z
		if mySum < lowestSum {
			lowestSum = mySum
			solution = k
		}
	}

	fmt.Printf("Solution: %d\n", solution.X+solution.Y+solution.Z)

}

func getNanobotsUnderPoint(nanobots Nanobots, point Point) int {

	var nanobotsInRange int

	for _, nanobot := range nanobots {
		var distance int
		distance = abs(point.X-nanobot.Point.X) + abs(point.Y-nanobot.Point.Y) + abs(point.Z-nanobot.Point.Z)
		//fmt.Println("Distance: ", distance, " Range", nanobot.Range)
		if distance <= nanobot.Range {
			nanobotsInRange++
		}
	}

	return nanobotsInRange
}

func getMostInRange(nanobots []Nanobot, maxRangeIndex int) Point {

	var point Point
	var currentNanobotsInRange int
	var maxRange = nanobots[maxRangeIndex].Range

	for x := 0 - (maxRange)/1000; x <= maxRange/1000; x++ {
		fmt.Println(x)
		for y := 0 - (maxRange)/1000; y <= maxRange/1000; y++ {
			for z := 0 - (maxRange)/1000; z <= maxRange/1000; z++ {
				var candidatePoint Point = Point{X: x, Y: y, Z: z}
				nanobotsInRange := getNanobotsUnderPoint(nanobots, candidatePoint)
				if nanobotsInRange > currentNanobotsInRange {
					currentNanobotsInRange = nanobotsInRange
					point = candidatePoint
				}
			}
		}
	}
	return point
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	nanobots, seeds := processFile(filename)
	largestBumberOfNanobots(nanobots, seeds)
	//	fmt.Println(nanobots)

}
