package main

import (
	"fmt"
)

func main() {
	x := 5
	y := 1
	fmt.Printf("x = %d, y = %d\n", x, y)
	x = x + y
	y = x - y
	x = x - y
	fmt.Printf("x = %d, y = %d\n", x, y)
	x = 6
	y = 2
	fmt.Printf("x = %d, y = %d\n", x, y)
	x, y = y, x
	fmt.Printf("x = %d, y = %d\n", x, y)
	x = 7
	y = 3
	fmt.Printf("x = %d, y = %d\n", x, y)
	x = x ^ y
	y = x ^ y
	x = x ^ y
	fmt.Printf("x = %d, y = %d\n", x, y)
}
