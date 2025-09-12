package main

import (
	"fmt"
)

type Human struct {
	X int
	Y int
}

type Action struct {
	Human
	Z int
}

func (h *Human) Print() {
	fmt.Printf("X = %d, Y = %d\n", h.X, h.Y)
}

func main() {
	action := Action{
		Human: Human{
			X: 5,
			Y: 10,
		},
		Z: 15,
	}
	action.Print()
	fmt.Println("Z = ", action.Z)
	action.X *= 2
	action.Y *= 2
	action.Z *= 2
	action.Print()
	fmt.Println("Z = ", action.Z)
}
