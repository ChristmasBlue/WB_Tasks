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
	str := scanner.Text()
	j := len(str)
	for i := len(str) - 1; i >= 0; i-- {
		if str[i] != ' ' {
			continue
		}
		fmt.Print(str[i+1:j], " ")
		j = i
	}
	fmt.Print(str[:j])
}
