package day11

import (
	"bufio"
	"container/list"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type Octopus struct {
	x, y int
}

type Octopuses map[Octopus]int

func parseFile() Octopuses {
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
	octopuses := make(Octopuses)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		for x, s := range strings.Split(line, "") {
			value, _ := strconv.Atoi(s)
			octopuses[Octopus{x: x, y: y}] = value
		}
		y += 1
	}
	return octopuses
}

func (octopuses Octopuses) IncrementValues() {
	for octo := range octopuses {
		octopuses[octo] += 1
	}
}

func (octopuses Octopuses) String() string {
	s := ""
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			s += strconv.Itoa(octopuses[Octopus{x, y}])
		}
		s += "\n"
	}
	return s
}

var vectors = [8][2]int{
	{-1, -1},
	{-1, 0},
	{-1, 1},
	{0, -1},
	{0, 1},
	{1, -1},
	{1, 0},
	{1, 1},
}

func (octopuses Octopuses) Flash() int {
	q := list.New()
	flashed := mapset.NewSet()
	for octo, value := range octopuses {
		if value > 9 {
			flashed.Add(octo)
			q.PushBack(octo)
			octopuses[octo] = 0
		}
	}
	for q.Len() > 0 {
		el := q.Front()
		q.Remove(el)
		octo := el.Value.(Octopus)
		for _, vec := range vectors {
			x, y := vec[0], vec[1]
			nextOcto := Octopus{x: octo.x + x, y: octo.y + y}
			_, exists := octopuses[nextOcto]
			if exists && !flashed.Contains(nextOcto) {
				octopuses[nextOcto] += 1
				if octopuses[nextOcto] > 9 {
					flashed.Add(nextOcto)
					q.PushBack(nextOcto)
					octopuses[nextOcto] = 0
				}
			}
		}
	}
	flashes := flashed.Cardinality()
	return flashes
}

func Day11() {
	octopuses := parseFile()
	for i := 1; i < 1000; i++ {
		octopuses.IncrementValues()
		flashes := octopuses.Flash()
		if flashes == 100 {
			fmt.Println(i)
			break
		}
	}
}
