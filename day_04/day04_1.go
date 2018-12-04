// Álvaro Castellano Vela 2018/12/04

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type GuardAction struct {
	dateTime time.Time
	id       int
	minute   int
	second   int
	action   string
}

type GuardActions []*GuardAction

func (x GuardActions) Len() int           { return len(x) }
func (x GuardActions) Less(i, j int) bool { return x[i].dateTime.Before(x[j].dateTime) }
func (x GuardActions) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

func processLine(line string, guardActions *[]*GuardAction) {

	var guard int
	var minute, second, guardString, action, dateString, timeSring string
	var dateTime time.Time

	var guardAction GuardAction

	guard = -1

	re := regexp.MustCompile("\\[([[:digit:]]{4}-[[:digit:]]{2}-[[:digit:]]{2}) ([[:digit:]]{2}):([[:digit:]]{2})\\] (Guard #[[:digit:]]+ )?(.+)$")
	match := re.FindAllStringSubmatch(line, -1)

	dateString = match[0][1]
	minute = match[0][2]
	second = match[0][3]
	guardString = match[0][4]
	action = match[0][5]

	timeSring = fmt.Sprintf("%s %s:%s", dateString, minute, second)

	fmt.Println(timeSring)
	dateTime, _ = time.Parse("0000-00-00 00:00", timeSring)

	if guardString != "" {
		guard, _ = strconv.Atoi(guardString[7:])
	}

	//fmt.Printf("Minute %d, Second %d\n", minute, second)
	//fmt.Printf("Guard %d\n", guard)
	//fmt.Printf("Action %s\n", action)
	//fmt.Printf("date %s - %d\n", timeSring, dateTime)

	guardAction.dateTime = dateTime
	guardAction.id = guard
	guardAction.minute, _ = strconv.Atoi(minute)
	guardAction.second, _ = strconv.Atoi(second)
	guardAction.action = action

	//fmt.Printf("%s\n", guardAction)
	*guardActions = append(*guardActions, &guardAction)
	//fmt.Printf("\n")
	//fmt.Printf("%s\n", guardAction)
}

func processFile(filename string, guards map[int]int) {

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
	for _, guardAction := range guardActions {
		fmt.Println(guardAction.dateTime)
		fmt.Println(guardAction.action)
	}
	fmt.Println("")
	sort.Sort(GuardActions(guardActions))
	for _, guardAction := range guardActions {
		fmt.Println(guardAction.dateTime)
		fmt.Println(guardAction.action)
	}
}

func main() {

	//guardActions := make([]GuardAction)
	guards := make(map[int]int)
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}

	filename := args[0]

	processFile(filename, guards)
	fmt.Printf("__\n")
}
