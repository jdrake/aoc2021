package day21

import (
	"fmt"
	"math"
)

func Generator() chan int {
	c := make(chan int)
	go func() {
		i := 1
		for {
			c <- i
			i += 1
			if i > 100 {
				i = 1
			}
		}
	}()
	return c
}

func Main() {
	p1, p2 := 4, 8
	// p1, p2 := 5, 8
	score1, score2 := 0, 0
	turn := 0
	c := Generator()
	for {
		roll1 := <-c
		roll2 := <-c
		roll3 := <-c
		sum := roll1 + roll2 + roll3
		// fmt.Println("rolls", roll1, roll2, roll3)
		if math.Mod(float64(turn), 2) == 1 {
			p0 := p2
			newPosition := p2 + sum
			if newPosition > 10 {
				newPosition = int(math.Mod(float64(newPosition), 10))
				if newPosition == 0 {
					newPosition = 10
				}
			}
			p2 = newPosition
			score2 += p2
			fmt.Println("player 2", "rolls", roll1, roll2, roll3, "=", sum, "position", p0, "->", p2, score2)
		} else {
			p0 := p1
			newPosition := p1 + sum
			if newPosition > 10 {
				newPosition = int(math.Mod(float64(newPosition), 10))
				if newPosition == 0 {
					newPosition = 10
				}
			}
			p1 = newPosition
			score1 += p1
			fmt.Println("player 1", "rolls", roll1, roll2, roll3, "=", sum, "position", p0, "->", p1, score1)
		}
		turn += 1
		if score1 >= 1000 || score2 >= 1000 {
			break
		}
	}

	var losingScore int
	if score1 > score2 {
		losingScore = score2
	} else {
		losingScore = score1
	}
	fmt.Println(turn * 3 * losingScore)
}
