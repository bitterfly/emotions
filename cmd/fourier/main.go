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

	n := 8
	// sr := 44100
	s := make([]float64, n, n)
	t := make([]fourier.Complex, n, n)

	for i := 0; i < n; i++ {
		s[i] = math.Pow(-1.0, float64(i))
		t[i] = fourier.Complex{Re: s[i], Im: 0.0}
	}

	a := make([]fourier.Complex, n, n)
	a[0] = fourier.Complex{Re: 1.0, Im: -1.0}
	a[1] = fourier.Complex{Re: -1.0, Im: 0.0}
	a[2] = fourier.Complex{Re: 0.0, Im: 1.0}
	a[3] = fourier.Complex{Re: 0.5, Im: 0.5}
	a[4] = fourier.Complex{Re: 0.5, Im: -0.5}
	a[5] = fourier.Complex{Re: -0.5, Im: 2.0}
	a[6] = fourier.Complex{Re: 2, Im: 1.5}
	a[7] = fourier.Complex{Re: -1.5, Im: 0.5}

	fmt.Printf("Slow fourier:\n")
	sc := fourier.Dft(a)
	fourier.PrintCoefficients(sc)

	fmt.Printf("Fast fourier:\n")
	fc := fourier.Fft(a)
	fourier.PrintCoefficients(fc)

	// frames := fourier.CutSliceIntoFrames(s, uint32(sr))
	// fourier.Bank(fourier.FftReal(frames[0]), sr, 3)

	// coefficients := fourier.FftReal(s)
	// fourier.PrintCoefficients(coefficients)
	// fourier.PlotCoefficients(coefficients, "bla.png")
	// fourier.PlotSignal(s, "signal.png")
}
