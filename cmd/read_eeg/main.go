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

// First arguments are files with floats on each line
// Second argument is the sign (+1, -1, 0)
// returns in the svm-lifght format
func main() {
	sign := os.Args[len(os.Args)-1]

	for f := 1; f < len(os.Args)-1; f++ {
		cbf := emotions.GetFourierForFile(os.Args[f], 19)

		for _, c := range cbf {
			if !isZero(c) {
				fmt.Printf("%s ", sign)
				for i, cc := range c {
					fmt.Printf("%d:%f ", i+1, cc)
				}
				fmt.Printf("\n")
			}
		}
	}
}
