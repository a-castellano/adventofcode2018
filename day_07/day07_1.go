// Álvaro Castellano Vela 2018/12/07

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
)

type Step struct {
	ID       byte
	status   string
	requires map[byte]int
}

type Steps []*Step

func (x Steps) Len() int           { return len(x) }
func (x Steps) Less(i, j int) bool { return false } //DEFINE
func (x Steps) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func processLine(line string) (byte, byte) {
	var state, beforeState byte

	re := regexp.MustCompile("Step ([[:alpha:]]) must be finished before step ([[:alpha:]]) can begin.$")
	match := re.FindAllStringSubmatch(line, -1)

	state = match[0][1][0]
	beforeState = match[0][2][0]

	return state, beforeState
}

func processFile(filename string) ([]*Step, map[byte]*Step) {

	var steps []*Step
	stepsMap := make(map[byte]*Step)

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		var id, beforeId byte
		id, beforeId = processLine(scanner.Text())
		if _, ok := stepsMap[id]; !ok {
			var step = new(Step)
			step.ID = id
			step.status = "Available"
			step.requires = map[byte]int{}
			stepsMap[id] = step
			steps = append(steps, step)
		}
		if _, ok := stepsMap[beforeId]; !ok {
			var step = new(Step)
			step.ID = beforeId
			step.status = "Blocked"
			step.requires = map[byte]int{}
			stepsMap[beforeId] = step
			steps = append(steps, step)
		}
		stepsMap[beforeId].status = "Blocked"
		stepsMap[beforeId].requires[id] = 1
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return steps, stepsMap
}

func main() {

	var stepts []*Step
	var stepsMap map[byte]*Step

	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	stepts, stepsMap = processFile(filename)

	fmt.Printf("%c\n", stepts[1].ID)
	fmt.Printf("%c\n", stepsMap['C'].ID)
	for _, step := range stepts {
		fmt.Printf("State %c:\n", step.ID)
		fmt.Printf("State %c is %s:\n", step.ID, step.status)
		for requiredStep, _ := range step.requires {
			fmt.Printf("\tRequires %c:\n", requiredStep)
		}
	}
	fmt.Printf("_\n")
}
