package day9

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
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

type Coord struct {
	row, col int
}

func FindLowPoints(coords map[Coord]int, maxRow int, maxCol int) []Coord {
	var lowPointCoords []Coord
	for coord, value := range coords {
		aboveIsHigher := coord.row-1 < 0 || coords[Coord{row: coord.row - 1, col: coord.col}] > value
		belowIsHigher := coord.row+1 >= maxRow || coords[Coord{row: coord.row + 1, col: coord.col}] > value
		leftIsHigher := coord.col-1 < 0 || coords[Coord{row: coord.row, col: coord.col - 1}] > value
		rightIsHigher := coord.col+1 >= maxCol || coords[Coord{row: coord.row, col: coord.col + 1}] > value
		if aboveIsHigher && belowIsHigher && leftIsHigher && rightIsHigher {
			lowPointCoords = append(lowPointCoords, coord)
		}
	}
	return lowPointCoords
}

func MakeCoords(grid [][]int) map[Coord]int {
	coords := make(map[Coord]int)
	for r, row := range grid {
		for c, value := range row {
			coords[Coord{row: r, col: c}] = value
		}
	}
	fmt.Println(coords)
	return coords
}

func FindBasin(coords map[Coord]int, lowPoint Coord, maxRow int, maxCol int) int {
	fmt.Println("----")
	fmt.Println("lowPoint", lowPoint)
	q := list.New()
	q.PushBack(lowPoint)
	visited := mapset.NewSet()
	visited.Add(lowPoint)
	basinSize := 0
	for q.Len() > 0 {
		el := q.Front()
		q.Remove(el)
		coord := el.Value.(Coord)
		if coords[coord] < 9 {
			fmt.Println("added coord to basin", coord, coords[coord])
			basinSize += 1
			if coord.row-1 >= 0 {
				newCoord := Coord{row: coord.row - 1, col: coord.col}
				if !visited.Contains(newCoord) {
					q.PushBack(newCoord)
					visited.Add(newCoord)
				}
			}
			if coord.row+1 < maxRow {
				newCoord := Coord{row: coord.row + 1, col: coord.col}
				if !visited.Contains(newCoord) {
					q.PushBack(newCoord)
					visited.Add(newCoord)
				}
			}
			if coord.col-1 >= 0 {
				newCoord := Coord{row: coord.row, col: coord.col - 1}
				if !visited.Contains(newCoord) {
					q.PushBack(newCoord)
					visited.Add(newCoord)
				}
			}
			if coord.col+1 < maxCol {
				newCoord := Coord{row: coord.row, col: coord.col + 1}
				if !visited.Contains(newCoord) {
					q.PushBack(newCoord)
					visited.Add(newCoord)
				}
			}
		}
	}
	fmt.Println("basinSize", basinSize)
	return basinSize
}

func Day9() {
	grid := parseFile()
	coords := MakeCoords(grid)
	lowPoints := FindLowPoints(coords, len(grid), len(grid[0]))
	fmt.Println(lowPoints)
	var basinSizes []int
	for _, lowPoint := range lowPoints {
		basinSize := FindBasin(coords, lowPoint, len(grid), len(grid[0]))
		basinSizes = append(basinSizes, basinSize)
	}
	sort.Ints(basinSizes)
	fmt.Println("basinSizes", basinSizes)
	total := 1
	for i := len(basinSizes) - 3; i < len(basinSizes); i++ {
		total *= basinSizes[i]
	}
	fmt.Println("total", total)
}
