package day2

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

func Day2() {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/day2.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var horizontal int64 = 0
	var depth int64 = 0
	var aim int64 = 0
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		v, err := strconv.ParseInt(parts[1], 10, 0)
		if err != nil {
			log.Fatal(err)
		}
		if parts[0] == "forward" {
			horizontal += v
			depth += aim * v
		} else if parts[0] == "up" {
			aim -= v
		} else if parts[0] == "down" {
			aim += v
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	result := horizontal * depth
	fmt.Println(result)
}
