package main

import "fmt"

func reverse(a []string) {
	for i, j := 0, len(a)-1; i < j; i, j = i+1, j-1 {
		a[i], a[j] = a[j], a[i]
	}
}

func main() {
	a := []string{"1", "2", "3"}
	fmt.Println(a)
	reverse(a)
	fmt.Println(a)
}
