package main

import (
	"fmt"
)

func main() {
	var s string
	fmt.Print("please insert a string an press enter")
	scan, err := fmt.Scan(&s)
	if err != nil {
		return
	} else {
		fmt.Printf("read string \"%v\" from  stdin \n", s)
		fmt.Printf("read string \"%v\" from  stdin \n", scan)
	}
}
