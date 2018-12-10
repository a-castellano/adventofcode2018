// Álvaro Castellano Vela 2018/12/07

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
)

type Step struct {
	ID       byte
	status   string
	requires map[byte]int
}

type Steps []*Step

func (x Steps) Len() int      { return len(x) }
func (x Steps) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x Steps) Less(i, j int) bool {

	var result bool

	if x[i].status == "Finished" && x[j].status == "Finished" {
		result = false
	} else {
		if x[i].status == x[j].status {
			result = x[i].ID < x[j].ID
		} else {
			if x[j].status == "Finished" {
				result = false
			} else {
				if x[i].status == "Finished" {
					result = true
				} else {
					if x[j].status == "Available" {
						result = false
					} else {
						if x[i].status == "Available" {
							result = true
						}
					}
				}
			}
		}
	}

	return result
}

func processLine(line string) (byte, byte) {
	var state, beforeState byte

	re := regexp.MustCompile("Step ([[:alpha:]]) must be finished before step ([[:alpha:]]) can begin.$")
	match := re.FindAllStringSubmatch(line, -1)

	state = match[0][1][0]
	beforeState = match[0][2][0]

	return state, beforeState
}

func processFile(filename string) []*Step {

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

	return steps
}

func order(steps []*Step) []byte {
	var order []byte
	var totalSteps, finishedSteps int = len(steps), 0
	sort.Sort(Steps(steps))
	for finishedSteps != totalSteps {
		var finishedID byte
		for _, step := range steps {
			if step.status == "Available" {
				step.status = "Finished"
				order = append(order, step.ID)
				finishedSteps++
				finishedID = step.ID
				break
			}
		}
		for _, step := range steps {
			if _, ok := step.requires[finishedID]; ok {
				delete(step.requires, finishedID)
				if len(step.requires) == 0 {
					step.status = "Available"
				}
			}
		}
		sort.Sort(Steps(steps))
	}
	return order
}

func main() {

	var steps []*Step

	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	steps = processFile(filename)

	fmt.Printf("Order: %s\n", order(steps))
}
