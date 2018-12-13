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

	n := 32768
	// sr := 44100
	s := make([]float64, n, n)
	t := make([]fourier.Complex, n, n)

	for i := 0; i < n; i++ {
		// s[i] = math.Cos(math.Pi * float64(2*i*22050*sr) / float64(n))
		// t[i] = fourier.Complex{Re: math.Cos(math.Pi * float64(2*i*22050*sr) / float64(n)), Im: 0.0}

		s[i] = math.Pow(-1.0, float64(i))
		t[i] = fourier.Complex{Re: s[i], Im: 0.0}
	}

	fmt.Printf("Fast fourier\n")
	fc := fourier.FftReal(s)
	fmt.Printf("Slow fourier\n")
	sc := fourier.Dft(t)

	fourier.PrintCoefficients(fc)
	fourier.PrintCoefficients(sc)
	// frames := fourier.CutSliceIntoFrames(s, uint32(sr))
	// fourier.Bank(fourier.FftReal(frames[0]), sr, 3)

	// coefficients := fourier.FftReal(s)
	// fourier.PrintCoefficients(coefficients)
	// fourier.PlotCoefficients(coefficients, "bla.png")
	// fourier.PlotSignal(s, "signal.png")
}
