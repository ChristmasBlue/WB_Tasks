package unpacking

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"
)

// SubseqString получает на вход строку, которую необходимо распаковать.
// перебирает каждый символ,
// функция рассматривает количество повторяющихся символо до 9 включительно,
// учитываются символы "a"..."z"(в любом регистре), "0"..."9", "\", при использовании других символов функция выдаст ошибку
func SubseqString(str string) (string, error) {

	if str == "" {
		return "", nil
	}

	var newStr strings.Builder
	var symb rune
	var ok bool

	for _, val := range str {

		switch {
		case unicode.IsDigit(val):
			switch {
			case symb == 0 || unicode.IsSpace(symb):
				return "", fmt.Errorf("symbol is invalid: %c", val)
			case symb == '\\':
				_, err := newStr.WriteRune(val)
				if err != nil {
					log.Printf("Error append rune %c in builder: %v\n", val, err)
					return "", err
				}
				ok = true
			default:
				if !ok {
					return "", fmt.Errorf("invalid string")
				}
				x, err := strconv.Atoi(string(val))
				if err != nil {
					return "", fmt.Errorf("can't convert symbol to number: %v", err)
				}
				err = addNewString(&newStr, symb, x-1)
				if err != nil {
					return "", fmt.Errorf("error append rune %c in builder: %v", val, err)
				}
				ok = false
			}
		case unicode.IsLetter(val):
			_, err := newStr.WriteRune(val)
			if err != nil {
				log.Printf("Error append rune %c in builder: %v\n", val, err)
				return "", err
			}
			ok = true
		default:
			if val != '\\' {
				return "", fmt.Errorf("invalid string")
			}
		}
		symb = val
	}

	if ok {
		return "", fmt.Errorf("invalid escape sequence at end of string")
	} else {
		return newStr.String(), nil
	}
}

func addNewString(str *strings.Builder, symb rune, amount int) error {
	for i := 0; i < amount; i++ {
		_, err := str.WriteRune(symb)
		if err != nil {
			log.Printf("Error append rune %c in builder: %v\n", symb, err)
			return err
		}
	}
	return nil
}
