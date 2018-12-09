package main

import (
	"fmt"
	"math"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	//	filename := os.Args[1]
	//	wf, err := fourier.Read(filename)
	//	if err != nil {
	//		panic(err.Error)
	//	}

	fmt.Printf("%d\n", fourier.FindClosestPower(513))

	n := 40000
	s := make([]float64, n, n)
	for i := 0; i < n; i++ {
		s[i] = math.Cos(math.Pi * float64(2*i*30) / float64(n))
	}

	// 8192
	frames := fourier.CutSliceIntoFrames(s, 40000/8192.0, 8192)
	fourier.PrintFrameSlice(frames)

	// coefficients := fourier.FftReal(s)
	// fourier.PrintCoefficients(coefficients)
	// fourier.PlotCoefficients(coefficients, "bla.png")
	// fourier.PlotSignal(s, "signal.png")
}
