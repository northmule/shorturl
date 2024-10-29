package main

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	if i == 10 {
		os.Exit(12)
	}
	return i * 2, nil
}

func main() {
	res, _ := mulfunc(5)
	fmt.Println(res)
	os.Exit(1) // want "direct function os.Exit call is prohibited"
}
