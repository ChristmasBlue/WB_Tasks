package main

import (
	"bufio"
	"fmt"
	"os"
	"task_26/str"
)

func main() {
	fmt.Print("Введите строку для проверки уникальности символов: ")
	scaner := bufio.NewScanner(os.Stdin)
	scaner.Scan()
	line := str.NewStr(scaner.Text())
	if line.CheckUniqSymb() {
		fmt.Println("true")
	} else {
		fmt.Println("false")
	}
}
