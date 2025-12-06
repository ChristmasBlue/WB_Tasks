package unpacking_test

import (
	"task9/unpacking"
	"testing"
)

func TestSubseqString(t *testing.T) {
	var testCases = []struct {
		text  string
		input string
		want  string
	}{
		{"Вход: a4bc2d5e\nВыход: aaaabccddddde", "a4bc2d5e", "aaaabccddddde"},
		{"Вход: abcd\nВыход: abcd (нет цифр — ничего не меняется)", "abcd", "abcd"},
		{"Вход: 45\nВыход: '' (некорректная строка, т.к. в строке только цифры — функция должна вернуть ошибку)", "45", ""},
		{"Вход: ''\nВыход: '' (пустая строка -> пустая строка)", "", ""},
		{`Вход: "qwe\4\5"\nВыход: "qwe45" (4 и 5 не трактуются как числа, т.к. экранированы)`, `qwe\4\5`, "qwe45"},
		{`Вход: "qwe\45"\nВыход: "qwe44444" (\4 экранирует 4, поэтому распаковывается только 5)`, `qwe\45`, "qwe44444"},
		{"Вход: 'abcd@'\nВыход: '' (недопустимый символ @)", "abcd@", ""},
	}

	for _, tt := range testCases {
		t.Run(tt.text, func(t *testing.T) {
			result, _ := unpacking.SubseqString(tt.input)

			if result != tt.want {
				t.Errorf("got %s, want %s", result, tt.want)
			}
		})
	}
}
