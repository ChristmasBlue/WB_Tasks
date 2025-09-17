package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var index int
	scaner := bufio.NewScanner(os.Stdin)
	fmt.Println("Введите текст для изменения: ")
	scaner.Scan()
	str := []rune(scaner.Text())
	lenStr := len(str)
	//использую цикл для валидации ввода индекса для удаления из среза,
	//пока ввод не окажется валидным, выхода из цикла не произойдёт
	for {
		fmt.Printf("Количество элементов в срезе %d, введите номер элемента для удаления: ", lenStr)
		fmt.Scanln(&index)
		if index <= lenStr && index > 0 {
			break
		}
		fmt.Println("Некорректный индекс элемента для удаления.\nПовторите попытку")
		fmt.Println()
	}
	//исльзовал функцию append для удаления элемента
	str = append(str[:index-1], str[index:]...)
	//также закомментировал второй вариант через функцию copy и обрезание хвоста
	//copy(str[index-1:], str[index:])
	//str = str[:lenStr-1]
	fmt.Println("Исходный текст с удалённым элементом:")
	for _, i := range str {
		fmt.Printf("%c", i)
	}
}
