package main

import (
	"fmt"
	"math"
)

func main() {
	var bit int
	var number, s, value int64
	fmt.Print("Введите число: ")
	fmt.Scanln(&number)
	for {
		fmt.Print("Введите номер бита который нужный изменить(нумерация битов начинается с единицы): ")
		fmt.Scanln(&bit)
		if bit > 0 && bit < 64 {
			s = int64(math.Pow(2.00, float64(bit-1)))
			break
		}
	}
	value = number &^ s
	if value == number {
		value = value | s
	}
	fmt.Printf("Число %d = %b \nНеобходимо изменить %d бит. \nРезультат: %d = %b", number, number, bit, value, value)
}
