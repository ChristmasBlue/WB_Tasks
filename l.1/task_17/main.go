package main

import (
	"fmt"
)

func main() {
	arr := []int{1, 3, 5, 7, 9, 11, 13, 15}
	fmt.Println(binarySearch(arr, 0, len(arr)-1, 3))
	fmt.Println(binarySearch(arr, 0, len(arr)-1, 13))
	fmt.Println(binarySearch(arr, 0, len(arr)-1, 6))
	fmt.Println(binarySearch(arr, 0, len(arr)-1, 16))
}

func binarySearch(arr []int, begin, end, trigger int) int {
	if begin > end {
		return -1
	}
	x := begin + (end-begin)/2

	switch {
	case arr[x] == trigger:
		return x
	case arr[x] < trigger:
		return binarySearch(arr, x+1, end, trigger)
	default: // arr[x]>trigger:
		return binarySearch(arr, begin, x-1, trigger)
	}
}
