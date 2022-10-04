package main

import (
	"fmt"
	"flag"
)

var max *int

func init() {
	max = flag.Int("max", 21, "The biggest fibonacci number to be displayed")

	flag.Parse()
}

func main() {
	num1, num2 := 0, 1
	fmt.Printf("%v, %v", num1, num2)

	for true {
		num3 := num1 + num2
		if num3 > *max {
			break
		}

		fmt.Printf(", %v", num3)

		num1 = num2
		num2 = num3
	}

	fmt.Println()
}
