package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"task_24/point"
)

func main() {
	var x, y float64
	scaner := bufio.NewScanner(os.Stdin)

	//использую бесконечный цикл для валидации ввода,
	//пока ввод не будет корректным выход из цикла не произойдёт
	for {
		fmt.Print("Введите координаты первой точки:\nX = ")
		scaner.Scan()
		val, err := strconv.ParseFloat(strings.TrimSpace(scaner.Text()), 64)
		if err != nil {
			fmt.Println("Введены некорректные координаты. Повторите попытку")
			continue
		}
		x = val
		fmt.Print("Y = ")
		scaner.Scan()
		val, err = strconv.ParseFloat(strings.TrimSpace(scaner.Text()), 64)
		if err != nil {
			fmt.Println("Введены некорректные координаты. Повторите попытку")
			continue
		}
		y = val
		break
	}
	point1 := point.NewPoint(x, y)
	for {
		fmt.Print("Введите координаты второй точки:\nX = ")
		scaner.Scan()
		val, err := strconv.ParseFloat(strings.TrimSpace(scaner.Text()), 64)
		if err != nil {
			fmt.Println("Введены некорректные координаты. Повторите попытку")
			continue
		}
		x = val
		fmt.Print("Y = ")
		scaner.Scan()
		val, err = strconv.ParseFloat(strings.TrimSpace(scaner.Text()), 64)
		if err != nil {
			fmt.Println("Введены некорректные координаты. Повторите попытку")
			continue
		}
		y = val
		break
	}
	point2 := point.NewPoint(x, y)
	fmt.Println("Расстояние между точками равно: ", point1.Distance(point2))
}
