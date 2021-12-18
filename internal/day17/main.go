package day17

import "fmt"

func MaxHeight(xMin int, xMax int, yMin int, yMax int, vx int, vy int) bool {
	vx0, vy0 := vx, vy
	x, y := 0, 0
	i := 0
	hMax := 0
	for x <= xMax && y >= yMin && i < 1000 {
		x += vx
		y += vy
		// fmt.Println(x, y)

		if y > hMax {
			hMax = y
		}

		if x >= xMin && x <= xMax && y >= yMin && y <= yMax {
			fmt.Println("landed", vx0, vy0)
			return true
		}

		if vx > 0 {
			vx -= 1
		} else if vx < 0 {
			vx += 1
		}
		vy -= 1

		i += 1
	}
	// fmt.Println("missed", vx0, vy0)
	return false
}

func Main() {
	// xMin, xMax, yMin, yMax := 20, 30, -10, -5
	xMin, xMax, yMin, yMax := 138, 184, -125, -71
	count := 0
	for vx := 1; vx < 1000; vx++ {
		for vy := -1000; vy < 1000; vy++ {
			if landed := MaxHeight(xMin, xMax, yMin, yMax, vx, vy); landed {
				count += 1
			}
		}
	}
	fmt.Println("---")
	fmt.Println(count)
}
