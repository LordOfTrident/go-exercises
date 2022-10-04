package main

import (
	"fmt"
	"math/rand"
	"time"
	"flag"
	"os"
)

var (
	min, max, depth *int
	invert *bool
)

var randGen *rand.Rand

type Node struct {
	data int

	left, right *Node
}

func nodeNew(p_data int) *Node {
	return &Node{data: p_data, left: nil, right: nil}
}

func (p_node *Node) addLeft(p_data int) *Node {
	p_node.left = nodeNew(p_data);

	return p_node.left
}

func (p_node *Node) addRight(p_data int) *Node {
	p_node.right = nodeNew(p_data);

	return p_node.right
}

func (p_node *Node) invert() {
	tmp := p_node.right
	p_node.right = p_node.left
	p_node.left  = tmp
}

func (p_node *Node) invertTree() {
	p_node.invert()

	if p_node.right != nil {
		p_node.right.invertTree()
	}

	if p_node.left != nil {
		p_node.left.invertTree()
	}
}

func (p_node *Node) printTree(p_indent int) {
	for i := 0; i < p_indent; i ++ {
		fmt.Print("  ")
	}

	fmt.Println(p_node.data)

	if p_node.right != nil {
		p_node.right.printTree(p_indent + 1)
	}

	if p_node.left != nil {
		p_node.left.printTree(p_indent + 1)
	}
}

func (p_node *Node) addRandBranches(p_depth, p_min, p_max int) {
	if (p_depth == 0) {
		return
	}

	p_depth --

	p_node.addLeft(genRandInt(p_min, p_max)).addRandBranches(p_depth, p_min, p_max)
	p_node.addRight(genRandInt(p_min, p_max)).addRandBranches(p_depth, p_min, p_max)
}

func genRandInt(p_min, p_max int) int {
	if p_min > p_max {
		return 0
	}

	return p_min + randGen.Intn(p_max - p_min + 1)
}

func init() {
	source := rand.NewSource(time.Now().UnixNano())
	randGen = rand.New(source)

	min    = flag.Int("min",     0,    "The minimal value of a node")
	max    = flag.Int("max",     255,  "The maximal value of a node")
	depth  = flag.Int("depth",   2,    "The amount of branches the tree has")
	invert = flag.Bool("invert", true, "Should an inverted variant be printed too")

	flag.Parse()
}

func main() {
	if *min > *max {
		fmt.Printf("Error: min > max (%v > %v)\n", *min, *max)

		os.Exit(1)
	}

	baseNode := nodeNew(genRandInt(*min, *max))
	baseNode.addRandBranches(*depth, *min, *max)

	baseNode.printTree(0)

	if *invert {
		baseNode.invertTree()

		fmt.Println("\nInverted:")
		baseNode.printTree(1)
	}
}
