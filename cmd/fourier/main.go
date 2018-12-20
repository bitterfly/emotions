package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	filename := os.Args[1]
	wf, _ := fourier.Read(filename, 0, 0.97)

	mfccs := fourier.MFCCs(wf, 13, 23)
	doubles := fourier.MFCCcDouble(mfccs)

	for i, d := range doubles {
		fmt.Printf("%d %v\n", i, d)
		fmt.Printf("\n")
	}
}
