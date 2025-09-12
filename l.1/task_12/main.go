package main

import (
	"fmt"
)

func main() {
	fmt.Println("Введите слова, для окончания ввода введите пустую строку")
	var arr []string
	for {
		var str string
		fmt.Print("Введите слово: ")
		fmt.Scanln(&str)
		if str == "" {
			fmt.Println("Ввод слов окончен")
			break
		}
		arr = append(arr, str)
	}
	cache := make(map[string]bool)
	fmt.Print("Уникальные слова:")
	for _, key := range arr {
		if cache[key] != true {
			fmt.Print(" ", key)
			cache[key] = true
		}
	}
}
