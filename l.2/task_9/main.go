package main

import (
	"fmt"
	"task9/unpacking"
)

func main() {
	str, err := unpacking.SubseqString(`\qwe\45\`)
	if err != nil {
		fmt.Printf("Error: %v", err)
	}
	fmt.Println(str)
}
