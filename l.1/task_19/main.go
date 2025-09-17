package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите строку для переворота: ")
	scanner.Scan()
	str := []rune(scanner.Text())
	for i := len(str) - 1; i >= 0; i-- {
		fmt.Printf("%c", str[i])
	}
}
