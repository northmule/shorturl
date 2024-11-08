package pkg2

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func main() {
	res, _ := mulfunc(5)
	fmt.Println(res)
	os.Exit(1)
}
