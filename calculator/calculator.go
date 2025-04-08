package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	MATH  = "Вывод ошибки, так как строка не является математической операцией."
	SINT  = "Вывод ошибки, так как формат математической операции не удовлетворяет заданию. Необходимо писать два операнда и один оператор (+, -, /, *)."
	RANGE = "Калькулятор умеет работать только с арабскими целыми числами или римскими цифрами от 1 до 10 включительно"
	SCALE = "Вывод ошибки, так как используются одновременно разные системы счисления."
	ZERO  = "Вывод ошибки, так как в римской системе нет числа 0 и отрицательных чисел."
)

var romanMap = map[string]int{
	"C":    100,
	"XC":   90,
	"L":    50,
	"XL":   40,
	"X":    10,
	"IX":   9,
	"VIII": 8,
	"VII":  7,
	"VI":   6,
	"V":    5,
	"IV":   4,
	"III":  3,
	"II":   2,
	"I":    1,
}

var convertIntToRoman = [14]int{
	100,
	90,
	50,
	40,
	10,
	9,
	8,
	7,
	6,
	5,
	4,
	3,
	2,
	1,
}

var a, b *int
var operators = map[string]func() int{
	"+": func() int { return *a + *b },
	"-": func() int { return *a - *b },
	"/": func() int { return *a / *b },
	"*": func() int { return *a * *b },
}
var stringExp []string

func main() {

	fmt.Println("Добро пожаловать, пользователь! Тебя приветствует Калькулус!\n" +
		"Я умею выполнять операции сложения, вычитания, умножения и деления с двумя ЦЕЛЫМИ числами от 1 до 10, записанными в строку. Например, A + B.\n" +
		"Выражения могут содержать как арабские, так и римские числа, но совмещать разные системы нельзя.\n" +
		"Результатом моей работы с арабскими числами могут быть отрицательные числа и ноль. Результатом работы с римскими числами могут быть только положительные числа.")
	for {
		fmt.Println("Введите выражение: ")
		text := reader()
		operator, stringsCheck, numbers, romanNumber := parse(strings.ToUpper(strings.TrimSpace(text)))
		operation(operator, stringsCheck, numbers, romanNumber)
	}
}

func reader() string {
	reader := bufio.NewReader(os.Stdin)
	var text string
	text, _ = reader.ReadString('\n')
	text = strings.ReplaceAll(text, " ", "")
	return text
}

func parse(text string) (string, int, []int, []string) {
	var operator string
	var stringsCheck int
	numbers := make([]int, 0)
	romans := make([]string, 0)
	for i := range operators {
		for _, val := range text {
			if i == string(val) {
				operator += i
				stringExp = strings.Split(text, operator)
			}
		}
	}

	for _, elem := range stringExp {
		num, err := strconv.Atoi(elem)
		if err != nil {
			stringsCheck++
			romans = append(romans, elem)
		} else {
			numbers = append(numbers, num)
		}
	}
	return operator, stringsCheck, numbers, romans
}

func funcArab(numbers []int, operator string) {
	if checkRange(numbers[0]) && checkRange(numbers[1]) {
		val := operators[operator]
		a, b = &numbers[0], &numbers[1]
		fmt.Printf("Результат работы: \n%v\n", val())
	} else {
		panic(RANGE)
	}
}

func funcRoman(romans []string, operator string) {
	romanToInt := make([]int, 0)
	for _, elem := range romans {
		val := romanMap[elem]
		if checkRange(val) {
			romanToInt = append(romanToInt, val)
		} else {
			panic(RANGE)
		}
	}
	if val, ok := operators[operator]; ok {
		a, b = &romanToInt[0], &romanToInt[1]
		intToRoman(val())
	}
}

func intToRoman(tmpRes int) {
	var romanRes string
	if tmpRes == 0 || tmpRes < 0 {
		panic(ZERO)
	}
	for tmpRes > 0 {
		for _, elem := range convertIntToRoman {
			for i := elem; i <= tmpRes; {
				for index, value := range romanMap {
					if value == elem {
						romanRes += index
						tmpRes -= elem
					}
				}
			}
		}
	}
	fmt.Printf("Результат работы: \n%v\n", romanRes)
}

func operation(operator string, stringsCheck int, numbers []int, romans []string) {
	checkOperator(operator)
	switch stringsCheck {
	case 1:
		panic(SCALE)
	case 0:
		funcArab(numbers, operator)
	case 2:
		funcRoman(romans, operator)
	}
}

func checkOperator(operator string) {
	switch {
	case len(operator) > 1:
		panic(SINT)
	case len(operator) < 1:
		panic(MATH)
	}
}

func checkRange(num int) bool {
	if num > 0 && num < 11 {
		return true
	} else {
		return false
	}
}
