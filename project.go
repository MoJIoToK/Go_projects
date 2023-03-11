package main

import "fmt"

func main() {
	fmt.Println("Hello world")

	var age int = 10
	var num float32 = 2.456
	var num1 = 2.4345
	var str = "Hi, man"
	//var str1 = "Whats, up?"
	var res int
	res = age + age
	const pi = 3.14

	var web string = "itProger"
	fmt.Println(age, num, num1, str, res)
	fmt.Println(len(web))
	fmt.Println(web + " is cool \nwebsite")

	var num2 float64 = 4.35564
	fmt.Printf("%f \n%.2f \n%T \n", num2, num2, num2)

	var isDone bool = true
	fmt.Printf("%t \n", isDone)

	if age < 8 {
		fmt.Println("Вам пора в детский сад")
	} else if age == 8 {
		fmt.Println("Вам пора в начальную школу")
	} else if (age > 5) && (age <= 18) {
		var grade = age - 5
		fmt.Println("Пора идти в ", grade, "класс")
	} else {
		fmt.Println("Вам пора в ВУЗ")
	}

	switch age {
	case 5:
		fmt.Println("Вам 5 лет")
	case 10:
		fmt.Println("Вам 10 лет")
	default:
		fmt.Println("Вам неизвестно сколько лет")
	}

	for i := 0; i < 5; i++ {
		fmt.Println(i)
	}

	var arr [3]int
	arr[0] = 45
	arr[1] = 90
	arr[2] = 76
	fmt.Println(arr[1])

	nums := [3]float64{4.23, 5.23, 98.1}
	for j, v := range nums {
		fmt.Println(v, j)
	}

	webSites := make(map[string]float64)

	webSites["itProger"] = 0.8
	webSites["yandex"] = 0.99
	fmt.Println(webSites["itProger"])

	var a = 3
	var b = 2

	// var r int
	// r = summ(a, b)
	// fmt.Println(r)

	a, b = summ(a, b)
	fmt.Println(a, b)

	var nnum = 3
	//Замыкания:
	multiple := func() int {
		nnum *= 2
		return nnum
	}
	fmt.Println(multiple())

	// fmt.Printf("\nОткладывание\n")
	// defer two()
	// one()

	fmt.Println("\nУказатели: ")
	var x = 0
	pointer(&x)
	fmt.Println(x)

	fmt.Println("\nСтруктуры (классы): ")
	bob := Cats{"Bob", 7, 0.87}
	fmt.Println("Bob age is", bob.age)
	fmt.Println("Bob function is", bob.test())
}

func summ(num_1 int, num_2 int) (int, int) {
	var res int
	res = num_1 + num_2
	return res, num_1 * num_2
}

//Откладывание
func one() {
	fmt.Println("1")
}
func two() {
	fmt.Println("2")
}

//Указатели (глобальные и локальные переменные)
func pointer(x *int) {
	*x = 2
}

// структура (классы)
type Cats struct {
	name      string
	age       int
	happiness float64
}

func (cat *Cats) test() float64 {
	return float64(cat.age) * cat.happiness
}
