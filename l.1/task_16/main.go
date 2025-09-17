package main

import (
	"fmt"
)

func main() {
	fmt.Println("Введите числа через Enter (Ctrl+D или нечисловой ввод для завершения):")
	arr := make([]int, 0)
	var x int
	for {
		_, err := fmt.Scan(&x)
		if err != nil {
			break
		}
		arr = append(arr, x)
	}
	newArr := quickSort(arr)
	fmt.Println("Не отсортированный массив: ", arr)
	fmt.Println("Отсортированный массив: ", newArr)
}

func quickSort(arr []int) []int {
	sorted := make([]int, len(arr))
	copy(sorted, arr)
	sort(sorted, 0, len(sorted)-1)
	return sorted
}

func sort(arr []int, low, high int) {
	if low >= high {
		return
	}
	j := low
	for i := low; i < high; i++ {
		if arr[i] < arr[high] {
			swap(arr, i, j)
			j++
		}
	}
	swap(arr, j, high)
	sort(arr, low, j-1)
	sort(arr, j+1, high)

}

func swap(arr []int, i, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}
