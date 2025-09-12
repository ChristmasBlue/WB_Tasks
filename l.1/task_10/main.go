package main

import (
	"fmt"
	"strconv"
)

func main() {
	var arr = []float32{-25.4, -27.0, 13.0, 19.0, 15.5, 24.5, -21.0, 32.5, 0.5, -0.5, 0.0, 5.5, -5.5, 110.5, -105.8}
	cache := make(map[string][]float32)
	var num int
	for i := 0; i < len(arr); i++ {
		x := arr[i]
		switch {
		case ((x >= 0) && (x < 10)):
			cache["+0"] = append(cache["+0"], x)
		case ((x < 0) && (x > -10)):
			cache["-0"] = append(cache["-0"], x)
		default:
			num = int(x)
			v := 0
			num = num / 10
			v = num
			key := strconv.Itoa(v * 10)
			cache[key] = append(cache[key], x)
		}
	}
	for x, y := range cache {
		fmt.Printf("%s: %.1f\n", x, y)
	}
}
