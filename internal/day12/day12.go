package day12

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

type Cave struct {
	name string
}

func (c Cave) IsSmall() bool {
	return strings.ToLower(c.name) == c.name
}

type Graph struct {
	caves map[Cave][]Cave
	paths [][]Cave
}

func parseFile() Graph {
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
	var graph = Graph{
		caves: make(map[Cave][]Cave),
		paths: make([][]Cave, 0),
	}
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		caveFrom := Cave{parts[0]}
		caveTo := Cave{parts[1]}
		_, caveFromExists := graph.caves[caveFrom]
		if caveFromExists {
			graph.caves[caveFrom] = append(graph.caves[caveFrom], caveTo)
		} else {
			graph.caves[caveFrom] = []Cave{caveTo}
		}
		// Don't add reverse path for start or end
		if parts[0] != "start" && parts[1] != "end" {
			_, caveToExists := graph.caves[caveTo]
			if caveToExists {
				graph.caves[caveTo] = append(graph.caves[caveTo], caveFrom)
			} else {
				graph.caves[caveTo] = []Cave{caveFrom}
			}
		}
	}
	return graph
}

func (graph *Graph) FindPath(path []Cave, smallCaves mapset.Set, fromCave Cave) {
	path = append(path, fromCave)
	if fromCave.IsSmall() {
		smallCaves.Add(fromCave)
	}
	if fromCave.name == "end" {
		fmt.Println(path)
		_path := make([]Cave, len(path))
		copy(_path, path)
		graph.paths = append(graph.paths, _path)
	} else {
		for _, toCave := range graph.caves[fromCave] {
			if !smallCaves.Contains(toCave) {
				graph.FindPath(path, smallCaves.Clone(), toCave)
			}
		}
	}
}

func (graph *Graph) FindPaths() {
	smallCaves := mapset.NewSet()
	var path []Cave
	graph.FindPath(path, smallCaves, Cave{"start"})
	for _, p := range graph.paths {
		var elems []string
		for _, c := range p {
			elems = append(elems, c.name)
		}
		fmt.Println(strings.Join(elems, ","))
	}
}

func Day12() {
	graph := parseFile()
	fmt.Println(graph)
	graph.FindPaths()
	fmt.Println(len(graph.paths))
}
