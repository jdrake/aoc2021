package day3

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func parseLine(line string) []int {
	parts := strings.Split(line, "")
	var digits []int
	for _, v := range parts {
		digit, err := strconv.Atoi(v)
		if err != nil {
			log.Fatal(err)
		}
		digits = append(digits, digit)
	}
	return digits
}

func convertBinarySliceToInt(slice []int) int64 {
	var stringValues []string
	for _, v := range slice {
		stringValues = append(stringValues, strconv.Itoa(v))
	}
	s := strings.Join(stringValues, "")
	v, err := strconv.ParseInt(s, 2, 0)
	if err != nil {
		log.Fatal(err)
	}
	return v
}

func parseFile() [][]int {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/day3.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var lines [][]int
	for scanner.Scan() {
		line := scanner.Text()
		digits := parseLine(line)
		lines = append(lines, digits)
	}
	return lines
}

func getCounterForBitIndex(lines [][]int, bitIndex int) int {
	counter := 0
	for _, digits := range lines {
		bit := digits[bitIndex]
		if bit == 1 {
			counter += 1
		} else {
			counter -= 1
		}
	}
	return counter
}

func filterLines(lines [][]int, bitIndex int, invert bool) [][]int {
	counter := getCounterForBitIndex(lines, bitIndex)
	var bitPredicate int
	if invert {
		if counter < 0 {
			// 0 is more common, keep 1
			bitPredicate = 1
		} else {
			// 1 or neither is more common, keep 0
			bitPredicate = 0
		}
	} else {
		if counter >= 0 {
			// 1 or neither is more common, keep 1
			bitPredicate = 1
		} else {
			// 0 is more common, keep 0
			bitPredicate = 0
		}
	}
	fmt.Println("bitIndex", bitIndex, "counter", counter, "predicate", bitPredicate)
	var filteredLines [][]int
	for _, digits := range lines {
		if digits[bitIndex] == bitPredicate {
			filteredLines = append(filteredLines, digits)
		}
	}
	return filteredLines
}

func getRating(lines [][]int, invert bool) int64 {
	filteredLines := make([][]int, len(lines))
	copy(filteredLines, lines)
	fmt.Println("starting line count", len(filteredLines))
	for bitIndex := range lines[0] {
		filteredLines = filterLines(filteredLines, bitIndex, invert)
		fmt.Println("bitIndex", bitIndex, "line count", len(filteredLines))
		if len(filteredLines) == 1 {
			break
		}
	}
	fmt.Println(filteredLines)
	rating := convertBinarySliceToInt(filteredLines[0])
	return rating
}

func Day3() {
	lines := parseFile()
	oxygenRating := getRating(lines, false)
	fmt.Println(oxygenRating)
	co2Rating := getRating(lines, true)
	fmt.Println(co2Rating)
	result := oxygenRating * co2Rating
	fmt.Println(result)
}
