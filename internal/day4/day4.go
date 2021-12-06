package main

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

func parseInputLine(line string) []int {
	parts := strings.Split(line, ",")
	var digits []int
	for _, v := range parts {
		if v != "" {
			digit, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			digits = append(digits, digit)
		}
	}
	return digits
}

func parseGridLine(line string) []int {
	parts := strings.Split(line, " ")
	var digits []int
	for _, v := range parts {
		if v != "" {
			digit, err := strconv.Atoi(v)
			if err != nil {
				log.Fatal(err)
			}
			digits = append(digits, digit)
		}
	}
	return digits
}

type Board struct {
	numbers   [][]int
	numberMap map[int][2]int
	grid      [][]int
	won       bool
}

func (b Board) Mark(number int) bool {
	location, ok := b.numberMap[number]
	if ok {
		x, y := location[0], location[1]
		b.grid[x][y] = 1
		return b.check(x, y)
	}
	return false
}

func (b Board) check(x int, y int) bool {
	rowFull, colFull := true, true
	for i := 0; i < 5; i++ {
		// Check the row
		if b.grid[x][i] == 0 {
			rowFull = false
		}
		// Check the column
		if b.grid[i][y] == 0 {
			colFull = false
		}
	}
	return rowFull || colFull
}

func (b Board) score(input int) int {
	sum := 0
	fmt.Println(b.numbers)
	fmt.Println(b.grid)
	for i, row := range b.grid {
		for j, v := range row {
			if v != 1 {
				sum += b.numbers[i][j]
			}
		}
	}
	fmt.Println(sum, input)
	return sum * input
}

func parseFile() ([]int, []Board) {
	_, fileName, _, _ := runtime.Caller(0)
	prefixPath := filepath.Dir(fileName)
	file, err := os.Open(prefixPath + "/day4.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	var input []int
	var board Board
	var boards []Board
	currGridRow := 0
	row := 0
	for scanner.Scan() {
		line := scanner.Text()
		if row == 0 {
			input = parseInputLine(line)
		} else if line == "" {
			board = Board{}
			board.numberMap = make(map[int][2]int)
			for i := 0; i < 5; i++ {
				board.grid = append(board.grid, make([]int, 5))
			}
			currGridRow = 0
		} else {
			numbers := parseGridLine(line)
			board.numbers = append(board.numbers, numbers)
			for column, value := range numbers {
				board.numberMap[value] = [2]int{currGridRow, column}
			}
			currGridRow += 1
			if currGridRow == 5 {
				boards = append(boards, board)
			}
		}
		row += 1

	}
	return input, boards
}

func Day4() {
	input, boards := parseFile()
	for _, input := range input {
		var remainingBoards []Board
		for _, board := range boards {
			board.won = board.Mark(input)
			if !board.won {
				remainingBoards = append(remainingBoards, board)
			}
		}
		fmt.Println("len(remainingBoards)", len(remainingBoards))
		if len(remainingBoards) == 0 {
			score := boards[0].score(input)
			fmt.Println(score)
			break
		}
		boards = remainingBoards
	}
}
