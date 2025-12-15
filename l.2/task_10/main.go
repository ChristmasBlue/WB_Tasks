package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

// Config содержит настройки сортировки
type Config struct {
	column      int    // -k: колонка для сортировки (с 1)
	numeric     bool   // -n: числовая сортировка
	reverse     bool   // -r: обратный порядок
	unique      bool   // -u: только уникальные строки
	month       bool   // -M: сортировка по месяцам
	human       bool   // -h: человекочитаемые числа
	check       bool   // -c: проверить отсортированность
	ignoreBlank bool   // -b: игнорировать хвостовые пробелы
	separator   string // разделитель колонок
}

var months = map[string]int{
	"jan": 1, "feb": 2, "mar": 3, "apr": 4,
	"may": 5, "jun": 6, "jul": 7, "aug": 8,
	"sep": 9, "oct": 10, "nov": 11, "dec": 12,
}

func main() {
	cfg := parseFlags()
	lines := readLines(cfg)

	if cfg.check {
		if isSorted(lines, cfg) {
			fmt.Fprintln(os.Stderr, "sorted")
			os.Exit(0)
		}
		fmt.Fprintln(os.Stderr, "not sorted")
		os.Exit(1)
	}

	lines = sortLines(lines, cfg)
	writeLines(lines)
}

// parseFlags парсит аргументы командной строки
func parseFlags() Config {
	cfg := Config{separator: "\t"}

	flag.IntVar(&cfg.column, "k", 0, "sort by column N")
	flag.BoolVar(&cfg.numeric, "n", false, "numeric sort")
	flag.BoolVar(&cfg.reverse, "r", false, "reverse sort")
	flag.BoolVar(&cfg.unique, "u", false, "unique lines only")
	flag.BoolVar(&cfg.month, "M", false, "month sort")
	flag.BoolVar(&cfg.human, "h", false, "human numeric sort")
	flag.BoolVar(&cfg.check, "c", false, "check if sorted")
	flag.BoolVar(&cfg.ignoreBlank, "b", false, "ignore trailing blanks")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [file]\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nIf file is not specified or is '-', reads from stdin.")
	}

	flag.Parse()

	return cfg
}

// readLines читает строки из файла или stdin
func readLines(cfg Config) []string {
	var scanner *bufio.Scanner
	var file *os.File
	var err error

	// Определяем источник ввода
	if flag.NArg() == 0 || flag.Arg(0) == "-" {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		scanner = bufio.NewScanner(file)
	}

	lines := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if cfg.ignoreBlank {
			line = strings.TrimRight(line, " \t")
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		os.Exit(1)
	}

	return lines
}

// writeLines выводит строки
func writeLines(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

// sortLines сортирует строки согласно конфигурации
func sortLines(lines []string, cfg Config) []string {
	if cfg.unique {
		lines = unique(lines)
	}

	sort.SliceStable(lines, func(i, j int) bool {
		less := compare(lines[i], lines[j], cfg)
		if cfg.reverse {
			return !less
		}
		return less
	})

	return lines
}

// compare сравнивает две строки для сортировки
func compare(a, b string, cfg Config) bool {
	// Получаем ключи для сравнения
	keyA := getKey(a, cfg)
	keyB := getKey(b, cfg)

	// Проверяем специальные типы сортировки
	if cfg.month {
		ma := monthValue(keyA)
		mb := monthValue(keyB)
		if ma != mb {
			return ma < mb
		}
		// Если месяцы равны, продолжаем сравнение
	}

	if cfg.human {
		ha := humanValue(keyA)
		hb := humanValue(keyB)
		if ha != hb {
			return ha < hb
		}
		// Если человекочитаемые числа равны, продолжаем сравнение
	}

	if cfg.numeric {
		na := numValue(keyA)
		nb := numValue(keyB)
		if na != nb {
			return na < nb
		}
		// Если числа равны, продолжаем сравнение
	}

	// Сравниваем как строки
	return keyA < keyB
}

// getKey возвращает ключ для сортировки строки
func getKey(line string, cfg Config) string {
	if cfg.column > 0 {
		parts := strings.Split(line, cfg.separator)
		if cfg.column-1 < len(parts) {
			return parts[cfg.column-1]
		}
		return ""
	}
	return line
}

// monthValue возвращает числовое значение месяца
func monthValue(s string) int {
	s = strings.ToLower(strings.TrimSpace(s))
	if len(s) >= 3 {
		if val, ok := months[s[:3]]; ok {
			return val
		}
	}
	return 0
}

// numValue пытается преобразовать строку в число
func numValue(s string) float64 {
	s = strings.TrimSpace(s)
	// Удаляем возможные пробелы внутри числа
	s = strings.ReplaceAll(s, " ", "")
	val, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return val
}

// humanValue парсит человекочитаемые числа (1K, 2M и т.д.)
func humanValue(s string) float64 {
	s = strings.TrimSpace(s)

	// Ищем числовую часть
	var numStr strings.Builder
	var suffix string

	for _, r := range s {
		if (r >= '0' && r <= '9') || r == '.' {
			numStr.WriteRune(r)
		} else if unicode.IsLetter(r) {
			// Берем только первый буквенный символ как суффикс
			if suffix == "" {
				suffix = strings.ToLower(string(r))
			}
			// После суффикса прекращаем парсинг
			break
		} else if r == ' ' || r == '\t' {
			// Пробельные символы означают конец числа
			break
		}
		// Игнорируем другие символы
	}

	if numStr.Len() == 0 {
		return 0
	}

	num, err := strconv.ParseFloat(numStr.String(), 64)
	if err != nil {
		return 0
	}

	// Применяем множитель в зависимости от суффикса
	switch suffix {
	case "k":
		return num * 1024
	case "m":
		return num * 1024 * 1024
	case "g":
		return num * 1024 * 1024 * 1024
	case "t":
		return num * 1024 * 1024 * 1024 * 1024
	default:
		return num
	}
}

// unique оставляет только уникальные строки
func unique(lines []string) []string {
	if len(lines) == 0 {
		return lines
	}

	seen := make(map[string]bool)
	result := make([]string, 0, len(lines))

	for _, line := range lines {
		if !seen[line] {
			seen[line] = true
			result = append(result, line)
		}
	}

	return result
}

// isSorted проверяет, отсортированы ли строки
func isSorted(lines []string, cfg Config) bool {
	for i := 1; i < len(lines); i++ {
		if compare(lines[i], lines[i-1], cfg) {
			return false
		}
	}
	return true
}
