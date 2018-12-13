package fourier

import (
	"fmt"
	"math"
)

func melToIndex(mel float64, sr int, n int) float64 {
	return (math.Pow(10, mel/2595.0) - 1) * float64(700*n) / float64(sr)
}

func IndToMel(ind float64, sr int, n int) float64 {
	return 2595 * math.Log10(1+float64(sr)*ind/float64(n)*700.0)
}

func melToFreq(mel float64) float64 {
	return (math.Pow(10, mel/2595.0) - 1) * 700
}

func freqToMel(freq float64) float64 {
	return 2595 * math.Log10(1+freq/700.0)
}

func Bank(coefficients []Complex, samplingRate int, numBanks int) []Complex {
	fmt.Printf("%d\n", len(coefficients))
	// fmt.Printf("1 to mel and back: %f\n", melToIndex(IndToMel(1, samplingRate, len(coefficients)), samplingRate, len(coefficients)))
	// fmt.Printf("1 to mel and back: %f\n", melToIndex(IndToMel(1, samplingRate, len(coefficients)), samplingRate, len(coefficients)))

	return nil
}
