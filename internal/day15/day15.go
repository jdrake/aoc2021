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
	"time"
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
	size := y

	for i := 1; i < 5; i++ {
		for relY := 0; relY < size; relY++ {
			for relX := 0; relX < size; relX++ {
				baseP := Point{
					x: (i-1)*size + relX,
					y: relY,
				}
				p := Point{
					x: i*size + relX,
					y: relY,
				}
				value := points[baseP] + 1
				if value > 9 {
					value = 1
				}
				points[p] = value
			}
		}
	}

	for i := 0; i < 5; i++ {
		for j := 1; j < 5; j++ {
			for relY := 0; relY < size; relY++ {
				for relX := 0; relX < size; relX++ {
					baseP := Point{
						x: (i * size) + relX,
						y: (j-1)*size + relY,
					}
					p := Point{
						x: (i * size) + relX,
						y: j*size + relY,
					}
					value := points[baseP] + 1
					if value > 9 {
						value = 1
					}
					points[p] = value
				}
			}
		}
	}

	return points
}

func h(start Point, goal Point) int {
	return int(math.Abs(float64(start.x-goal.x))) + int(math.Abs(float64(start.y-goal.y)))
}

// func reconstructPath(cameFrom map[Point]Point, current Point) []Point {
// 	var path []Point
// 	path = append(path, current)
// 	for {
// 		next, exists := cameFrom[current]
// 		if exists {
// 			current = next
// 			path = append(path, current)
// 		} else {
// 			break
// 		}
// 	}
// 	return path
// }

type AStarError struct {
	Message string
}

func (e *AStarError) Error() string {
	return fmt.Sprintf("A* error: %s", e.Message)
}

var vectors = [4][2]int{
	{-1, 0},
	{0, -1},
	{0, 1},
	{1, 0},
}

func AStar(points map[Point]int) (int, error) {
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
		// fmt.Println("current", current.value, current.priority)
		if current.value == goal {
			return current.priority, nil
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

	return 0, &AStarError{"no path to goal"}
}

func PrintGrid(grid map[Point]int) {
	for y := 0; y < 50; y++ {
		row := ""
		for x := 0; x < 50; x++ {
			row += strconv.Itoa(grid[Point{x, y}])
		}
		fmt.Println(row)
	}
}

func Day15() {
	start := time.Now()
	points := parseFile()
	t1 := time.Now()
	elapsed := t1.Sub(start)
	fmt.Println(elapsed)
	// PrintGrid(points)
	fmt.Println(AStar(points))
	t2 := time.Now()
	elapsed = t2.Sub(t1)
	fmt.Println(elapsed)
}
