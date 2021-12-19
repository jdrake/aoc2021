package day18

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Node struct {
	parent *Node
	left   *Node
	right  *Node
	value  int
}

func (node Node) String() string {
	if !node.IsTerminal() {
		return "[" + node.left.String() + ", " + node.right.String() + "]"
	} else {
		return fmt.Sprint(node.value)
	}
}

func (node *Node) IsTerminal() bool {
	return node.left == nil && node.right == nil
}

func (node *Node) Height(depth int) int {
	if node.IsTerminal() {
		return depth + 1
	} else {
		leftHeight := node.left.Height(depth)
		rightHeight := node.right.Height(depth)
		return int(math.Max(float64(leftHeight), float64(rightHeight))) + 1
	}
}

func (node *Node) Split() {
	left := node.value / 2
	right := int(math.Ceil(float64(node.value) / 2))
	node.left = &Node{
		parent: node,
		value:  left,
	}
	node.right = &Node{
		parent: node,
		value:  right,
	}
	node.value = -1
	// node.PrintRoot("split")
}

func (node *Node) AddLeft() {
	curr := node
	for curr.parent != nil {
		child := curr
		curr = curr.parent
		if curr.left != child {
			curr = curr.left
			for !curr.IsTerminal() {
				curr = curr.right
			}
			curr.value += node.value
			// curr.PrintRoot("explode add left")
			break
		}
	}
}

func (node *Node) AddRight() {
	curr := node
	for curr.parent != nil {
		child := curr
		curr = curr.parent
		if curr.right != child {
			curr = curr.right
			for !curr.IsTerminal() {
				curr = curr.left
			}
			curr.value += node.value
			// curr.PrintRoot("explode add right")
			break
		}
	}
}

func (node *Node) PrintRoot(label string) {
	curr := node
	for curr.parent != nil {
		curr = curr.parent
	}
	fmt.Printf("%20s %s\n", label, curr)
}

func (node *Node) Reduce(level int, reductionMethod string) (*Node, bool) {
	var reduced bool
	if node.IsTerminal() {
		if node.value >= 10 && reductionMethod == "split" {
			node.Split()
			reduced = true
		}
	} else {
		if level >= 4 && reductionMethod == "explode" {
			// Explode
			node.left.AddLeft()
			node.right.AddRight()
			node.left = nil
			node.right = nil
			node.value = 0
			// node.PrintRoot("explode")
			reduced = true
		} else {
			node.left, reduced = node.left.Reduce(level+1, reductionMethod)
			if !reduced {
				node.right, reduced = node.right.Reduce(level+1, reductionMethod)
			}
		}
	}
	return node, reduced
}

func parseTree(line string) *Node {
	stack := list.New()
	var root *Node
	for _, char := range strings.Split(line, "") {
		curr := stack.Back()
		if char == "[" {
			node := &Node{value: -1}
			if curr == nil {
				root = node
			} else {
				parent := curr.Value.(*Node)
				if parent.left == nil {
					parent.left = node
				} else {
					parent.right = node
				}
				node.parent = parent
			}
			stack.PushBack(node)
		} else if char == "]" {
			stack.Remove(curr)
		} else if char != "," {
			node := curr.Value.(*Node)
			value, _ := strconv.Atoi(char)
			if node.left == nil {
				node.left = &Node{
					parent: node,
					value:  value,
				}
			} else if node.right == nil {
				node.right = &Node{
					parent: node,
					value:  value,
				}
			} else {
				log.Fatal("unknown state")
			}
		}
	}
	return root
}

func parseFile(name string) []string {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/" + name + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}
	return lines
}

func (tree *Node) ReduceFull() *Node {
	method := "explode"
	for i := 0; i < 1000; i++ {
		// fmt.Println("action", i, method)
		_, reduced := tree.Reduce(0, method)
		height := tree.Height(0)
		if height < 6 {
			if reduced {
				method = "split"
			} else {
				break
			}
		} else {
			method = "explode"
		}
	}
	return tree
}

func (node *Node) Magnitude() int {
	if node.IsTerminal() {
		return node.value
	} else {
		return 3*node.left.Magnitude() + 2*node.right.Magnitude()
	}
}

func Main() {
	// var lines = []string{
	// 	"[[[[[9,8],1],2],3],4]",
	// 	"[7,[6,[5,[4,[3,2]]]]]",
	// 	"[[6,[5,[4,[3,2]]]],1]",
	// 	"[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
	// 	"[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]",
	// }
	// var lines = []string{
	// 	"[1,1]",
	// 	"[2,2]",
	// 	"[3,3]",
	// 	"[4,4]",
	// 	"[5,5]",
	// 	"[6,6]",
	// }

	lines := parseFile("input")

	var combos []string
	for _, a := range lines {
		for _, b := range lines {
			if a != b {
				combos = append(combos, "["+a+","+b+"]")
			}
		}
	}

	max := 0
	for _, line := range combos {
		tree := parseTree(line)
		tree.ReduceFull()
		if mag := tree.Magnitude(); mag > max {
			max = mag
			fmt.Println("new max", line, mag)
		}
	}

	// tree := parseTree("[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]")
	// fmt.Println(tree)
	// tree.ReduceFull()
	// fmt.Println(tree)
	// fmt.Println(tree.Magnitude())
}
