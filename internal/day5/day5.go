package day5

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

func parseCoordinatePair(s string) []int {
	var pair []int
	for _, s := range strings.Split(s, ",") {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		pair = append(pair, v)
	}
	return pair
}

func Max(nums ...int) int {
	var max int
	for _, num := range nums {
		if num > max {
			max = num
		}
	}
	return max
}

func parseFile() ([][][]int, int, int) {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/day5.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var coordinatePairs [][][]int
	var maxX int
	var maxY int
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)
		start := parseCoordinatePair(words[0])
		finish := parseCoordinatePair(words[2])
		maxX = Max(maxX, start[0], finish[0])
		maxY = Max(maxY, start[1], finish[1])
		coordinatePair := [][]int{start, finish}
		coordinatePairs = append(coordinatePairs, coordinatePair)
	}
	return coordinatePairs, maxX, maxY
}

func Day5() {
	coordinatePairs, maxX, maxY := parseFile()
	fmt.Println("max", maxX, maxY)
	// Create an empty grid
	grid := make([][]int, maxY+1)
	for i := range grid {
		grid[i] = make([]int, maxX+1)
	}
	for _, pair := range coordinatePairs {
		x1, y1 := pair[0][0], pair[0][1]
		x2, y2 := pair[1][0], pair[1][1]
		var startX int
		var finishX int
		var startY int
		var finishY int
		if x1 < x2 {
			startX, finishX = x1, x2
		} else {
			startX, finishX = x2, x1
		}
		if y1 < y2 {
			startY, finishY = y1, y2
		} else {
			startY, finishY = y2, y1
		}
		if startX == finishX {
			for i := startY; i <= finishY; i++ {
				grid[i][startX] += 1
			}
		} else if startY == finishY {
			for j := startX; j <= finishX; j++ {
				grid[startY][j] += 1
			}
		} else {
			fmt.Println("PAIR", pair)
			if x1 < x2 {
				if y1 < y2 {
					j := y1
					for i := x1; i <= x2; i++ {
						grid[j][i] += 1
						j += 1
					}
				} else {
					j := y1
					for i := x1; i <= x2; i++ {
						grid[j][i] += 1
						j -= 1
					}
				}
			} else {
				if y1 < y2 {
					j := y1
					for i := x1; i >= x2; i-- {
						grid[j][i] += 1
						j += 1
					}
				} else {
					j := y1
					for i := x1; i >= x2; i-- {
						grid[j][i] += 1
						j -= 1
					}
				}
			}
		}
	}
	count := 0
	for _, row := range grid {
		for _, v := range row {
			if v >= 2 {
				count += 1
			}
		}
	}
	// fmt.Println(grid)
	fmt.Println(count)
}
