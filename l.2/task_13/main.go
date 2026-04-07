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

type Config struct {
	//входные аргументы, нужно парсить
	FieldsStr    string
	DelimiterStr string
	//аргументы которые можно использовать(уже распаршены)
	Separated bool
	Delimiter rune
	Fields    []int
	//файлы из которых нужно читать
	Files []string
}

func main() {
	cfg, err := parseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(cfg.Files) > 0 {
		for _, filename := range cfg.Files {
			err := processFile(filename, cfg)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		}
	} else {
		err := processRead(os.Stdin, cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
}

// processFile обработкак входных данных через файл
func processFile(filename string, cfg *Config) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open file: %s, %v", filename, err)
	}

	defer file.Close()
	err = processRead(file, cfg)
	if err != nil {
		return err
	}
	return nil
}

// processRead обработка данных чтением
func processRead(reader *os.File, config *Config) error {

	scanner := bufio.NewScanner(reader)
	delimiter := string(config.Delimiter)

	for scanner.Scan() {
		line := scanner.Text()

		if config.Separated && !strings.ContainsRune(line, config.Delimiter) {
			continue
		}

		if config.Fields == nil {
			fmt.Println(line)
			continue
		}

		str := strings.Split(line, delimiter)
		maxStrLen := len(str)

		for i, field := range config.Fields {
			if field > maxStrLen {
				continue
			}

			if i != 0 {
				fmt.Printf("%s", delimiter)
			}

			fmt.Printf("%s", str[field-1])
		}
		fmt.Println()

	}

	return scanner.Err()
}

// parseFlags обработка входных данных
func parseFlags() (*Config, error) {
	var cfg Config

	flag.StringVar(&cfg.FieldsStr, "f", "", "список полей для анализа, разделенных запятыми")
	flag.StringVar(&cfg.DelimiterStr, "d", "\t", "разделитель поля")
	flag.BoolVar(&cfg.Separated, "s", false, "использовать разделительные поля")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Использование: %s [опции]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Если файл не указан, читает из STDIN.")
	}

	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		cfg.Files = []string{}
	} else {
		cfg.Files = args
	}

	err := cfg.ParseDelimiter()
	if err != nil {
		return nil, err
	}

	err = cfg.ParseFields()
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

// ParseDelimiter парсит входные аргументы флага -d
func (c *Config) ParseDelimiter() error {
	delim := []rune(c.DelimiterStr)

	if len(delim) == 1 {
		c.Delimiter = delim[0]
		return nil
	}

	if len(delim) > 2 || len(delim) == 0 {
		return fmt.Errorf("invalid delimiter")
	}

	if delim[0] != '\\' {
		return fmt.Errorf("invalid delimiter")
	}

	switch delim[1] {
	case 't':
		c.Delimiter = '\t'

	case 'n':
		c.Delimiter = '\n'

	case 'r':
		c.Delimiter = '\r'

	case '\\':
		c.Delimiter = '\\'

	case '0':
		c.Delimiter = 0

	default:
		return fmt.Errorf("invalid delimiter")
	}

	return nil

}

// ParsedField структура для отслеживания чисел
type ParsedField struct {
	Number      int          //одно число
	Builder     []rune       //собираем число
	NumberStart int          //начало последовательности
	NumberEnd   int          //конец последовательности
	Subsequence bool         //флаг последовательности
	Result      map[int]bool //собираем числа в мапе, чтобы избежать дублирования
}

// ParseFields парсит входные аргументы флага -f
func (c *Config) ParseFields() error {

	if c.FieldsStr == "" {
		c.Fields = nil
		return nil
	}

	parsedFields := &ParsedField{
		Number:      0,
		Builder:     make([]rune, 0),
		NumberStart: 0,
		NumberEnd:   0,
		Subsequence: false,
		Result:      make(map[int]bool),
	}

	var err error

	for _, field := range c.FieldsStr {
		switch {
		case unicode.IsDigit(field): //проверяем, что символ является цифрой
			parsedFields.Builder = append(parsedFields.Builder, field) //отправляем его в билдер

		case field == ',': //если встречаем запятую обрабатываем число или последовательность
			if parsedFields.Subsequence { //проверяем флаг последовательности
				parsedFields.NumberEnd, err = strconv.Atoi(string(parsedFields.Builder)) //конвертируем и сохраняем число на котором остановить последовательность
				if err != nil {
					return err
				}
				parsedFields.Builder = parsedFields.Builder[:0] //очищаем биледер
				if parsedFields.NumberStart > parsedFields.NumberEnd {
					return fmt.Errorf("invalid number")
				}
			} else {
				parsedFields.Number, err = strconv.Atoi(string(parsedFields.Builder)) //конвертируем и сохраняем как одиночное число
				if err != nil {
					return err
				}
				parsedFields.Builder = parsedFields.Builder[:0] //очищаем биледер
				if parsedFields.Number <= 0 {
					return fmt.Errorf("field numbers must be positive, got %d", parsedFields.Number)
				}
			}

			parsedFields.AddNumberResult() //добавляем число или последовательность в мапу

		case field == '-':
			if parsedFields.Subsequence {
				return fmt.Errorf("invalid field: multiple '-' in sequence")
			}
			parsedFields.Subsequence = true //если встречаем "-", то значит это последовательность, выставляем флаг
			if len(parsedFields.Builder) == 0 {
				parsedFields.NumberStart = 1
			} else {
				parsedFields.NumberStart, err = strconv.Atoi(string(parsedFields.Builder)) // конвертируем и сохраняем начало последовательсноти
				if err != nil {
					return err
				}
				parsedFields.Builder = parsedFields.Builder[:0] //очищаем билдер
				if parsedFields.NumberStart <= 0 {
					return fmt.Errorf("field numbers must be positive, got %d", parsedFields.NumberStart)
				}
			}

		default:
			return fmt.Errorf("invalid field")
		}
	}

	if len(parsedFields.Builder) != 0 {
		if parsedFields.Subsequence { //проверяем флаг последовательности
			parsedFields.NumberEnd, err = strconv.Atoi(string(parsedFields.Builder)) //конвертируем и сохраняем число на котором остановить последовательность
			if err != nil {
				return err
			}
			parsedFields.Builder = parsedFields.Builder[:0] //очищаем биледер
			if parsedFields.NumberStart > parsedFields.NumberEnd {
				return fmt.Errorf("invalid number")
			}
		} else {
			parsedFields.Number, err = strconv.Atoi(string(parsedFields.Builder)) //конвертируем и сохраняем как одиночное число
			if err != nil {
				return err
			}
			parsedFields.Builder = parsedFields.Builder[:0] //очищаем биледер
		}

		parsedFields.AddNumberResult() //добавляем число или последовательность в мапу
	} else {
		if parsedFields.Subsequence {
			return fmt.Errorf("invalid subsequence")
		}
	}

	res := make([]int, 0, len(parsedFields.Result))
	for num, _ := range parsedFields.Result {
		res = append(res, num)
	}

	sort.Ints(res)
	c.Fields = res
	return nil
}

// AddNumberResult добавляем уникальные числа
func (p *ParsedField) AddNumberResult() {
	if p.Subsequence {
		for i := p.NumberStart; i <= p.NumberEnd; i++ {
			p.Result[i] = true
		}

		p.NumberStart = 0
		p.NumberEnd = 0
		p.Subsequence = false
	} else {
		p.Result[p.Number] = true
		p.Number = 0
	}
}
