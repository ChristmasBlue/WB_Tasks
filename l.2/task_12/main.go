package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	contextA int
	contextB int
	contextC int
	template string   //шаблон для поиска
	files    []string //файлы в которых нужно совершать поиск
}

func main() {
	cfg := parseFlags()

	if cfg.template == "" {
		log.Println("Ошибка, необходимо указать шаблон")
		flag.Usage()
		os.Exit(1)
	}

	if len(cfg.files) > 0 {
		for _, filename := range cfg.files {
			processFile(filename, &cfg)
		}
	} else {
		processRead(os.Stdin, &cfg)
	}

}

// обработка входных данных
func parseFlags() Config {
	var cfg Config

	flag.IntVar(&cfg.contextA, "A", 0, "после каждой найденной строки дополнительно вывести N строк после неё.")
	flag.IntVar(&cfg.contextB, "B", 0, "вывести N строк до каждой найденной строки.")
	flag.IntVar(&cfg.contextC, "C", 0, "вывести N строк контекста вокруг найденной строки.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: %s [опции] шаблон [файл]...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Если файл не указан, читает из STDIN.")
	}

	flag.Parse()

	args := flag.Args()

	//если нет аргументов, возвращаем пустой конфиг (ошибка)
	if len(args) == 0 {
		log.Println("Нет аргументов.")
		return cfg
	}

	cfg.template = args[0]

	if len(args) > 1 {
		cfg.files = args[1:]
	}

	return cfg
}

// обработкак входных данных через файл
func processFile(filename string, cfg *Config) {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Ошибка открытия файла: %s, %v", filename, err)
		return
	}

	defer file.Close()
	processRead(file, cfg)
}

// обработка данных чтением
func processRead(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)

	//lineNum := 1

	if cfg.contextB != 0 {
		strArr := make([]string, cfg.contextB)

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, cfg.template) {
				for _, str := range strArr {
					if str != "" {
						fmt.Println(str)
					}
				}
				fmt.Println(line)
				fmt.Println("--------")
			}

			for i := 1; i < len(strArr); i++ {
				strArr[i-1] = strArr[i]
			}
			strArr[len(strArr)-1] = line
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка чтения: %v", err)
		}
	} else if cfg.contextA != 0 {
		lineNum := 0
		strArr := make(map[int][]string)
		countLines := make([]int, 0)

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, cfg.template) {
				strArr[lineNum] = make([]string, 0)
				countLines = append(countLines, lineNum)
				lineNum++
			}

			delEl := make([]int, 0)
			for i := 0; i < len(countLines); i++ {
				strArr[countLines[i]] = append(strArr[countLines[i]], line)
				if len(strArr[countLines[i]]) == cfg.contextA+1 {
					delEl = append(delEl, i)
				}
			}

			for _, i := range delEl {
				countLines = append(countLines[:i], countLines[i+1:]...)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка чтения: %v", err)
		}

		for i := 0; i < len(strArr); i++ {
			for _, str := range strArr[i] {
				if str != "" {
					fmt.Println(str)
				}
			}
			fmt.Println("--------")
		}
	} else {
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, cfg.template) {
				fmt.Println(line)
			}
			//lineNum++
		}

		if err := scanner.Err(); err != nil {
			log.Printf("Ошибка чтения: %v", err)
		}
	}
}

func ContextA(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	lineNum := 0
	strArr := make(map[int][]string)
	countLines := make([]int, 0)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, cfg.template) {
			strArr[lineNum] = make([]string, 0)
			countLines = append(countLines, lineNum)
			lineNum++
		}

		delEl := make([]int, 0)
		for i := 0; i < len(countLines); i++ {
			strArr[countLines[i]] = append(strArr[countLines[i]], line)
			if len(strArr[countLines[i]]) == cfg.contextA+1 {
				delEl = append(delEl, i)
			}
		}

		for _, i := range delEl {
			countLines = append(countLines[:i], countLines[i+1:]...)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
	}

	for i := 0; i < len(strArr); i++ {
		for _, str := range strArr[i] {
			if str != "" {
				fmt.Println(str)
			}
		}
		fmt.Println("--------")
	}
}

func ContextB(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)

	strArr := make([]string, cfg.contextB)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, cfg.template) {
			for _, str := range strArr {
				if str != "" {
					fmt.Println(str)
				}
			}
			fmt.Println(line)
			fmt.Println("--------")
		}

		for i := 1; i < len(strArr); i++ {
			strArr[i-1] = strArr[i]
		}
		strArr[len(strArr)-1] = line
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
	}
}
