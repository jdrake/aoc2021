package day15

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func parseFile() map[Point]int {
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
	points := make(map[Point]int)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, s := range strings.Split(line, "") {
			value, _ := strconv.Atoi(s)
			points[Point{x, y}] = value
		}
		y += 1
	}
	return points
}

func h(start Point, goal Point) int {
	return int(math.Abs(float64(start.x-goal.x))) + int(math.Abs(float64(start.y-goal.y)))
}

func reconstructPath(cameFrom map[Point]Point, current Point) []Point {
	var path []Point
	path = append(path, current)
	for {
		next, exists := cameFrom[current]
		if exists {
			current = next
			path = append(path, current)
		} else {
			break
		}
	}
	return path
}

type AStarError struct {
	Message string
}

func (e *AStarError) Error() string {
	return fmt.Sprintf("A* error: %s", e.Message)
}

var vectors = [4][2]int{
	// {-1, -1},
	{-1, 0},
	// {-1, 1},
	{0, -1},
	{0, 1},
	// {1, -1},
	{1, 0},
	// {1, 1},
}

func AStar(points map[Point]int) ([]Point, error) {
	size := int(math.Sqrt(float64(len(points))))
	start, goal := Point{0, 0}, Point{size - 1, size - 1}
	openSet := make(PriorityQueue, 1)
	openSet[0] = &Item{
		value:    start,
		priority: h(start, goal),
		index:    0,
	}
	heap.Init(&openSet)
	openSetCache := make(map[Point]bool)
	openSetCache[start] = true

	cameFrom := make(map[Point]Point)

	gScore := make(map[Point]int)
	gScore[start] = 0

	fScore := make(map[Point]int)
	fScore[start] = h(start, goal)

	for openSet.Len() > 0 {
		current := heap.Pop(&openSet).(*Item)
		delete(openSetCache, current.value)
		fmt.Println("current", current.value, current.priority)
		if current.value == goal {
			return reconstructPath(cameFrom, current.value), nil
		}

		for _, vector := range vectors {
			x := current.value.x + vector[0]
			y := current.value.y + vector[1]
			if x < 0 || y < 0 || x >= size || y >= size {
				continue
			}
			neighbor := Point{x, y}
			tentativeGScore := gScore[current.value] + points[neighbor]
			neighborGScore, exists := gScore[neighbor]
			if !exists {
				gScore[neighbor] = math.MaxInt
				neighborGScore = math.MaxInt
			}
			if tentativeGScore < neighborGScore {
				cameFrom[neighbor] = current.value
				gScore[neighbor] = tentativeGScore
				fScore[neighbor] = tentativeGScore + h(neighbor, goal)
				if _, exists := openSetCache[neighbor]; !exists {
					item := &Item{
						value:    neighbor,
						priority: fScore[neighbor],
					}
					heap.Push(&openSet, item)
					openSetCache[neighbor] = true
				}
			}
		}
	}

	return nil, &AStarError{"no path to goal"}
}

func Day15() {
	points := parseFile()
	fmt.Println(AStar(points))
}
