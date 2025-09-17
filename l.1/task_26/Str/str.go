package str

import "strings"

type Str struct {
	str   string
	cache map[rune]bool
}

//NewStr конструктор
func NewStr(str string) *Str {
	return &Str{str: str, cache: make(map[rune]bool)}
}

//CheckUniqSymb функция проверки уникальности символов, пробелы тоже считаются символами
func (s *Str) CheckUniqSymb() bool {
	s.str = strings.ToLower(s.str)
	for _, symb := range s.str {
		if s.cache[symb] {
			return false
		}
		s.cache[symb] = true
	}
	return true
}
