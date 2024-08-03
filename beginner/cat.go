package main

import (
	"bytes"
	"fmt"
	"strconv"
)

func main() {
	var buffer bytes.Buffer

	for i := 0; i < 32; i++ {
		buffer.WriteString("a")
	}
	fmt.Printf("%v\n", buffer)

	s := buffer.String()
	fmt.Printf("%v\n", s)

	for i := 0; i < 8192; i++ {
		s += strconv.Itoa(i)
	}
	fmt.Println(s)

	s += buffer.String()
	for i := 0; i < 512; i++ {
		s = "pre" + s
	}
	fmt.Println(s)
}
