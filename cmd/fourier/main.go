package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	filename := os.Args[1]
	wf, _ := fourier.Read(filename, 0, 0.97)

	frames := fourier.CutWavFileIntoFrames(wf)
	c, _ := fourier.FftReal(frames[0])
	b := fourier.Bank(c, wf.GetSampleRate(), 23)

	mfccs := fourier.MFCC(b, 13)
	fmt.Printf("%.9f\n", mfccs)
}
