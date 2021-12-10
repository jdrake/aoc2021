package day10

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

func parseFile() [][]string {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var lines [][]string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, strings.Split(line, ""))
	}
	return lines
}

var charMap = map[string]string{
	")": "(",
	"}": "{",
	"]": "[",
	">": "<",
}

func FirstInvalidChar(chars []string) (string, *list.List) {
	stack := list.New()
	stack.PushBack(chars[0])
	for i := 1; i < len(chars); i++ {
		char := chars[i]
		openingBrace, isClosingBrace := charMap[char]
		if isClosingBrace {
			// Trying to close a brace
			el := stack.Back()
			if openingBrace == el.Value.(string) {
				// Brace is closed
				stack.Remove(el)
			} else {
				// invalid brace
				return char, stack
			}
		} else {
			// Add an opening brace
			stack.PushBack(char)
		}
	}
	return "", stack
}

var pointMap = map[string]int{
	"(": 1,
	"[": 2,
	"{": 3,
	"<": 4,
}

func StackScore(stack *list.List) int {
	score := 0
	braces := ""
	for stack.Len() > 0 {
		el := stack.Back()
		stack.Remove(el)
		brace := el.Value.(string)
		braces += brace
		score *= 5
		score += pointMap[brace]
	}
	fmt.Println(braces, score)
	return score
}

func Day10() {
	lines := parseFile()
	var scores []int
	for _, line := range lines {
		char, stack := FirstInvalidChar(line)
		if char == "" {
			score := StackScore(stack)
			scores = append(scores, score)
		}
	}
	sort.Ints(scores)
	midScore := scores[(len(scores)-1)/2]
	fmt.Println("midScore", midScore)
}
