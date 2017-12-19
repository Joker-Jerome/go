package main

import (
	"fmt"
	"bufio"
	"os"
)

func main() {
	h1 := Hex{2, 3}
	h2 := Hex{5, 7}
	h3 := Hex{2, 3}

	fmt.Println(h1 == h3)
	fmt.Println(h1.Distance(h2))
	fmt.Println(h1)
	h1.Closer(h2)
	fmt.Println(h1)
	fmt.Println(h1.Cartesian())

	in := bufio.NewReader(os.Stdin)

	w := NewWorld(12)
	for i := 0; i < 500; i++ {
		fmt.Println(w)
		w.Update()
		fmt.Println()

		in.ReadString('\n')
	}
	fmt.Println(w)

	p := NewPredator(w)
	s := NewScent(p)
	s.Update()
	fmt.Println(s.Origin())
	fmt.Println(s.Position())
}
