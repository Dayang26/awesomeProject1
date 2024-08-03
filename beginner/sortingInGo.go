package main

import (
	"fmt"
	"sort"
)

func main() {
	strs := []string{"c", "b", "a", "e", "d", "z", "g"}
	sort.Strings(strs)
	fmt.Println("Strings:", strs)

	ints := []int{10, 3, 2, 5, 6, 8, 3, 0, 4, 7}
	sort.Ints(ints)
	fmt.Println("Ints: ", ints)

	s := sort.IntsAreSorted(ints)
	fmt.Println("Sorted: ", s)
}
