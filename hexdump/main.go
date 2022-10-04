package main

import (
	"fmt"
	"os"
	"flag"
)

var lineSplit *int

func hexDump(p_path string, p_lineSplit int) error {
	data, err := os.ReadFile(p_path)
	if err != nil {
		return err
	}

	if len(data) == 0 {
		return nil
	}

	for i, char := range data {
		if i % p_lineSplit == 0 {
			if i != 0 {
				fmt.Println()
			}

			fmt.Printf("%08X ", i)
		}

		fmt.Printf("%02X ", char)
	}

	fmt.Printf("\n%08X\n", len(data))

	return nil
}

func init() {
	flag.Usage = func() {
		fmt.Printf("Usage: %v FILENAME...\n", os.Args[0])
		fmt.Println("Options:")

		flag.PrintDefaults()
	}

	lineSplit = flag.Int("linesplit", 16, "The byte at which the line gets split")

	flag.Parse()
}

func main() {
	if len(flag.Args()) == 0 {
		fmt.Println("Error: Expected a file name")
		fmt.Printf("Try '%v -h'\n", os.Args[0])

		return
	}

	for i, arg := range flag.Args() {
		if i > 0 {
			fmt.Printf("\n'%v':\n", arg)
		} else if len(flag.Args()) > 1 {
			fmt.Printf("'%v':\n", arg)
		}

		if err := hexDump(arg, *lineSplit); err != nil {
			fmt.Printf("Error: File '%v' does not exist\n", arg)
			fmt.Printf("Try '%v -h'\n", os.Args[0])
		}
	}
}
