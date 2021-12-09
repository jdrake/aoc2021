package day8

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	mapset "github.com/deckarep/golang-set"
)

func parseLine(line string) ([]string, []string) {
	var signals []string
	var outputs []string
	foundDelimiter := false
	for _, s := range strings.Fields(line) {
		if s == "|" {
			foundDelimiter = true
		} else if foundDelimiter {
			outputs = append(outputs, s)
		} else {
			signals = append(signals, s)
		}
	}
	return signals, outputs
}

type Entry struct {
	signals []string
	outputs []string
}

func (e Entry) GetOutputDigitFrequencies() []int {
	counts := make([]int, 10)
	for _, output := range e.outputs {
		if len(output) == 2 {
			counts[1] += 1
		} else if len(output) == 3 {
			counts[7] += 1
		} else if len(output) == 4 {
			counts[4] += 1
		} else if len(output) == 7 {
			counts[8] += 1
		}
	}
	return counts
}

type Signal struct {
	encoded string
	value   int
	charSet mapset.Set
}

func MakeCharSet(s string) mapset.Set {
	chars := make([]interface{}, len(s))
	for _, s := range strings.Split(s, "") {
		if s != "" {
			chars = append(chars, s)
		}
	}
	return mapset.NewSetFromSlice(chars)
}

func (e Entry) DecipherSignals() []Signal {
	signals := make([]Signal, 10)
	for _, encoded := range e.signals {
		if len(encoded) == 2 {
			signals[1] = Signal{
				encoded: encoded,
				value:   1,
				charSet: MakeCharSet(encoded),
			}
		} else if len(encoded) == 3 {
			signals[7] = Signal{
				encoded: encoded,
				value:   7,
				charSet: MakeCharSet(encoded),
			}
		} else if len(encoded) == 4 {
			signals[4] = Signal{
				encoded: encoded,
				value:   4,
				charSet: MakeCharSet(encoded),
			}
		} else if len(encoded) == 7 {
			signals[8] = Signal{
				encoded: encoded,
				value:   8,
				charSet: MakeCharSet(encoded),
			}
		}
	}

	// fmt.Println(signals)

	// Handle 6-segment signals
	for _, encoded := range e.signals {
		charSet := MakeCharSet(encoded)
		if len(encoded) == 6 {
			diff := signals[4].charSet.Difference(charSet).ToSlice()
			// fmt.Println(signals[4].charSet, "diff", encoded, "=", diff)
			if len(diff) == 0 {
				// fmt.Println("found 9")
				signals[9] = Signal{
					encoded: encoded,
					value:   9,
					charSet: MakeCharSet(encoded),
				}
			} else if len(diff) == 1 {
				if signals[1].charSet.Contains(diff[0]) {
					// fmt.Println("found 6")
					signals[6] = Signal{
						encoded: encoded,
						value:   6,
						charSet: MakeCharSet(encoded),
					}
				} else {
					// fmt.Println("found 0")
					signals[0] = Signal{
						encoded: encoded,
						value:   0,
						charSet: MakeCharSet(encoded),
					}
				}
			}
		}
	}

	// fmt.Println(signals)

	// Handle 5-segment signals
	for _, encoded := range e.signals {
		charSet := MakeCharSet(encoded)
		if len(encoded) == 5 {
			diff := signals[1].charSet.Difference(charSet).ToSlice()
			// fmt.Println(signals[1].charSet, "diff", encoded, "=", diff)
			if len(diff) == 0 {
				// fmt.Println("found 3")
				signals[3] = Signal{
					encoded: encoded,
					value:   3,
					charSet: MakeCharSet(encoded),
				}
			} else {
				diff = charSet.Difference(signals[9].charSet).ToSlice()
				// fmt.Println(charSet, "diff", signals[9].charSet, "=", diff)
				if len(diff) == 0 {
					// fmt.Println("found 5")
					signals[5] = Signal{
						encoded: encoded,
						value:   5,
						charSet: MakeCharSet(encoded),
					}
				} else {
					// fmt.Println("found 2")
					signals[2] = Signal{
						encoded: encoded,
						value:   2,
						charSet: MakeCharSet(encoded),
					}
				}
			}
		}
	}

	// fmt.Println(signals)

	return signals
}

func SortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}

func (e Entry) TranslateOutput(signals []Signal) int {
	// fmt.Println("--------------------")
	sum := 0
	signalMap := make(map[string]int)
	for _, signal := range signals {
		signalMap[SortString(signal.encoded)] = signal.value
	}
	fmt.Println(signalMap)
	for i := 0; i < len(e.outputs); i++ {
		s := SortString(e.outputs[i])
		value := signalMap[s]
		sum += value * int(math.Pow10(len(e.outputs)-i-1))
		// fmt.Println(e.outputs[i], s, value, int(math.Pow10(i)))
	}
	fmt.Println(e.outputs, sum)
	return sum
}

func parseFile() []Entry {
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
	var entries []Entry
	for scanner.Scan() {
		line := scanner.Text()
		signals, outputs := parseLine(line)
		entries = append(entries, Entry{signals, outputs})
	}
	return entries
}

func Sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

func Day8() {
	entries := parseFile()
	sum := 0
	for _, entry := range entries {
		signals := entry.DecipherSignals()
		sum += entry.TranslateOutput(signals)
		// break
	}
	fmt.Println(sum)
}
