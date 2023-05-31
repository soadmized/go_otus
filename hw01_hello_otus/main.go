package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	str := "Hello, OTUS!"
	res := stringutil.Reverse(str)
	fmt.Println(res)
}
