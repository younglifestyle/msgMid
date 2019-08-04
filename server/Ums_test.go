package server

import (
	"fmt"
	"testing"
)

func TestRunHttpServer(t *testing.T) {
	us := make([]bool, 8, 10)

	for _, value := range us {
		fmt.Println(value)
	}
	fmt.Println()

	u2 := [8]bool{0: true}
	for _, value := range u2 {
		fmt.Println(value)
	}

}
