package domain

type StoppedProg int

const (
	StopWithoutErr      = StoppedProg(iota) //программа остановлена без ошибок
	StopErrGetTime                          //программа остановлена из-за ошибки получения времени с сервера
	StopErrClearConsole                     //программа остановлена из-за ошибки очистки консоли
)
