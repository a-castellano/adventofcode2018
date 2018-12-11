// Álvaro Castellano Vela 2018/12/11

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

type Step struct {
	ID         byte
	status     string
	requires   map[byte]int
	workedTime int
	timeToWork int
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
					if x[j].status == "In progress" {
						result = false
					} else {
						if x[i].status == "In progress" {
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
			step.timeToWork = int(id) - 64 + 60
			stepsMap[id] = step
			steps = append(steps, step)
		}
		if _, ok := stepsMap[beforeId]; !ok {
			var step = new(Step)
			step.ID = beforeId
			step.status = "Blocked"
			step.requires = map[byte]int{}
			step.timeToWork = int(beforeId) - 64 + 60
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

func order(steps []*Step, workers int) ([]byte, int) {
	var order []byte
	var seconds int = -1
	workerStatus := make(map[byte]*Step)
	finishedStepsMap := make(map[byte]bool)

	var totalSteps, finishedSteps int = len(steps), 0

	sort.Sort(Steps(steps))
	for finishedSteps != totalSteps {

		seconds++

		var finishedID byte

		for index, worker := range workerStatus {
			worker.workedTime++
			if worker.workedTime == worker.timeToWork {
				worker.status = "Finished"
				delete(workerStatus, index)
			}
		}

		for _, step := range steps {
			if _, ok := finishedStepsMap[step.ID]; !ok && step.status == "Finished" {
				finishedStepsMap[step.ID] = true
				order = append(order, step.ID)
				finishedID = step.ID
				finishedSteps++

				for _, blockedStep := range steps {
					if _, ok := blockedStep.requires[finishedID]; ok {
						delete(blockedStep.requires, finishedID)
						if len(blockedStep.requires) == 0 {
							blockedStep.status = "Available"

						}
					}
				}

			}
		}
		for _, step := range steps {
			if step.status == "Available" && len(workerStatus) < workers {
				step.status = "In progress"
				workerStatus[step.ID] = step
			}
		}
		sort.Sort(Steps(steps))
	}
	return order, seconds
}

func main() {

	var steps []*Step

	args := os.Args[1:]

	if len(args) != 2 {
		log.Fatal("You must supply a file to process and the number of workers you want.")
	}

	filename := args[0]
	workers, err := strconv.Atoi(args[1])

	if err != nil {
		log.Fatal("Cannot workers.\n")
	}

	steps = processFile(filename)

	fmt.Printf("Workers: %d\n", workers)
	stepsOrder, seconds := order(steps, workers)
	fmt.Printf("Order: %s\n", stepsOrder)
	fmt.Printf("Seconds: %d\n", seconds)
}
