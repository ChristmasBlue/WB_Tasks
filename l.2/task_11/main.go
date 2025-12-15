package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	var arr = []string{
		"пятак", "пятка", "тяпка", "листок", "слиток", "столик", "стол",
	}
	angrams := Angrams(arr)
	for word, words := range angrams {
		fmt.Printf("%s:%v \n", word, words)
	}
}

func Angrams(arr []string) map[string][]string {
	angrams := make(map[string]map[string]bool)
	firstWord := make([]string, 0)
	for _, str := range arr {
		str = strings.ToLower(str)
		word := cutSumb(str)
		if _, ok := angrams[word]; !ok {
			firstWord = append(firstWord, str)
			angrams[word] = make(map[string]bool)
		}
		angrams[word][str] = true

	}
	angramsWords := make(map[string][]string)
	for _, word := range firstWord {
		wordRunes := string(cutSumb(word))
		if len(angrams[wordRunes]) > 1 {
			for str, _ := range angrams[wordRunes] {
				angramsWords[word] = append(angramsWords[word], str)
			}
			sort.Slice(angramsWords[word], func(i, j int) bool {
				return angramsWords[word][i] < angramsWords[word][j]
			})
		}
	}
	return angramsWords
}

// cutSumb разбирает слово на срез рун и сортирует его
func cutSumb(str string) string {
	arr := make([]rune, 0)
	for _, sumb := range str {
		arr = append(arr, sumb)
	}
	sort.Slice(arr, func(i, j int) bool {
		return arr[i] < arr[j]
	})
	return string(arr)
}
