package main

import (
	"fmt"
	"strconv"
)

func main() {
	var arr1 []int
	var arr2 []int
	fmt.Println("Введите числа, для окончания ввода введите пустую строку")
	for {
		var x string
		fmt.Print("Введите число в первый срез: ")
		fmt.Scanln(&x)
		if x == "" {
			break
		}
		str, err := strconv.Atoi(x)
		if err == nil {
			arr1 = append(arr1, str)
		}
	}
	for {
		var x string
		fmt.Print("Введите число во второй срез: ")
		fmt.Scanln(&x)
		if x == "" {
			break
		}
		str, err := strconv.Atoi(x)
		if err == nil {
			arr2 = append(arr2, str)
		}
	}
	fmt.Println("Ввод чисел окончен")
	cache1 := make(map[int]bool)
	cache2 := make(map[int]bool)
	for _, key := range arr1 {
		cache1[key] = true
	}
	for _, key := range arr2 {
		if cache1[key] {
			cache2[key] = true
		}
	}
	fmt.Print("Пересечение множеств:")
	for key := range cache2 {
		fmt.Print(" ", key)
	}
}
