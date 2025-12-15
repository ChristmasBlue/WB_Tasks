package main

import (
	"flag"
	"os"
	"strings"
	"testing"
)

// TestParseFlags тестирует парсинг флагов командной строки
func TestParseFlags(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected Config
	}{
		{
			name: "no flags",
			args: []string{"program"},
			expected: Config{
				separator: "\t",
			},
		},
		{
			name: "numeric sort",
			args: []string{"program", "-n"},
			expected: Config{
				numeric:   true,
				separator: "\t",
			},
		},
		{
			name: "reverse sort",
			args: []string{"program", "-r"},
			expected: Config{
				reverse:   true,
				separator: "\t",
			},
		},
		{
			name: "unique lines",
			args: []string{"program", "-u"},
			expected: Config{
				unique:    true,
				separator: "\t",
			},
		},
		{
			name: "month sort",
			args: []string{"program", "-M"},
			expected: Config{
				month:     true,
				separator: "\t",
			},
		},
		{
			name: "human numeric sort",
			args: []string{"program", "-h"},
			expected: Config{
				human:     true,
				separator: "\t",
			},
		},
		{
			name: "check sorted",
			args: []string{"program", "-c"},
			expected: Config{
				check:     true,
				separator: "\t",
			},
		},
		{
			name: "ignore blanks",
			args: []string{"program", "-b"},
			expected: Config{
				ignoreBlank: true,
				separator:   "\t",
			},
		},
		{
			name: "column sort",
			args: []string{"program", "-k", "2"},
			expected: Config{
				column:    2,
				separator: "\t",
			},
		},
		{
			name: "multiple flags",
			args: []string{"program", "-n", "-r", "-u"},
			expected: Config{
				numeric:   true,
				reverse:   true,
				unique:    true,
				separator: "\t",
			},
		},
		{
			name: "file argument",
			args: []string{"program", "test.txt"},
			expected: Config{
				separator: "\t",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Сохраняем оригинальные аргументы
			oldArgs := os.Args
			defer func() { os.Args = oldArgs }()

			// Устанавливаем тестовые аргументы
			os.Args = tt.args

			// Сбрасываем флаги для чистого теста
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

			cfg := parseFlags()

			if cfg.numeric != tt.expected.numeric {
				t.Errorf("numeric: got %v, want %v", cfg.numeric, tt.expected.numeric)
			}
			if cfg.reverse != tt.expected.reverse {
				t.Errorf("reverse: got %v, want %v", cfg.reverse, tt.expected.reverse)
			}
			if cfg.unique != tt.expected.unique {
				t.Errorf("unique: got %v, want %v", cfg.unique, tt.expected.unique)
			}
			if cfg.month != tt.expected.month {
				t.Errorf("month: got %v, want %v", cfg.month, tt.expected.month)
			}
			if cfg.human != tt.expected.human {
				t.Errorf("human: got %v, want %v", cfg.human, tt.expected.human)
			}
			if cfg.check != tt.expected.check {
				t.Errorf("check: got %v, want %v", cfg.check, tt.expected.check)
			}
			if cfg.ignoreBlank != tt.expected.ignoreBlank {
				t.Errorf("ignoreBlank: got %v, want %v", cfg.ignoreBlank, tt.expected.ignoreBlank)
			}
			if cfg.column != tt.expected.column {
				t.Errorf("column: got %v, want %v", cfg.column, tt.expected.column)
			}
			if cfg.separator != tt.expected.separator {
				t.Errorf("separator: got %q, want %q", cfg.separator, tt.expected.separator)
			}
		})
	}
}

// TestCompare тестирует функцию сравнения строк
func TestCompare(t *testing.T) {
	tests := []struct {
		name     string
		a        string
		b        string
		cfg      Config
		expected bool // true если a < b
	}{
		{
			name:     "basic string sort",
			a:        "apple",
			b:        "banana",
			cfg:      Config{separator: "\t"},
			expected: true,
		},
		{
			name:     "basic string sort reverse",
			a:        "banana",
			b:        "apple",
			cfg:      Config{separator: "\t"},
			expected: false,
		},
		{
			name:     "numeric sort",
			a:        "10",
			b:        "2",
			cfg:      Config{numeric: true, separator: "\t"},
			expected: false, // 10 > 2
		},
		{
			name:     "numeric sort equal numbers",
			a:        "5",
			b:        "5.0",
			cfg:      Config{numeric: true, separator: "\t"},
			expected: false, // 5 == 5.0, затем сравниваем как строки: "5" < "5.0"
		},
		{
			name:     "month sort",
			a:        "Feb",
			b:        "Jan",
			cfg:      Config{month: true, separator: "\t"},
			expected: false, // Feb (2) > Jan (1)
		},
		{
			name:     "human numeric sort equal",
			a:        "1K",
			b:        "1024",
			cfg:      Config{human: true, separator: "\t"},
			expected: false, // 1K == 1024, затем "1K" > "1024" как строки
		},
		{
			name:     "column sort",
			a:        "z\t1",
			b:        "a\t5",
			cfg:      Config{column: 2, separator: "\t"},
			expected: true, // 1 < 5
		},
		{
			name:     "ignore blanks",
			a:        "apple   ",
			b:        "banana",
			cfg:      Config{ignoreBlank: true, separator: "\t"},
			expected: true, // "apple" < "banana" после обрезки пробелов
		},
		{
			name:     "equal strings",
			a:        "apple",
			b:        "apple",
			cfg:      Config{separator: "\t"},
			expected: false, // a == b, поэтому a < b false
		},
		{
			name:     "numeric with trailing spaces",
			a:        "5   ",
			b:        "5.0",
			cfg:      Config{numeric: true, separator: "\t"},
			expected: false, // 5 == 5.0
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compare(tt.a, tt.b, tt.cfg)
			if result != tt.expected {
				t.Errorf("compare(%q, %q) = %v, want %v", tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

// TestGetKey тестирует функцию получения ключа сортировки
func TestGetKey(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		cfg      Config
		expected string
	}{
		{
			name:     "no column",
			line:     "apple\tbanana\tcherry",
			cfg:      Config{separator: "\t"},
			expected: "apple\tbanana\tcherry",
		},
		{
			name:     "with column",
			line:     "apple\tbanana\tcherry",
			cfg:      Config{column: 2, separator: "\t"},
			expected: "banana",
		},
		{
			name:     "column out of range",
			line:     "apple\tbanana",
			cfg:      Config{column: 5, separator: "\t"},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getKey(tt.line, tt.cfg)
			if result != tt.expected {
				t.Errorf("getKey(%q) = %q, want %q", tt.line, result, tt.expected)
			}
		})
	}
}

// TestMonthValue тестирует функцию преобразования месяца в число
func TestMonthValue(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"Jan", 1},
		{"JAN", 1},
		{"january", 1},
		{"jan", 1},
		{"Feb", 2},
		{"Mar", 3},
		{"Apr", 4},
		{"May", 5},
		{"Jun", 6},
		{"Jul", 7},
		{"Aug", 8},
		{"Sep", 9},
		{"Oct", 10},
		{"Nov", 11},
		{"Dec", 12},
		{"Unknown", 0},
		{"", 0},
		{"Ja", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := monthValue(tt.input)
			if result != tt.expected {
				t.Errorf("monthValue(%q) = %d, want %d", tt.input, result, tt.expected)
			}
		})
	}
}

// TestNumValue тестирует функцию преобразования строки в число
func TestNumValue(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"123", 123},
		{"-123", -123},
		{"123.45", 123.45},
		{"  123  ", 123},
		{"abc", 0},
		{"", 0},
		{"123abc", 0},
		{"12 34", 0}, // пробел внутри числа
		{"5.0", 5.0},
		{"5.000", 5.0},
		{"0", 0},
		{"-0", 0},
		{"3.14159", 3.14159},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := numValue(tt.input)
			// Используем epsilon для сравнения чисел с плавающей точкой
			if result != tt.expected && !(result-tt.expected < 0.0001 && result-tt.expected > -0.0001) {
				t.Errorf("numValue(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// TestHumanValue тестирует функцию преобразования человекочитаемых чисел
func TestHumanValue(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"1K", 1024},
		{"1k", 1024},
		{"1.5K", 1.5 * 1024},
		{"2M", 2 * 1024 * 1024},
		{"3G", 3 * 1024 * 1024 * 1024},
		{"4T", 4 * 1024 * 1024 * 1024 * 1024},
		{"512", 512},
		{"1.5", 1.5},
		{"abc", 0},
		{"", 0},
		{"123KB", 123}, // только K учитывается, B игнорируется
		{"2.5k", 2.5 * 1024},
		{"1 K", 1}, // пробел перед K - не парсим как суффикс
		{"10", 10},
		{"0.5m", 0.5 * 1024 * 1024},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := humanValue(tt.input)
			// Используем epsilon для сравнения чисел с плавающей точкой
			if result != tt.expected && !(result-tt.expected < 0.0001 && result-tt.expected > -0.0001) {
				t.Errorf("humanValue(%q) = %f, want %f", tt.input, result, tt.expected)
			}
		})
	}
}

// TestUnique тестирует функцию удаления дубликатов
func TestUnique(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "no duplicates",
			input:    []string{"apple", "banana", "cherry"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "with duplicates",
			input:    []string{"apple", "banana", "apple", "cherry", "banana"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "empty slice",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "all duplicates",
			input:    []string{"apple", "apple", "apple"},
			expected: []string{"apple"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := unique(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("unique() length = %d, want %d", len(result), len(tt.expected))
				return
			}
			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("unique()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}

// TestIsSorted тестирует функцию проверки отсортированности
func TestIsSorted(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		cfg      Config
		expected bool
	}{
		{
			name:     "sorted basic",
			lines:    []string{"a", "b", "c"},
			cfg:      Config{separator: "\t"},
			expected: true,
		},
		{
			name:     "not sorted",
			lines:    []string{"c", "a", "b"},
			cfg:      Config{separator: "\t"},
			expected: false,
		},
		{
			name:     "sorted numeric",
			lines:    []string{"1", "2", "10"},
			cfg:      Config{numeric: true, separator: "\t"},
			expected: true,
		},
		{
			name:     "not sorted numeric",
			lines:    []string{"10", "2", "1"},
			cfg:      Config{numeric: true, separator: "\t"},
			expected: false,
		},
		{
			name:     "empty slice",
			lines:    []string{},
			cfg:      Config{separator: "\t"},
			expected: true,
		},
		{
			name:     "single element",
			lines:    []string{"apple"},
			cfg:      Config{separator: "\t"},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSorted(tt.lines, tt.cfg)
			if result != tt.expected {
				t.Errorf("isSorted() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestSortLines тестирует основную функцию сортировки
func TestSortLines(t *testing.T) {
	tests := []struct {
		name     string
		lines    []string
		cfg      Config
		expected []string
	}{
		{
			name:     "basic sort",
			lines:    []string{"banana", "apple", "cherry"},
			cfg:      Config{separator: "\t"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "reverse sort",
			lines:    []string{"banana", "apple", "cherry"},
			cfg:      Config{reverse: true, separator: "\t"},
			expected: []string{"cherry", "banana", "apple"},
		},
		{
			name:     "numeric sort",
			lines:    []string{"10", "2", "1", "20"},
			cfg:      Config{numeric: true, separator: "\t"},
			expected: []string{"1", "2", "10", "20"},
		},
		{
			name:     "numeric reverse sort",
			lines:    []string{"10", "2", "1", "20"},
			cfg:      Config{numeric: true, reverse: true, separator: "\t"},
			expected: []string{"20", "10", "2", "1"},
		},
		{
			name:     "unique lines",
			lines:    []string{"banana", "apple", "banana", "cherry", "apple"},
			cfg:      Config{unique: true, separator: "\t"},
			expected: []string{"apple", "banana", "cherry"},
		},
		{
			name:     "month sort",
			lines:    []string{"Mar", "Jan", "Feb", "Dec"},
			cfg:      Config{month: true, separator: "\t"},
			expected: []string{"Jan", "Feb", "Mar", "Dec"},
		},
		{
			name:     "human numeric sort",
			lines:    []string{"1K", "512", "2M", "1G"},
			cfg:      Config{human: true, separator: "\t"},
			expected: []string{"512", "1K", "2M", "1G"},
		},
		{
			name:     "ignore blanks - строки уже обработаны в readLines",
			lines:    []string{"banana", "  apple", "cherry"}, // уже без хвостовых пробелов
			cfg:      Config{ignoreBlank: true, separator: "\t"},
			expected: []string{"  apple", "banana", "cherry"},
		},
		{
			name:     "ignore blanks with trailing spaces",
			lines:    []string{"banana   ", "  apple", "cherry"},
			cfg:      Config{ignoreBlank: true, separator: "\t"},
			expected: []string{"  apple", "banana", "cherry"}, // banana без пробелов
		},
		{
			name:     "column sort",
			lines:    []string{"z\t5", "a\t3", "b\t1"},
			cfg:      Config{column: 2, separator: "\t"},
			expected: []string{"b\t1", "a\t3", "z\t5"},
		},
		{
			name:     "combined numeric and reverse",
			lines:    []string{"3", "1", "4", "2"},
			cfg:      Config{numeric: true, reverse: true, separator: "\t"},
			expected: []string{"4", "3", "2", "1"},
		},
		{
			name:     "empty input",
			lines:    []string{},
			cfg:      Config{separator: "\t"},
			expected: []string{},
		},
		{
			name:     "single line",
			lines:    []string{"single"},
			cfg:      Config{separator: "\t"},
			expected: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Для тестов с ignoreBlank нужно предварительно обработать строки
			testLines := make([]string, len(tt.lines))
			copy(testLines, tt.lines)

			if tt.cfg.ignoreBlank {
				for i := range testLines {
					testLines[i] = strings.TrimRight(testLines[i], " \t")
				}
			}

			result := sortLines(testLines, tt.cfg)

			if len(result) != len(tt.expected) {
				t.Errorf("sortLines() length = %d, want %d", len(result), len(tt.expected))
				t.Errorf("Result: %v", result)
				t.Errorf("Expected: %v", tt.expected)
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("sortLines()[%d] = %q, want %q", i, result[i], tt.expected[i])
				}
			}
		})
	}
}
