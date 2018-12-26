// Álvaro Castellano Vela 2018/12/26

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Point struct {
	Y int
	X int
}

type Stream struct {
	Position Point
	Ended    bool
}

type ReservoirResearch struct {
	Ground  *[2000][2000]byte
	Streams []Stream
}

func renderLine(line string, ground *[2000][2000]byte) (int, int) {

	re := regexp.MustCompile("(x|y)=([[:digit:]]+), (x|y)=([[:digit:]]+)..([[:digit:]]+)$")
	match := re.FindAllStringSubmatch(line, -1)

	firstCoordinate := match[0][1]
	firstCoordinateValue, _ := strconv.Atoi(match[0][2])
	secondCoordinateStartRange, _ := strconv.Atoi(match[0][4])
	secondCoordinateEndRange, _ := strconv.Atoi(match[0][5])

	var xStart, xEnd, yStart, yEnd int

	if firstCoordinate == "x" {
		xStart = firstCoordinateValue
		xEnd = firstCoordinateValue
		yStart = secondCoordinateStartRange
		yEnd = secondCoordinateEndRange
	} else {
		yStart = firstCoordinateValue
		yEnd = firstCoordinateValue
		xStart = secondCoordinateStartRange
		xEnd = secondCoordinateEndRange
	}
	for i := yStart; i <= yEnd; i++ {
		for j := xStart; j <= xEnd; j++ {
			(*ground)[i][j] = '#'
		}
	}
	(*ground)[0][500] = '+'

	return yEnd, yStart
}

func renderGround(filename string, ground *[2000][2000]byte) (int, int) {

	var maxY int = 0
	var minY int = 2001

	for i := 0; i < 2000; i++ {
		for j := 0; j < 2000; j++ {
			(*ground)[i][j] = '.'
		}
	}

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		var yCandidateEnd, yCandidateStart int = renderLine(scanner.Text(), ground)
		if yCandidateEnd > maxY {
			maxY = yCandidateEnd
		}
		if yCandidateStart > maxY {
			maxY = yCandidateStart
		}
		if yCandidateEnd < minY {
			minY = yCandidateEnd
		}
		if yCandidateStart < minY {
			minY = yCandidateStart
		}

	}

	return maxY, minY
}

func stream(ground *[2000][2000]byte, maxY int) {
	var reservoirResearch ReservoirResearch
	reservoirResearch.Ground = ground
	reservoirResearch.Streams = append(reservoirResearch.Streams, Stream{Position: Point{Y: 1, X: 500}, Ended: false})

	var allFilled bool = false

	for allFilled == false {
		for streamID, stream := range reservoirResearch.Streams {
			// Down until # or maxY
			if stream.Ended == false {
				for reservoirResearch.Ground[stream.Position.Y+1][stream.Position.X] != '#' && stream.Position.Y+1 <= maxY {
					reservoirResearch.Ground[stream.Position.Y][stream.Position.X] = '|'
					stream.Position.Y++
				}
				// Fill

				for stream.Ended == false {
					var offset int
					var fillLine bool = true
					var rightLimit, leftLimit int

					for offset = stream.Position.X; reservoirResearch.Ground[stream.Position.Y][offset] != '#'; offset++ {
						var valueUnder byte = reservoirResearch.Ground[stream.Position.Y+1][offset]
						if valueUnder != '~' && valueUnder != '#' {
							fillLine = false
							break
						}
					}
					if reservoirResearch.Ground[stream.Position.Y][offset] != '#' {
						fillLine = false
					}
					rightLimit = offset - 1
					for offset = stream.Position.X; reservoirResearch.Ground[stream.Position.Y][offset] != '#'; offset-- {
						var valueUnder byte = reservoirResearch.Ground[stream.Position.Y+1][offset]
						if valueUnder != '~' && valueUnder != '#' {
							fillLine = false
							break
						}
					}
					if reservoirResearch.Ground[stream.Position.Y][offset] != '#' {
						fillLine = false
					}
					leftLimit = offset + 1
					if fillLine {
						for offset = leftLimit; offset <= rightLimit; offset++ {
							reservoirResearch.Ground[stream.Position.Y][offset] = '~'
						}
						stream.Position.Y--
					} else {
						for offset = leftLimit; offset <= rightLimit; offset++ {
							reservoirResearch.Ground[stream.Position.Y][offset] = '|'
						}
						reservoirResearch.Streams[streamID].Ended = true
						reservoirResearch.Streams[streamID].Position.Y = stream.Position.Y
						reservoirResearch.Streams[streamID].Position.X = offset
						stream.Ended = true
						if reservoirResearch.Ground[stream.Position.Y][leftLimit-1] == '.' {
							reservoirResearch.Ground[stream.Position.Y][leftLimit-1] = '|'
							reservoirResearch.Streams = append(reservoirResearch.Streams, Stream{Position: Point{Y: stream.Position.Y, X: leftLimit - 1}, Ended: false})
						}
						if reservoirResearch.Ground[stream.Position.Y][rightLimit+1] == '.' {
							reservoirResearch.Ground[stream.Position.Y][rightLimit+1] = '|'
							reservoirResearch.Streams = append(reservoirResearch.Streams, Stream{Position: Point{Y: stream.Position.Y, X: rightLimit + 1}, Ended: false})
						}
					}
				}
			}
		}
		allFilled = true
		for _, stream := range reservoirResearch.Streams {
			if stream.Ended == false {
				allFilled = false
				break
			}
		}
	}
}

func main() {
	var ground [2000][2000]byte
	var maxY int
	var minY int
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	maxY, minY = renderGround(filename, &ground)
	stream(&ground, maxY)
	//for i := 0; i < 2000; i++ {
	//	fmt.Println(string(ground[i][:2000]))
	//}
	var count int
	for i := minY; i <= maxY; i++ {
		for j := 0; j < 2000; j++ {
			if ground[i][j] == '~' || ground[i][j] == '|' {
				count++
			}
		}
	}
	fmt.Printf("Final Count: %d\n", count)
}
