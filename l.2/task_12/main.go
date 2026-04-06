package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Config struct {
	contextA int
	contextB int
	contextC int
	nFlag    bool
	cFlag    bool
	FFlag    bool
	iFlag    bool
	vFlag    bool
	countStr int
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
	flag.BoolVar(&cfg.cFlag, "c", false, "выводит только количество строк, совпадающих с шаблоном")
	flag.BoolVar(&cfg.nFlag, "n", false, "выводит номер строки передкаждой найденной строкой")
	flag.BoolVar(&cfg.FFlag, "F", false, "воспринимает шаблон как фиксированную строку, а не регулярное выражение")
	flag.BoolVar(&cfg.iFlag, "i", false, "игнорирует регистр")
	flag.BoolVar(&cfg.vFlag, "v", false, "инвертирует фильтр: выводить строки не содержащие шаблон")

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

	cfg.countStr = 0

	switch {
	case cfg.cFlag:

		count := countMatches(reader, cfg)
		fmt.Println(count)

	case cfg.contextC != 0:
		ContextC(reader, cfg)

	case cfg.contextA != 0 && cfg.contextB != 0:

		ContextBoth(reader, cfg)

	case cfg.contextA != 0:

		ContextA(reader, cfg)

	case cfg.contextB != 0:

		ContextB(reader, cfg)

	default:
		SimpleGrep(reader, cfg)
	}
}

func ContextA(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	separator := false
	lastSeparator := false
	outStr := 0 // счётчик вывода строк после найденного совпадения

	for scanner.Scan() {
		line := scanner.Text()

		if cfg.Matches(line) {
			outStr = cfg.contextA
			cfg.countStr++
			cfg.Print(line)
			separator = true
			continue
		}

		if outStr > 0 {
			cfg.countStr++
			cfg.Print(line)
			lastSeparator = false
			outStr--
			continue
		}

		if outStr == 0 && separator {
			fmt.Println("--------")
			separator = false
			lastSeparator = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return
	}

	if !lastSeparator {
		fmt.Println("--------")
	}
}

func ContextB(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	separator := false
	prevStr := make([]string, 0, cfg.contextB)

	for scanner.Scan() {

		line := scanner.Text()

		if cfg.Matches(line) {
			if len(prevStr) == cap(prevStr) && separator {
				fmt.Println("--------")
			}
			for _, str := range prevStr {
				cfg.countStr++
				cfg.Print(str)
			}
			cfg.countStr++
			cfg.Print(line)
			prevStr = prevStr[:0]
			separator = true
			continue
		}

		if len(prevStr) == cfg.contextB {
			for i := 1; i < len(prevStr); i++ {
				prevStr[i-1] = prevStr[i]
			}
			prevStr[len(prevStr)-1] = line
		} else {
			prevStr = append(prevStr, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return
	}

	fmt.Println("--------")
}

func ContextC(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	prevStr := make([]string, 0, cfg.contextC)
	separator := false
	lastSeparator := false
	outStr := 0 //счётчик вывода строк после найденного совпадения

	for scanner.Scan() {
		line := scanner.Text()

		if cfg.Matches(line) {
			for _, str := range prevStr {
				cfg.countStr++
				cfg.Print(str)
			}
			cfg.countStr++
			cfg.Print(line)
			lastSeparator = false
			prevStr = prevStr[:0]
			outStr = cfg.contextC
			separator = true
			continue
		}

		if outStr > 0 {
			cfg.countStr++
			cfg.Print(line)
			outStr--
			continue
		}

		if outStr == 0 && separator {
			fmt.Println("--------")
			separator = false
			lastSeparator = true
		}

		if len(prevStr) == cfg.contextC {
			for i := 1; i < len(prevStr); i++ {
				prevStr[i-1] = prevStr[i]
			}
			prevStr[len(prevStr)-1] = line
		} else {
			prevStr = append(prevStr, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return
	}

	if !lastSeparator {
		fmt.Println("--------")
	}

}

func ContextBoth(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	prevStr := make([]string, 0, cfg.contextB)
	separator := false
	lastSeparator := false
	outStr := 0

	for scanner.Scan() {
		line := scanner.Text()

		if cfg.Matches(line) {
			for _, str := range prevStr {
				cfg.countStr++
				cfg.Print(str)
			}
			cfg.countStr++
			cfg.Print(line)
			lastSeparator = false
			prevStr = prevStr[:0]
			outStr = cfg.contextA
			separator = true
			continue
		}

		if outStr > 0 {
			cfg.countStr++
			cfg.Print(line)
			lastSeparator = false
			outStr--
			continue
		}

		if outStr == 0 && separator {
			fmt.Println("--------")
			separator = false
			lastSeparator = true
		}

		if len(prevStr) == cfg.contextB {
			for i := 1; i < len(prevStr); i++ {
				prevStr[i-1] = prevStr[i]
			}
			prevStr[len(prevStr)-1] = line
		} else {
			prevStr = append(prevStr, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return
	}

	if !lastSeparator {
		fmt.Println("--------")
	}
}

func (cfg *Config) Matches(str string) bool {

	var ok bool
	templateStr := cfg.template

	if cfg.iFlag {
		templateStr = strings.ToLower(cfg.template)
		str = strings.ToLower(str)
	}

	if cfg.FFlag {
		ok = strings.Contains(str, templateStr)
	} else {
		re, err := regexp.Compile(templateStr)
		if err != nil {
			ok = strings.Contains(str, templateStr)
		} else {
			ok = re.MatchString(str)
		}
	}

	if cfg.vFlag {
		return !ok
	}

	return ok
}

func (cfg *Config) Print(str string) {
	if cfg.nFlag {
		fmt.Printf("%d. %s\n", cfg.countStr, str)
		return
	}

	fmt.Println(str)
}

func SimpleGrep(reader *os.File, cfg *Config) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.Matches(line) {
			cfg.countStr++
			cfg.Print(line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return
	}
}

func countMatches(reader *os.File, cfg *Config) int {
	scanner := bufio.NewScanner(reader)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.Matches(line) {
			count++
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Ошибка чтения: %v", err)
		return 0
	}

	return count
}
