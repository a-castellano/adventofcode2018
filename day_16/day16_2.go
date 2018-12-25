// Álvaro Castellano Vela 2018/12/25

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func addr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	registers[registerC] = registers[registerA] + registers[registerB]
	return registers
}

func addi(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	registers[registerC] = registers[registerA] + immediate
	return registers
}

func mulr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	registers[registerC] = registers[registerA] * registers[registerB]
	return registers
}

func muli(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	registers[registerC] = registers[registerA] * immediate
	return registers
}

func banr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	registers[registerC] = registers[registerA] & registers[registerB]
	return registers
}

func bani(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	registers[registerC] = registers[registerA] & immediate
	return registers
}

func borr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	registers[registerC] = registers[registerA] | registers[registerB]
	return registers
}

func bori(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	registers[registerC] = registers[registerA] | immediate
	return registers
}

func setr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	registers[registerC] = registers[registerA]
	return registers
}

func seti(registers [4]int, immediate int, noOneCares int, registerC int) [4]int {
	registers[registerC] = immediate
	return registers
}

func gtir(registers [4]int, immediate int, registerB int, registerC int) [4]int {
	if immediate > registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func gtri(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	if registers[registerA] > immediate {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func gtrr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	if registers[registerA] > registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqir(registers [4]int, immediate int, registerB int, registerC int) [4]int {
	if immediate == registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqri(registers [4]int, registerA int, immediate int, registerC int) [4]int {
	if registers[registerA] == immediate {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqrr(registers [4]int, registerA int, registerB int, registerC int) [4]int {
	if registers[registerA] == registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func checkFunctions(beforeString string, instructionString string, afterString string, functionIDs *map[int]int, options *map[int]map[int]bool, assignedFunctionIDs *map[int]int) {

	f := []func([4]int, int, int, int) [4]int{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}

	var hasOptions bool

	beforeRE := regexp.MustCompile("Before: \\[([[:digit:]]+), ([[:digit:]]+), ([[:digit:]]+), ([[:digit:]]+)\\]$")
	instructionRE := regexp.MustCompile("([[:digit:]]+) ([[:digit:]]+) ([[:digit:]]+) ([[:digit:]]+)$")
	afterRE := regexp.MustCompile("After:  \\[([[:digit:]]+), ([[:digit:]]+), ([[:digit:]]+), ([[:digit:]]+)\\]$")

	beforeMatch := beforeRE.FindAllStringSubmatch(beforeString, -1)
	instuctionMatch := instructionRE.FindAllStringSubmatch(instructionString, -1)
	afterMatch := afterRE.FindAllStringSubmatch(afterString, -1)

	var beforeRegisters [4]int
	for i, register := range beforeMatch[0][1:] {
		beforeRegisters[i], _ = strconv.Atoi(register)
	}

	var instuction [4]int
	for i, register := range instuctionMatch[0][1:] {
		instuction[i], _ = strconv.Atoi(register)
	}

	var afterRegisters [4]int
	for i, register := range afterMatch[0][1:] {
		afterRegisters[i], _ = strconv.Atoi(register)
	}

	if len((*options)[instuction[0]]) == 0 {
		hasOptions = false
		(*options)[instuction[0]] = make(map[int]bool)
	} else {
		hasOptions = true
	}

	if _, ok := (*functionIDs)[instuction[0]]; !ok { //We do not know yet what function it is
		candidates := make(map[int]bool)
		for instructionID, _ := range f {
			if _, ok := (*assignedFunctionIDs)[instructionID]; !ok {
				afterCandidate := f[instructionID](beforeRegisters, instuction[1], instuction[2], instuction[3])
				var equal bool = true
				for i, _ := range afterCandidate {
					if afterCandidate[i] != afterRegisters[i] {
						equal = false
						break
					}
				}
				if equal {
					if hasOptions == false {
						(*options)[instuction[0]][instructionID] = true
					} else {
						candidates[instructionID] = true
					}
				}
			}
		}
		if hasOptions == false {
			hasOptions = true
		} else {
			for currentOption, _ := range (*options)[instuction[0]] {
				if _, ok := candidates[currentOption]; !ok {
					delete((*options)[instuction[0]], currentOption)
				}
			}
		}
	}
}

func processFile(filename string, functionIDs *map[int]int, assignedFunctionIDs *map[int]int) {

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	options := make(map[int]map[int]bool)
	for i := 0; i < 16; i++ {
		options[i] = make(map[int]bool)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		beforeString := scanner.Text()
		scanner.Scan()
		instructionString := scanner.Text()
		scanner.Scan()
		afterString := scanner.Text()

		scanner.Scan()

		checkFunctions(beforeString, instructionString, afterString, functionIDs, &options, assignedFunctionIDs)
	}
	for functionID, mached := range options {
		if len(mached) == 1 {
			for opcode, _ := range mached {
				(*functionIDs)[functionID] = opcode
				(*assignedFunctionIDs)[opcode] = functionID
			}
		}
	}
}

func evaluateOpcodes(filename string, functionIDs map[int]int) [4]int {
	f := []func([4]int, int, int, int) [4]int{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	var registers [4]int = [4]int{0, 0, 0, 0}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		instructionString := scanner.Text()
		instructionRE := regexp.MustCompile("([[:digit:]]+) ([[:digit:]]+) ([[:digit:]]+) ([[:digit:]]+)$")
		instuctionMatch := instructionRE.FindAllStringSubmatch(instructionString, -1)
		var instuction [4]int
		for i, register := range instuctionMatch[0][1:] {
			instuction[i], _ = strconv.Atoi(register)
		}

		registers = f[functionIDs[instuction[0]]](registers, instuction[1], instuction[2], instuction[3])
	}

	return registers
}

func main() {
	args := os.Args[1:]
	functionIDs := make(map[int]int)
	assignedFunctionIDs := make(map[int]int)

	if len(args) != 2 {
		log.Fatal("You must supply a two files to process.")
	}
	rules := args[0]
	opcodes := args[1]
	for len(functionIDs) != 16 {
		processFile(rules, &functionIDs, &assignedFunctionIDs)
	}
	registers := evaluateOpcodes(opcodes, functionIDs)
	fmt.Printf("First Register: %d\n", registers[0])
}
