package day9

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

func parseFile() [][]int {
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
	var grid [][]int
	for scanner.Scan() {
		line := scanner.Text()
		var row []int
		for _, n := range strings.Split(line, "") {
			v, _ := strconv.Atoi(n)
			row = append(row, v)
		}
		grid = append(grid, row)
	}
	return grid
}

func Search(grid [][]int) {
	risk := 0
	for r, row := range grid {
		for c, value := range row {
			aboveIsHigher := r-1 < 0 || grid[r-1][c] > value
			belowIsHigher := r+1 >= len(grid) || grid[r+1][c] > value
			leftIsHigher := c-1 < 0 || grid[r][c-1] > value
			rightIsHigher := c+1 >= len(row) || grid[r][c+1] > value
			if aboveIsHigher && belowIsHigher && leftIsHigher && rightIsHigher {
				risk += value + 1
			}
		}
	}
	fmt.Println(risk)
}

func Day9() {
	grid := parseFile()
	Search(grid)
}
