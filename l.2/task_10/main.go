package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

func main() {
	lines, err := ReadFile("txt.txt")
	if err != nil {
		log.Printf("%v", err)
		return
	}
	lines = SortLines(lines)
	for _, i := range lines {
		fmt.Println(i)
	}
}

func SortLines(lines []string) []string {
	sort.Strings(lines)
	return lines
}

func ReadFile(name string) ([]string, error) {
	file, err := os.Open(name)
	if err != nil {
		log.Printf("error open file, file name: %s\n", name)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
