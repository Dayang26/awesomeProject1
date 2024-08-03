package main

import (
	"fmt"
)

var strarray1 = []string{"lorem", "ipsum", "dolor", "sit", "amet"}
var intarray = []int{1, 2, 4, 8, 16}
var mapone1 = map[int]string{}
var mapone2 = map[string]interface{}{}

func main() {
	for i := 0; i < 5; i++ {
		fmt.Println(intarray[i], "\t", strarray1[i])
		mapone1[intarray[i]] = strarray1[i]
		mapone2[strarray1[i]] = mapone1
	}

	fmt.Println(mapone1)
	fmt.Println(mapone2)
}
