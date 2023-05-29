package main

import "fmt"

func main() {
	for i := 0; i < 3; i++ {
		i2 := make([]int, 0)
		i2 = append(i2, i)
		fmt.Println(i2)
	}
}
