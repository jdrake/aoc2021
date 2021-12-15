package day14

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Pair [2]string

func parseFile() (map[Pair]int, map[Pair]string, string) {
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
	pairs := make(map[Pair]int)
	insertionMap := make(map[Pair]string)
	row := 0
	lastChar := ""
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		} else if row == 0 {
			chars := strings.Split(line, "")
			lastChar = chars[len(chars)-1]
			for i := 0; i < len(chars)-1; i++ {
				pairs[Pair{chars[i], chars[i+1]}] += 1
			}
			row += 1
		} else {
			parts := strings.Fields(line)
			chars := strings.Split(parts[0], "")
			insertionMap[Pair{chars[0], chars[1]}] = parts[2]
		}
	}
	return pairs, insertionMap, lastChar
}

func Insert(pairs map[Pair]int, insertionMap map[Pair]string) map[Pair]int {
	newPairs := make(map[Pair]int)
	for pair, count := range pairs {
		char, exists := insertionMap[pair]
		if exists {
			newPairs[pair] -= count
			newPairs[Pair{pair[0], char}] += count
			newPairs[Pair{char, pair[1]}] += count
		}
	}
	for pair, count := range newPairs {
		pairs[pair] += count
	}
	return pairs
}

func MinMax(pairs map[Pair]int, lastChar string) (string, int, string, int) {
	chars := make(map[string]int)
	for pair, count := range pairs {
		chars[pair[0]] += count
	}
	chars[lastChar] += 1
	minChar, maxChar := "", ""
	minCount, maxCount := math.MaxInt, 0
	for char, count := range chars {
		if count < minCount {
			minCount = count
			minChar = char
			// fmt.Println("new min", char, count)
		}
		if count > maxCount {
			maxCount = count
			maxChar = char
			// fmt.Println("new max", char, count)
		}
	}

	return minChar, minCount, maxChar, maxCount
}

func Day14() {
	pairs, insertionMap, lastChar := parseFile()
	fmt.Println(pairs)
	for step := 0; step < 40; step++ {
		fmt.Println("-----")
		pairs = Insert(pairs, insertionMap)
		fmt.Println(pairs)
		minChar, minCount, maxChar, maxCount := MinMax(pairs, lastChar)
		fmt.Println(minChar, minCount, maxChar, maxCount)
		fmt.Println(maxCount - minCount)
	}
}
