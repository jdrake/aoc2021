package day13

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

type Point struct {
	x, y int
}

type Fold struct {
	direction string
	value     int
}

type Grid struct {
	points map[Point]int
	folds  []Fold
}

func parseFile() Grid {
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
	var grid = Grid{
		points: make(map[Point]int),
	}
	parsingPoints := true
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			parsingPoints = false
		} else {
			if parsingPoints {
				parts := strings.Split(line, ",")
				x, _ := strconv.Atoi(parts[0])
				y, _ := strconv.Atoi(parts[1])
				point := Point{x, y}
				grid.points[point] = 1
			} else {
				fields := strings.Fields(line)
				parts := strings.Split(fields[len(fields)-1], "=")
				value, _ := strconv.Atoi(parts[1])
				grid.folds = append(grid.folds, Fold{
					direction: parts[0],
					value:     value,
				})
			}
		}
	}
	return grid
}

func (grid Grid) String() string {
	s, maxX, maxY := "", 0, 0
	for point := range grid.points {
		if point.x > maxX {
			maxX = point.x
		}
		if point.y > maxY {
			maxY = point.y
		}
	}
	for y := 0; y <= maxY; y++ {
		for x := 0; x <= maxX; x++ {
			_, exists := grid.points[Point{x, y}]
			if exists {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

func (grid Grid) Fold(direction string, value int) {
	if direction == "y" {
		for point := range grid.points {
			if point.y > value {
				grid.points[Point{point.x, value - (point.y - value)}] += 1
				delete(grid.points, point)
			}
		}
	} else {
		for point := range grid.points {
			if point.x > value {
				grid.points[Point{value - (point.x - value), point.y}] += 1
				delete(grid.points, point)
			}
		}
	}
}

func Day13() {
	grid := parseFile()
	// fmt.Println(grid)
	for _, fold := range grid.folds {
		grid.Fold(fold.direction, fold.value)
		// break
	}
	fmt.Println(grid)
	count := 0
	for _, value := range grid.points {
		if value > 0 {
			count += 1
		}
	}
	fmt.Println(count)
}
