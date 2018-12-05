// Álvaro Castellano Vela 2018/12/05

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type GuardAction struct {
	dateSeconds int64
	id          int
	hour        int
	minute      int
	action      string
	line        string
}

type GuardActions []*GuardAction

func (x GuardActions) Len() int           { return len(x) }
func (x GuardActions) Less(i, j int) bool { return x[i].dateSeconds < x[j].dateSeconds }
func (x GuardActions) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func processLine(line string, guardActions *[]*GuardAction) {

	var minute, hour, day, month, year int
	var guardString, action string

	var guardAction GuardAction

	re := regexp.MustCompile("\\[([[:digit:]]{4})-([[:digit:]]{2})-([[:digit:]]{2}) ([[:digit:]]{2}):([[:digit:]]{2})\\] (Guard #[[:digit:]]+ )?(.+)$")
	match := re.FindAllStringSubmatch(line, -1)

	year, _ = strconv.Atoi(match[0][1])
	month, _ = strconv.Atoi(match[0][2])
	day, _ = strconv.Atoi(match[0][3])
	hour, _ = strconv.Atoi(match[0][4])
	minute, _ = strconv.Atoi(match[0][5])
	guardString = match[0][6]
	action = match[0][7]

	if guardString != "" {
		guardString = strings.Trim(guardString[7:], " ")
		guardAction.id, _ = strconv.Atoi(guardString)
	} else {
		guardAction.id = -1
	}

	guardAction.dateSeconds = int64(minute*60 + hour*3600 + day*86400 + month*2678400 + year*32140800)
	guardAction.minute = minute
	guardAction.hour = hour
	guardAction.action = action
	guardAction.line = line

	*guardActions = append(*guardActions, &guardAction)
}

func processFile(filename string) []*GuardAction {

	var guardActions []*GuardAction

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		processLine(scanner.Text(), &guardActions)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	sort.Sort(GuardActions(guardActions))

	return guardActions
}

func findSleepyhead(guardActions *[]*GuardAction, guards map[int]map[int]int) int {
	actions := len(*guardActions)
	var i, maxId int = 0, -1
	var moreMinuesAslept, maxTimesAslept int = -1, -1

	var guardID = -1
	for i < actions-1 {
		var preAction, postAction string
		if (*guardActions)[i].id > 0 {
			guardID = (*guardActions)[i].id
		}
		i++
		preAction = (*guardActions)[i].action
		if preAction == "falls asleep" {
			i++
			postAction = (*guardActions)[i].action
			if postAction == "wakes up" {
				if guards[guardID] == nil {
					guards[guardID] = map[int]int{}
				}
				for j := (*guardActions)[i-1].minute; j < (*guardActions)[i].minute; j++ {
					guards[guardID][j]++
				}
			}
		}
	}
	for guardID, _ := range guards {
		var candidateMaxTimes = -1
		var candidateMinute = -1
		for minute, timesAslept := range guards[guardID] {
			if candidateMaxTimes < timesAslept {
				candidateMaxTimes = timesAslept
				candidateMinute = minute
			}
		}
		if maxTimesAslept < candidateMaxTimes {
			maxTimesAslept = candidateMaxTimes
			moreMinuesAslept = candidateMinute
			maxId = guardID
		}
	}

	return maxId * moreMinuesAslept
}

func main() {

	var guardActions []*GuardAction
	var sleepyheadID int
	guards := make(map[int]map[int]int)
	args := os.Args[1:]

	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	guardActions = processFile(filename)
	sleepyheadID = findSleepyhead(&guardActions, guards)

	fmt.Printf("Guard ID -> %d\n", sleepyheadID)
}
