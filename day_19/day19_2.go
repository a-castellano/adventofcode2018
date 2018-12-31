// Álvaro Castellano Vela 2018/12/31

// Thanks to https://rosettacode.org/wiki/Factors_of_an_integer

// My input was 10551410

package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func printFactors(nr int) {
	var sum int
	if nr < 1 {
		fmt.Println("\nFactors of", nr, "not computed")
		return
	}
	fs := make([]int, 1)
	fs[0] = 1
	apf := func(p int, e int) {
		n := len(fs)
		for i, pp := 0, p; i < e; i, pp = i+1, pp*p {
			for j := 0; j < n; j++ {
				fs = append(fs, fs[j]*pp)
			}
		}
	}
	e := 0
	for ; nr&1 == 0; e++ {
		nr >>= 1
	}
	apf(2, e)
	for d := 3; nr > 1; d += 2 {
		if d*d > nr {
			d = nr
		}
		for e = 0; nr%d == 0; e++ {
			nr /= d
		}
		if e > 0 {
			apf(d, e)
		}
	}
	for _, number := range fs {
		sum += number
	}
	fmt.Println("Result: ", sum)
}
func main() {
	args := os.Args[1:]
	if len(args) != 1 {
		log.Fatal("You must supply a number.")
	}
	number, _ := strconv.Atoi(args[0])
	printFactors(number)
}
