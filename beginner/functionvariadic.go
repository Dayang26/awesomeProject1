package main

import "fmt"

func main() {
	variadic(1, 2, 3, 4, 5, 6, 7, 8, 9)
	variadic()
}

func variadic(numbers ...int) {
	fmt.Printf("Type: %T\t Content: %d\n", numbers, numbers)
}
