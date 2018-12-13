package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	filename := os.Args[1]
	wf, err := fourier.Read(filename)
	if err != nil {
		panic(err.Error)
	}

	fmt.Printf("len wav file: %d\n", len(wf.GetData()))
	// fourier.PlotSignal(wf.GetData(), "signal.png")
	fmt.Printf("Dft\n")
	fourier.PlotCoefficients(fourier.DftReal(wf.GetData()), "spectrum.png")
	fmt.Printf("end\n")

	frames := fourier.CutWavFileIntoFrames(wf)

	fmt.Printf("frames: %d\n", len(frames))
	for i, f := range frames {
		fmt.Printf("%d, len: %d\n ", i, len(f))
		fourier.PlotSignal(f, fmt.Sprintf("signal%d.png", i))
		fourier.PlotCoefficients(fourier.FftReal(f), fmt.Sprintf("spectrum%d.png", i))
	}

	// fourier.Bank(fourier.FftReal(frames[0]), sr, 3)

	// coefficients := fourier.FftReal(s)
	// fourier.PrintCoefficients(coefficients)
	// fourier.PlotCoefficients(coefficients, "bla.png")
	// fourier.PlotSignal(s, "signal.png")
}
