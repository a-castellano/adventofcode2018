// Álvaro Castellano Vela 2018/12/14

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Pot struct {
	Index    int
	Next     *Pot
	Previous *Pot
	PotType  byte
}

type Pots struct {
	FirstPot *Pot
	LastPot  *Pot
}

func processPotsFromLine(line string) Pots {

	var pots Pots
	var pot Pot
	var currentPot *Pot
	var initialStateString string
	var index int = 0

	re := regexp.MustCompile("initial state: ([#.]+)")
	match := re.FindAllStringSubmatch(line, -1)

	initialStateString = match[0][1]

	pot.PotType = initialStateString[0]
	pot.Previous = nil
	currentPot = &pot

	pots.FirstPot = &pot
	pot.Index = index

	for _, potByte := range initialStateString[1:] {
		index++
		var newPot Pot
		newPot.PotType = byte(potByte)
		newPot.Index = index

		newPot.Previous = currentPot
		currentPot.Next = &newPot

		currentPot = currentPot.Next
	}

	pots.LastPot = currentPot
	currentPot.Next = nil

	return pots
}

func addRuleFromLine(rules *map[string]byte, line string) {

	re := regexp.MustCompile("([#.]{5}) => (#|.)$")
	match := re.FindAllStringSubmatch(line, -1)
	(*rules)[match[0][1]] = match[0][2][0]
}

func processFile(filename string) (Pots, map[string]byte) {

	var pots Pots
	rules := make(map[string]byte)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	pots = processPotsFromLine(scanner.Text())
	scanner.Scan()

	for scanner.Scan() {
		addRuleFromLine(&rules, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pots, rules
}

func (pot *Pot) getString() string {
	s := make([]byte, 5)
	s[0] = pot.Previous.Previous.PotType
	s[1] = pot.Previous.PotType
	s[2] = pot.PotType
	s[3] = pot.Next.PotType
	s[4] = pot.Next.Next.PotType

	return string(s)
}

func (pots *Pots) generation(rules *map[string]byte) {

	var pot, first, last *Pot
	newValues := make(map[int]byte)

	for i := 0; i < 4; i++ {
		pot = new(Pot)
		pot.Index = pots.FirstPot.Index - 1
		pot.PotType = 46
		pot.Next = pots.FirstPot
		pots.FirstPot.Previous = pot
		pots.FirstPot = pot
	}
	pot.Previous = nil
	first = pots.FirstPot.Next.Next

	for i := 0; i < 4; i++ {
		pot = new(Pot)
		pot.Index = pots.LastPot.Index + 1
		pot.PotType = 46
		pot.Previous = pots.LastPot
		pots.LastPot.Next = pot
		pots.LastPot = pot
	}
	pot.Next = nil
	last = pots.LastPot.Previous.Previous

	for matchPot := first; matchPot != last; matchPot = matchPot.Next {

		if newValue, _ := (*rules)[matchPot.getString()]; newValue != 0 {
			newValues[matchPot.Index] = newValue
		} else {
			newValues[matchPot.Index] = 46
		}
	}
	for matchPot := first; matchPot != last; matchPot = matchPot.Next {
		matchPot.PotType = newValues[matchPot.Index]
	}

}

func main() {
	var plants int = 0

	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]
	pots, rules := processFile(filename)
	for i := 0; i < 20; i++ {
		pots.generation(&rules)
	}

	for pot := pots.FirstPot; pot != nil; pot = pot.Next {
		if pot.PotType == 35 {
			plants += pot.Index
		}
	}
	fmt.Printf("Total Plants: %d\n", plants)

}
