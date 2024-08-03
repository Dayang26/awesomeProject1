package main

import "fmt"

type Person struct {
	Name string
	Age  int
}

func main() {
	p := Person{Name: "John", Age: 42}

	fmt.Printf("%v,%T,%#v\n", p, p, p)

	fmt.Printf("%v,%T,%#v\n", p.Name, p.Name, p.Name)

	fmt.Printf("%v,%T,%#v\n", p.Age, p.Age, p.Age)

	p.Name = "aaron"
	p.Age = 22
	fmt.Println(p)
	fmt.Printf("%v", p)
}
