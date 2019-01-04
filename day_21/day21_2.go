// Álvaro Castellano Vela 2019/01/04

package main

import (
	"fmt"
)

func main() {

	var r2, r3, r5, last_result int
	var not_repeted bool = true

	results := make(map[int]bool)

	for not_repeted {
		r3 = r5 | 65536
		r5 = 8725355
		for true {
			r5 += r3 & 255
			r5 &= 16777215
			r5 *= 65899
			r5 &= 16777215
			if 256 > r3 {
				break
			}
			r2 = 0
			for (r2+1)*256 <= r3 {
				r2++
			}
			r3 = r2
		}
		if _, ok := results[r5]; !ok {
			results[r5] = true
			last_result = r5
		} else {
			not_repeted = false
		}
	}
	fmt.Printf("Lowest non-negative integer value for register 0 that causes the program to halt after executing the most instructions: %d\n", last_result)
}
