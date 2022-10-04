package main

import (
	"fmt"
	"math/rand"
	"time"
	"strconv"
	"flag"
)

var (
	randGen *rand.Rand
	max *int
)

type Node struct {
	passengers int
	name       string

	next *Node
}

func nodeNew(p_passengers int, p_name string) *Node {
	return &Node{passengers: p_passengers, name: p_name, next: nil}
}

func (p_node *Node) add(p_toAdd *Node) *Node {
	node := p_node
	for node.next != nil {
		node = node.next
	}

	node.next = p_toAdd

	return p_toAdd
}

func (p_node *Node) print() {
	for node := p_node; node != nil; node = node.next {
		fmt.Printf("[ %v | %v ]", node.name, node.passengers)

		if node.next != nil {
			fmt.Print("+-")
		}
	}

	fmt.Println()
}

func (p_node *Node) length() int {
	len := 0
	for node := p_node; node != nil; node = node.next {
		len ++
	}

	return len
}

func (p_node *Node) updateList() *Node {
	var prev *Node = nil
	node := p_node
	for node != nil {
		if node.passengers > 0 {
			prev = node
			node = node.next

			continue
		}

		if prev != nil {
			prev.next = node.next
		} else {
			p_node = node.next
		}

		node = node.next
	}

	return p_node
}

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(source)

	max = flag.Int("max", 6, "Maximal number of cars generated")

	flag.Parse()
}

func main() {
	train := nodeNew(randGen.Intn(*max), "C1")

	count := randGen.Intn(*max)
	for i := 0; i < count; i ++ {
		train.add(nodeNew(randGen.Intn(*max), "C" + strconv.Itoa(i)))
	}

	fmt.Print("  ")
	train.print()

	train = train.updateList()

	fmt.Print("\nRemoved empty cars:\n  ")
	train.print()

	fmt.Printf("\nCars count: %v\n", train.length())
}
