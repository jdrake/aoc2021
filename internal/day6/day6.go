package day6

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

func parseNumbers(s string) []int {
	var numbers []int
	for _, s := range strings.Split(s, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		numbers = append(numbers, v)
	}
	return numbers
}

func parseFile() []int {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/day6.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var numbers []int
	for scanner.Scan() {
		line := scanner.Text()
		numbers = parseNumbers(line)
		break
	}
	return numbers
}

func sum(numbers []int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func Day6() {
	numbers := parseFile()
	counts := make([]int, 9)
	for _, n := range numbers {
		counts[n] += 1
	}
	for day := 0; day < 256; day++ {
		fmt.Println(counts)
		newFishCount := counts[0]
		for i := 1; i <= 8; i++ {
			count := counts[i]
			if count > 0 {
				counts[i] -= count
				counts[i-1] += count
			}
		}
		counts[0] -= newFishCount
		counts[6] += newFishCount
		counts[8] += newFishCount
	}
	fmt.Println(sum(counts))
}
