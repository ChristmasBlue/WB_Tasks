package main

import (
	"fmt"
	"os"
)

type Service interface {
	Read(string) string
	Save(string)
}

type Client interface {
	Write()
}

type Adapter struct {
	Service
}

func (a *Adapter) Write() {
	var str string
	fmt.Scanln(&str)
	a.Save(str)

}

type Ser struct {
	cache map[int]string
	count int
}

func (s *Ser) Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(data)
}

func (s *Ser) Save(str string) {
	if s.cache == nil {
		s.cache = make(map[int]string)
	}
	s.count++
	s.cache[s.count] = str
}

type Cli struct {
	mess string
}

func (c *Cli) Write() {
	fmt.Scan(&c.mess)
}

func clientCode(c Client) {
	c.Write()
}

func main() {
	myService := &Ser{}
	serviceAdapter := Adapter{Service: myService}

	directClient := &Cli{}
	directClient.Write()
	fmt.Println(directClient.mess)

	clientCode(&serviceAdapter)
	for key, val := range myService.cache {
		fmt.Println(key, " ", val)
	}
}
