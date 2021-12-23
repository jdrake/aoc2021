package day20

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

type Point struct {
	x, y int
}

type Algo []int

type Image map[Point]int

func (image Image) MinMax() (int, int, int, int) {
	minx, maxx, miny, maxy := math.MaxInt, -math.MaxInt, math.MaxInt, -math.MaxInt
	for point := range image {
		if point.x < minx {
			minx = point.x
		}
		if point.x > maxx {
			maxx = point.x
		}
		if point.y < miny {
			miny = point.y
		}
		if point.y > maxy {
			maxy = point.y
		}
	}
	return minx, maxx, miny, maxy
}

func (image Image) String() string {
	minx, maxx, miny, maxy := image.MinMax()
	s := ""
	for y := miny; y <= maxy; y++ {
		for x := minx; x <= maxx; x++ {
			if image[Point{x, y}] == 1 {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

func (image Image) LitCount() int {
	count := 0
	for _, bit := range image {
		if bit == 1 {
			count += 1
		}
	}
	return count
}

func Bin2Int(s string) int {
	i, _ := strconv.ParseInt(s, 2, 64)
	return int(i)
}

var buffer int = 0

func (image Image) Apply(algo Algo) Image {
	newImage := make(Image)
	minx, maxx, miny, maxy := image.MinMax()
	for x := minx - buffer; x <= maxx+buffer; x++ {
		for y := miny - buffer; y <= maxy+buffer; y++ {
			point := Point{x, y}
			var points = []Point{
				{point.x - 1, point.y - 1},
				{point.x, point.y - 1},
				{point.x + 1, point.y - 1},
				{point.x - 1, point.y},
				point,
				{point.x + 1, point.y},
				{point.x - 1, point.y + 1},
				{point.x, point.y + 1},
				{point.x + 1, point.y + 1},
			}
			s := ""
			for _, p := range points {
				s += strconv.Itoa(image[p])
			}
			num := Bin2Int(s)
			bit := algo[num]
			newImage[point] = bit
			// fmt.Println(point, s, num, bit)
		}
	}
	fmt.Println(newImage)
	return newImage
}

func parseFile(name string) (Algo, Image) {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/" + name + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	scanner.Scan()
	var algo Algo
	for _, char := range strings.Split(scanner.Text(), "") {
		bit := 0
		if char == "#" {
			bit = 1
		}
		algo = append(algo, bit)
	}
	scanner.Scan()
	row := 0
	image := make(Image)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		} else {
			for col, char := range strings.Split(line, "") {
				point := Point{col, row}
				bit := 0
				if char == "#" {
					bit = 1
				}
				image[point] = bit
			}
		}
		row += 1
	}
	return algo, image
}

func Main() {
	algo, image := parseFile("input")
	fmt.Println(algo)
	fmt.Println(image)
	image = image.Apply(algo)
	image = image.Apply(algo)
	fmt.Println()
	fmt.Println(image.LitCount())
}
