package main

import (
	"fmt"
	"os"
	"bufio"
	"math/rand"
	"time"
	"flag"
	"strconv"
)

var (
	randGen *rand.Rand

	outFile, inFile *string
	genRandNums     *bool

	min, max, maxVal *int
)

func bubblesort(p_list []int) {
	if len(p_list) < 2 {
		return
	}

	sorted := false
	for !sorted {
		sorted = true

		for i := 0; i < len(p_list) - 1; i ++ {
			if p_list[i + 1] < p_list[i] {
				p_list[i + 1], p_list[i] = p_list[i], p_list[i + 1]

				sorted = false
			}
		}
	}
}

func genNumFile(p_path string, p_min, p_max, p_maxVal int) error {
	size := p_min + randGen.Intn(p_max - p_min + 1)

	file, err := os.Create(p_path)
	if err != nil {
		return err
	}
	defer file.Close()

	for i := 0; i < size; i ++ {
		file.WriteString(strconv.Itoa(randGen.Intn(p_maxVal + 1)) + "\n")
	}

	return nil
}

func genSortedFile(p_path string, p_list []int) error {
	file, err := os.Create(p_path)
	if err != nil {
		return err
	}
	defer file.Close()

	bubblesort(p_list)

	for _, val := range p_list {
		file.WriteString(strconv.Itoa(val) + "\n")
	}

	return nil
}

func sortNumFile(p_inPath, p_outPath string) error {
	file, err := os.Open(p_inPath)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)

	list := []int{}
	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			return err
		}

		list = append(list, num)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	file.Close()

	if err = genSortedFile(p_outPath, list); err != nil {
		return err
	}

	return nil
}

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(source)

	inFile  = flag.String("in",  "nums.txt",   "The input file containing numbers to sort")
	outFile = flag.String("out", "sorted.txt", "The output file to contain sorted numbers")

	genRandNums = flag.Bool("gennums", false, "Generate a file (in) with random numbers")

	min    = flag.Int("min",    5,    "Minimal amount of numbers generated")
	max    = flag.Int("max",    12,   "Maximal amount of numbers generated")
	maxVal = flag.Int("maxval", 1024, "Biggest number to be generated")

	flag.Parse()
}

func main() {
	if *min > *max {
		fmt.Printf("Error: min (%v) is greater than max (%v)\n", *min, *max)

		os.Exit(1)
	}

	if *genRandNums {
		if err := genNumFile(*inFile, *min, *max, *maxVal); err != nil {
			fmt.Printf("Error: %v", err.Error())

			os.Exit(1)
		}

		fmt.Printf("Generated into '%v'\n", *inFile)
	} else {
		if err := sortNumFile(*inFile, *outFile); err != nil {
			fmt.Printf("Error: %v", err.Error())

			os.Exit(1)
		}

		fmt.Printf("Sorted '%v' into '%v'\n", *inFile, *outFile)
	}
}
