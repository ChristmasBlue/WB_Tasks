/*package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}*/

//выводится "error"
//переменная err имеет тип error, где error - интерфейс,
//переменной err присваивается указатель на структуру, которая удовлетворяет интерфейсу error,
//при проверке err != nil, err не равен nil т.к. err - интерфейс, который хранит указатель на структуру customError,
//и указатель на значение равное nil,
//интерфейс равен nil, только в том случае если динамический тип интерфейса и динамическое значение интерфейса равны nil.