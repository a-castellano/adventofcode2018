// Álvaro Castellano Vela 2018/12/24

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

func checkFunctions(beforeString string, instructionString string, afterString string) int {
	f := []func([4]int, int, int, int) [4]int{addr, addi, mulr, muli, banr, bani, borr, bori, setr, seti, gtir, gtri, gtrr, eqir, eqri, eqrr}
	var matches, toReturn int

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

	for instruction, _ := range f {
		afterCandidate := f[instruction](beforeRegisters, instuction[1], instuction[2], instuction[3])
		var equal bool = true
		for i, _ := range afterCandidate {
			if afterCandidate[i] != afterRegisters[i] {
				equal = false
				break
			}
		}
		if equal {
			matches++
		}
	}

	if matches >= 3 {
		toReturn = 1
	}
	return toReturn
}

func processFile(filename string) int {

	var machedSamples int

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
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

		machedSamples += checkFunctions(beforeString, instructionString, afterString)
	}
	return machedSamples
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	samples := processFile(filename)
	fmt.Printf("Samples: %d\n", samples)
}
