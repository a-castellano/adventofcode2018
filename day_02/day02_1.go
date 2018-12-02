// Álvaro Castellano Vela 2018/12/02

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func count(line string) (int, int) {
	var two, three int
	letterCount := make(map[rune]int)
	for _, letter := range line {
		letterCount[letter] += 1
	}
	for _, count := range letterCount {
		switch count {
		case 2:
			two = 1
		case 3:
			three = 1
		}
	}
	return two, three
}

func processFile(filename string, twoAndThree map[int]int) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		two, three := count(scanner.Text())
		twoAndThree[2] += two
		twoAndThree[3] += three
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	twoAndThree := make(map[int]int)
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	twoAndThree[2] = 0
	twoAndThree[3] = 0

	filename := args[0]
	processFile(filename, twoAndThree)
	fmt.Printf("2's: %d 3's:%d\nChecksum %d\n", twoAndThree[2], twoAndThree[3], twoAndThree[2]*twoAndThree[3])
}
