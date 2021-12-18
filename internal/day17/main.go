package day17

import "fmt"

func MaxHeight(xMin int, xMax int, yMin int, yMax int, vx int, vy int) int {
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
			fmt.Println("landed", x, y, "hMax", hMax)
			return hMax
		}

		if vx > 0 {
			vx -= 1
		} else if vx < 0 {
			vx += 1
		}
		vy -= 1

		i += 1
	}
	return 0
}

func Main() {
	xMin, xMax, yMin, yMax := 138, 184, -125, -71
	hMax := 0
	for vx := 1; vx < 1000; vx++ {
		for vy := 1; vy < 1000; vy++ {
			if v := MaxHeight(xMin, xMax, yMin, yMax, vx, vy); v > hMax {
				hMax = v
			}
		}
	}
	fmt.Println("---")
	fmt.Println(hMax)
}
