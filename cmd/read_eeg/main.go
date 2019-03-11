package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/emotions"
)

// First arguments are files with floats on each line
// Second argument is the sign (+1, -1, 0)
// returns in the svm-lifght format
func main() {
	sign := os.Args[len(os.Args)-1]

	for f := 1; f < len(os.Args)-1; f++ {
		cbf := emotions.GetFourierForFile(os.Args[f], 19)

		for _, c := range cbf {
			if !emotions.IsZero(c) {
				fmt.Printf("%s ", sign)
				for i, cc := range c {
					fmt.Printf("%d:%f ", i+1, cc)
				}
				fmt.Printf("\n")
			}
		}
	}
}
