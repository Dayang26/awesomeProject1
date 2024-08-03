package main

import (
	"fmt"
	"strings"
)

var variable1 string
var strarray []string

func main() {
	variable1 = "Lorem Ipsum Dolor sit Amet"
	fmt.Println(variable1)

	strarray = strings.Split(variable1, " ")
	//strarray = strings.Fields(variable1)
	for i := 0; i < len(strarray); i++ {
		fmt.Println(strarray[i])
	}
}
