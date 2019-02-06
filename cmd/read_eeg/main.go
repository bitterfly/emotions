package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func isZero(x []float64) bool {
	for _, xx := range x {
		if xx > 0.00001 || xx > -0.00001 {
			return false
		}
	}
	return true
}

func main() {
	cbf := emotions.GetFourierForFile(os.Args[1], 19)

	for _, c := range cbf {
		if !isZero(c) {
			for i, cc := range c {
				fmt.Printf("%d:%f ", i+1, cc)
			}
			fmt.Printf("\n")
		}
	}
}
