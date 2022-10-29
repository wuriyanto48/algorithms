// Wuriyanto 2022
package main

import "fmt"

func main() {
	res := wordDistance("kitten", "sitting")
	fmt.Println(res)
	if res < 2 {
		fmt.Println("text is similar")
	} else {
		fmt.Println("text is not similar")
	}

}

func wordDistance(a, b string) int {
	aLen := len(a) + 1
	bLen := len(b) + 1

	d := make([][]int, aLen)
	for i := 0; i < aLen; i++ {
		d[i] = make([]int, bLen)
	}

	for i := 0; i <= aLen-1; i++ {

		d[i][0] = i
	}

	for j := 0; j <= bLen-1; j++ {
		d[0][j] = j
	}

	min := func(a, b, c int) int {
		if b < a {
			a = b
		}

		if c < a {
			return c
		}

		return a
	}

	var i int
	var j int
	for i = 0; i < aLen-1; i++ {
		for j = 0; j < bLen-1; j++ {

			if a[i] == b[j] {
				d[i+1][j+1] = d[(i+1)-1][(j+1)-1]
			} else {

				d[i+1][j+1] = min(d[(i+1)-1][j+1],
					d[i+1][(j+1)-1],
					d[(i+1)-1][(j+1)-1]) + 1
			}
		}

	}

	for x := 0; x < aLen; x++ {
		for y := 0; y < bLen; y++ {
			fmt.Printf("%d ", d[x][y])
		}
		fmt.Println()
	}

	return d[aLen-1][bLen-1]
}
