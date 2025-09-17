package main

import (
	"fmt"
	"math/big"
)

func main() {
	for {
		var op int
		fmt.Println("Выберите операцию:\n1. Сложение.\n2. Вычитание.\n3. Умножение.\n4. Деление.\n0. Выход.")
		fmt.Print("Введите цифру нужной операции: ")
		fmt.Scanln(&op)
		if op == 0 {
			break
		}
		switch op {
		case 1: //сложение
			Sum()
		case 2: //вычитание
			Subtraction()
		case 3: //умножение
			Multiplication()
		case 4: //деление
			Division()
		default:
			fmt.Println("Операция не распознанна. Повторите ещё раз.")
		}
	}
}

func Sum() {
	var a, b string
	num1 := new(big.Int)
	num2 := new(big.Int)
	//использую бесконечный цикл для валидации ввода, пока ввод не валидный выхода из цикла не произойдёт
	for {
		fmt.Print("Введите число a: ")
		fmt.Scanln(&a)
		fmt.Print("Введите число b: ")
		fmt.Scanln(&b)
		//делаю проверку на вилдацию используя встроенную функцию SetString, и проверяю флаг успешности операции,
		//проверяю флаг, для избежания использования неактуальных данных
		_, success1 := num1.SetString(a, 10)
		_, success2 := num2.SetString(b, 10)
		if success1 && success2 {
			break
		}
		fmt.Println("Некорректное число, повторите попытку.")
		fmt.Println()
	}
	result := new(big.Int).Add(num1, num2)
	fmt.Println("Результат операции сложения: ", result)
}

func Subtraction() {
	var a, b string
	num1 := new(big.Int)
	num2 := new(big.Int)
	//использую бесконечный цикл для валидации ввода, пока ввод не валидный выхода из цикла не произойдёт
	for {
		fmt.Print("Введите число a (уменьшаемое): ")
		fmt.Scanln(&a)
		fmt.Print("Введите число b (вычитаемое): ")
		fmt.Scanln(&b)
		//делаю проверку на вилдацию используя встроенную функцию SetString, и проверяю флаг успешности операции,
		//проверяю флаг, для избежания использования неактуальных данных
		_, success1 := num1.SetString(a, 10)
		_, success2 := num2.SetString(b, 10)
		if success1 && success2 {
			break
		}
		fmt.Println("Некорректное число, повторите попытку.")
		fmt.Println()
	}
	result := new(big.Int).Sub(num1, num2)
	fmt.Println("Результат операции вычитания: ", result)
}

func Multiplication() {
	var a, b string
	num1 := new(big.Int)
	num2 := new(big.Int)
	//использую бесконечный цикл для валидации ввода, пока ввод не валидный выхода из цикла не произойдёт
	for {
		fmt.Print("Введите число a: ")
		fmt.Scanln(&a)
		fmt.Print("Введите число b: ")
		fmt.Scanln(&b)
		//делаю проверку на вилдацию используя встроенную функцию SetString, и проверяю флаг успешности операции,
		//проверяю флаг, для избежания использования неактуальных данных
		_, success1 := num1.SetString(a, 10)
		_, success2 := num2.SetString(b, 10)
		if success1 && success2 {
			break
		}
		fmt.Println("Некорректное число, повторите попытку.")
		fmt.Println()
	}
	result := new(big.Int).Mul(num1, num2)
	fmt.Println("Результат операции умножения: ", result)
}

// Division в функции во время деления числа могут не делиться нацело,
// поэтому в этом случае я сделал округление до 10 символов после запятой
func Division() {
	var a, b string
	num1 := new(big.Float).SetPrec(34)
	num2 := new(big.Float).SetPrec(34)
	//использую бесконечный цикл для валидации ввода, пока ввод не валидный выхода из цикла не произойдёт
	for {
		fmt.Print("Введите число a (делимое): ")
		fmt.Scanln(&a)
		fmt.Print("Введите число b (делитель): ")
		fmt.Scanln(&b)
		//делаю проверку на вилдацию используя встроенную функцию SetString, и проверяю флаг успешности операции,
		//проверяю флаг, для избежания использования неактуальных данных
		_, success1 := num1.SetString(a)
		_, success2 := num2.SetString(b)
		if !(success1 && success2) {
			fmt.Println("Некорректное число, повторите попытку.")
			fmt.Println()
			continue
		}
		if num2.Sign() == 0 {
			fmt.Println("Деление на 0 невозможно, повторите попытку.")
			fmt.Println()
			continue
		}
		break
	}
	result := new(big.Float).Quo(num1, num2).SetPrec(34)
	if result.IsInt() {
		fmt.Printf("Результат операции деления: %s\n", result.Text('f', 0))
	} else {
		fmt.Printf("Результат операции деления: %s\n", result.Text('f', 10))
	}
}
