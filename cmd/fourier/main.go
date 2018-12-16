package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	filename := os.Args[1]
	wf, err := fourier.Read(filename)
	if err != nil {
		panic(err.Error)
	}

	// fmt.Printf("len wav file: %d\n", len(wf.GetData()))
	// fourier.PlotSignal(wf.GetData(), "signal.png")
	// fmt.Printf("Dft\n")

	// f, err := os.Create("/tmp/framesprof")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// pprof.StartCPUProfile(f)
	start := time.Now()
	frames := fourier.CutWavFileIntoFrames(wf)
	fmt.Printf("Frames: %s\n", time.Since(start))
	// pprof.StopCPUProfile()

	fmt.Printf("frames: %d\n", len(frames))

	start = time.Now()
	for i, f := range frames {
		bank := fourier.Bank(fourier.FftReal(f), wf.GetSampleRate(), 16)
		fmt.Printf("%d ", i)
		_ = fourier.MFCC(bank)
	}
	fmt.Printf("\nBanks: %s\n", time.Since(start))

	// fmt.Printf("frames: %d\n", len(frames))
	// for i, f := range frames {
	// 	fmt.Printf("%d, len: %d\n ", i, len(f))
	// 	fourier.PlotSignal(f, fmt.Sprintf("signal%d.png", i))
	// 	fourier.PlotCoefficients(fourier.FftReal(f), fmt.Sprintf("spectrum%d.png", i))
	// }

	// coefficients := fourier.FftReal(s)
	// fourier.PrintCoefficients(coefficients)
	// fourier.PlotCoefficients(coefficients, "bla.png")
	// fourier.PlotSignal(s, "signal.png")
}
