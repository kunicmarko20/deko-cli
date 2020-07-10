package util

import (
	"fmt"
	"os"
)

func Exit(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}
