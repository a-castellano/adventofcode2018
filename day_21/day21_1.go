// Álvaro Castellano Vela 2019/01/03

package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Instruction struct {
	Name      string
	RegisterA int
	RegisterB int
	RegisterC int
}

func addr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	registers[registerC] = registers[registerA] + registers[registerB]
	return registers
}

func addi(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	registers[registerC] = registers[registerA] + immediate
	return registers
}

func mulr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	registers[registerC] = registers[registerA] * registers[registerB]
	return registers
}

func muli(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	registers[registerC] = registers[registerA] * immediate
	return registers
}

func banr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	registers[registerC] = registers[registerA] & registers[registerB]
	return registers
}

func bani(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	registers[registerC] = registers[registerA] & immediate
	return registers
}

func borr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	registers[registerC] = registers[registerA] | registers[registerB]
	return registers
}

func bori(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	registers[registerC] = registers[registerA] | immediate
	return registers
}

func setr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	registers[registerC] = registers[registerA]
	return registers
}

func seti(registers [6]int, immediate int, noOneCares int, registerC int) [6]int {
	registers[registerC] = immediate
	return registers
}

func gtir(registers [6]int, immediate int, registerB int, registerC int) [6]int {
	if immediate > registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func gtri(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	if registers[registerA] > immediate {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func gtrr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	if registers[registerA] > registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqir(registers [6]int, immediate int, registerB int, registerC int) [6]int {
	if immediate == registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqri(registers [6]int, registerA int, immediate int, registerC int) [6]int {
	if registers[registerA] == immediate {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func eqrr(registers [6]int, registerA int, registerB int, registerC int) [6]int {
	if registers[registerA] == registers[registerB] {
		registers[registerC] = 1
	} else {
		registers[registerC] = 0
	}
	return registers
}

func processFile(filename string) ([]Instruction, [6]int, int) {

	var instructions []Instruction
	var registers [6]int

	var ip int

	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()
	ipString := scanner.Text()
	ipStringRe := regexp.MustCompile("#ip ([[:digit:]]+)$")
	matchIp := ipStringRe.FindAllStringSubmatch(ipString, -1)

	ip, _ = strconv.Atoi(matchIp[0][1])

	for scanner.Scan() {
		var instruction Instruction
		instructionString := scanner.Text()
		re := regexp.MustCompile("([[:alpha:]]+) ([[:digit:]]+) ([[:digit:]]+) ([[:digit:]]+)$")
		match := re.FindAllStringSubmatch(instructionString, -1)

		instruction.Name = match[0][1]
		instruction.RegisterA, _ = strconv.Atoi(match[0][2])
		instruction.RegisterB, _ = strconv.Atoi(match[0][3])
		instruction.RegisterC, _ = strconv.Atoi(match[0][4])

		instructions = append(instructions, instruction)
	}
	return instructions, registers, ip
}

func run(instructions []Instruction, registers [6]int, ip int) [6]int {

	f := map[string]func([6]int, int, int, int) [6]int{"addr": addr, "addi": addi, "mulr": mulr, "muli": muli, "banr": banr, "bani": bani, "borr": borr, "bori": bori, "setr": setr, "seti": seti, "gtir": gtir, "gtri": gtri, "gtrr": gtrr, "eqir": eqir, "eqri": eqri, "eqrr": eqrr}
	//maxInstruction := len(instructions) - 1

	//for registers[ip] <= maxInstruction {
	for i := 0; i < 370700; i++ {
		fmt.Printf("%s %d %d %d\n", instructions[registers[ip]].Name, instructions[registers[ip]].RegisterA, instructions[registers[ip]].RegisterB, instructions[registers[ip]].RegisterC)
		registers = f[instructions[registers[ip]].Name](registers, instructions[registers[ip]].RegisterA, instructions[registers[ip]].RegisterB, instructions[registers[ip]].RegisterC)
		fmt.Println(registers)
		registers[ip]++
		fmt.Println(registers)
	}

	return registers
}

func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a file to process.")
	}
	filename := args[0]
	instructions, registers, ip := processFile(filename)
	//yyp	registers[0] = 4682012
	//registers[0] = 7282971
	registers = run(instructions, registers, ip)
	fmt.Printf("Register 0 has %d value.\n", registers[0])
}
