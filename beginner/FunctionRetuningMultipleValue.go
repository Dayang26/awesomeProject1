package main

import "fmt"

func myfunc(p, q int) (int, int, int) {
	return p - q, p * q, p + q
}

func main() {
	var1, var2, var3 := myfunc(4, 2)

	fmt.Printf("Result is:%d", var1)
	fmt.Printf("\nResult is:%d", var2)
	fmt.Printf("\nResult is:%d", var3)
}
