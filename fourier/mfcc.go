package fourier

// import (
// 	"fmt"
// 	"math"
// )

// func melToInd(mel float64) float64 {
// 	return (math.Pow(10, mel/2595.0) - 1) * 700 * n / float64(sr)
// }

// func IndToMel(ind float64, sr int, n int) float64 {
// 	return 2595 * math.Log10(1+i*sr/n*700.0)
// }

// func Check() {

// 	// freq := 500 //Hz

// 	// sr = 1000
// 	// n = 2000

// 	//sec = 2

// 	panic(fmt.Sprintf("%f %f", freqToMel(500), indToMel(1000)))

// }

// func melToFreq(mel float64) float64 {
// 	return (math.Pow(10, mel/2595.0) - 1) * 700
// }

// func freqToMel(freq float64) float64 {
// 	return 2595 * math.Log10(1+freq/700.0)
// }

// //sr
// //n
// //
// //sec = sr/n
// //
// //i
// //freq = i / sec
// //
// //i * n / sr

// func bank(coefficients []Complex, numBanks int) []Complex {
// 	return nil
// }
