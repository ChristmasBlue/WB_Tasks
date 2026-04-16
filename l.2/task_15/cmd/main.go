package main

import (
	"bufio"
	"fmt"
	"os"
	"task_15/internal/domain"
	"task_15/internal/service"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		strs, err := service.ParseCommand(scanner.Text())
		if err != nil {
			fmt.Println(err)
			continue
		}

		command := domain.NewConditionals()

		command.ParseCommands(strs)

		command.Execute()

	}
}
