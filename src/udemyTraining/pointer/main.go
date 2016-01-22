package main

import "fmt"

func main() {

	c := 234

	fmt.Println("c = ", c)  // variable itself
	fmt.Println("c = ", &c) // reference adress

	var p *int = &c

	*(&c)++
	fmt.Println("c = ", *(&c)) // dereferences The Address

	fmt.Println(p)
	fmt.Println(*p)

	for {
		break
	}

}
