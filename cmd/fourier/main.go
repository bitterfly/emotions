package main

import (
	"fmt"
	"os"
	"time"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	filename := os.Args[1]
	// wf, err := fourier.Read(filename, 0)
	wf, err := fourier.Read(filename, 0.1)
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
	// // pprof.StopCPUProfile()

	fmt.Printf("frames: %d\n", len(frames))
	M := 16
	banks := make([][]float64, len(frames), len(frames))
	start = time.Now()

	for i, f := range frames {
		banks[i] = fourier.Bank(fourier.FftReal(f), wf.GetSampleRate(), M)
		fmt.Printf("%d ", i)
	}
	fmt.Printf("\nBanks: %s\n", time.Since(start))

	start = time.Now()
	mfccs := fourier.MFCCS(banks)
	fmt.Printf("\nMFCCs: %s\n", time.Since(start))

	for i, mfcc := range mfccs {
		fourier.PlotSignal(mfcc, fmt.Sprintf("/tmp/mfcc%d.png", i))
	}

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
