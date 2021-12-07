package day7

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
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
	file, err := os.Open(prefixPath + "/day7.txt")
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

func getMedian(numbers []int) int {
	var median int
	mid := len(numbers) / 2
	if math.Mod(float64(len(numbers)), 2) == 0 {
		median = (numbers[mid-1] + numbers[mid]) / 2
	} else {
		median = numbers[mid]
	}
	return median
}

func calculateFuel(numberMap map[int]int, numbers []int, mid int) int {
	v, ok := numberMap[mid]
	if ok {
		return v
	}
	totalFuel := 0
	for _, number := range numbers {
		n := int(math.Abs(float64(number - mid)))
		fuel := n * (n + 1) / 2
		// fmt.Println(number, n, fuel)
		totalFuel += fuel
	}
	numberMap[mid] = totalFuel
	return totalFuel
}

func Day7() {
	numbers := parseFile()
	sort.Ints(numbers)
	i := getMedian(numbers)
	numberMap := make(map[int]int, len(numbers))
	fmt.Println(i)
	for 0 <= i && i < len(numbers) {
		prev := calculateFuel(numberMap, numbers, i-1)
		curr := calculateFuel(numberMap, numbers, i)
		next := calculateFuel(numberMap, numbers, i+1)
		if prev > curr && curr < next {
			fmt.Println(curr)
			break
		} else if prev < curr {
			i -= 1
		} else if next < curr {
			i += 1
		}
	}
}
