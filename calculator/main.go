package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

func isOperator(p_op string) bool {
	switch p_op {
	case "+", "-", "*", "/", "%", "^": return true
	default: return false
	}
}

func calc(p_num1 float64, p_num2 float64, p_op string) float64 {
	switch p_op {
	case "+": return p_num1 + p_num2
	case "-": return p_num1 - p_num2
	case "*": return p_num1 * p_num2
	case "/": return p_num1 / p_num2
	case "%": return float64(int64(p_num1) % int64(p_num2))
	case "^": return math.Pow(p_num1, p_num2)

	default: panic("Operator is not valid")
	}

	return 0
}

func printOptions() {
	fmt.Print("[+  -  *  /  %  ^]")
}

func getNumPrompt() (float64, error) {
	fmt.Print("Number: ")

	var num float64
	_, err := fmt.Scan(&num)

	if err != nil {
		fmt.Println("  -> Error: Expected a number\n")
		fmt.Scanln()

		return 0, err
	}

	return num, nil
}

func getNumArgs(p_idx int) (float64, error) {
	num, err := strconv.ParseFloat(os.Args[p_idx], 64)
	if err != nil {
		fmt.Printf("Error: Expected number, got '%v'\n", os.Args[p_idx])
		tryUsage()

		return 0, err
	}

	return num, nil
}

func mainLoop() {
	quit := false

	for !quit {
		printOptions()
		fmt.Println(", ^C/q to quit\n")
		fmt.Print("> ")

		var op string
		fmt.Scanf("%s", &op)

		if !isOperator(op) {
			if op == "q" {
				quit = true
			}

			continue
		}

		var num1, num2 float64
		var err error

		fmt.Print("  ")
		num1, err = getNumPrompt()
		if err != nil {
			continue
		}

		fmt.Print("  ")
		num2, err = getNumPrompt()
		if err != nil {
			continue
		}

		result := calc(num1, num2, op)
		fmt.Printf("\n    %v %v %v = %v\n\n", num1, op, num2, result)
	}
}

func tryUsage() {
	fmt.Printf("Try '%v -h' for usage\n", os.Args[0])
}

func usage() {
	fmt.Println("Usage: app [NUM NUM OP] | [OPTIONS]")
	fmt.Println("Options:")
	fmt.Println("  -h   Show this help message")
}

func calcFromArgs() {
	var num1, num2 float64
	var err error

	num1, err = getNumArgs(1)
	if err != nil {
		os.Exit(1)
	}

	op := os.Args[2]
	if !isOperator(op) {
		fmt.Print("Error: Expected ")
		printOptions()
		fmt.Printf(", got '%v'\n", op)
		tryUsage()

		os.Exit(1)
	}

	num2, err = getNumArgs(3)
	if err != nil {
		os.Exit(1)
	}

	result := calc(num1, num2, op)
	fmt.Printf("%v %v %v = %v\n", num1, op, num2, result)
}

func main() {
	if len(os.Args) > 1 {
		if len(os.Args) == 2 && os.Args[1] == "-h" {
			usage()
		} else if len(os.Args) == 4 {
			calcFromArgs()
		} else {
			fmt.Println("Error: Unexpected arguments")
			tryUsage()

			os.Exit(1)
		}
	} else {
		mainLoop()
	}
}
