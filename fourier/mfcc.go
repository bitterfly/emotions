package fourier

import (
	"math"
)

// func melToIndex(M int, m int, sr int, samplesPerFrame int, maxMel float64) int {
// 	return int(melToFreq(maxMel*float64(m)/float64(M+1)) * float64(samplesPerFrame) / float64(sr))
// }

// func IndToMel(M int, i int, sr int, samplesPerFrame int, maxMel float64) int {
// 	return int(freqToMel(float64(sr*i)/float64(samplesPerFrame)) * float64(M+1) / maxMel)
// }

func melToIndex(M int, m float64, sr int, n int, maxMel float64) float64 {
	// fmt.Printf("m2int %f\n", melToFreq(maxMel*float64(m)/float64(M+1))*float64((2*n))/float64(sr))
	return melToFreq(maxMel*float64(m)/float64(M+1)) * float64((2 * n)) / float64(sr)
}

func IndToMel(M int, i float64, sr int, n int, maxMel float64) float64 {
	// fmt.Printf("int2m %f\n", freqToMel(float64(sr*i)/float64(2*n))*float64(M+1)/maxMel)
	return freqToMel(float64(sr)*i/float64(2*n)) * float64(M+1) / maxMel
}

func melToFreq(mel float64) float64 {
	return (math.Pow(10, mel/2595.0) - 1) * 700
}

func freqToMel(freq float64) float64 {
	return 2595 * math.Log10(1+freq/700.0)
}

func triangleBank(coefficients []Complex, s, e, center int) float64 {
	sum := 0.0

	for i := s; i <= e; i++ {
		if i < center {
			sum += Magnitude(coefficients[i]) * float64(i-s) / float64(center-e)
		} else {
			sum += Magnitude(coefficients[i]) * float64(e-i) / float64(e-center)
		}
	}

	return sum
}

func Bank(coefficients []Complex, sampleRate int, M int) []float64 {
	maxMel := freqToMel(float64(sampleRate) / 2.0)

	banks := make([]float64, M, M)
	for m := 0; m < M; m++ {
		s := int(melToIndex(M, float64(m), sampleRate, len(coefficients), maxMel))
		center := int(melToIndex(M, float64(m+1), sampleRate, len(coefficients), maxMel))
		e := int(melToIndex(M, float64(m+2), sampleRate, len(coefficients), maxMel))

		banks[m] = triangleBank(coefficients, s, e, center)
	}

	// m := IndToMel(numBanks, 1.0, samplingRate, len(coefficients), maxMel)
	// fmt.Printf("100 to mel: %f\n", m)
	// fmt.Printf("mel to 100: %f\n", melToIndex(numBanks, m, samplingRate, len(coefficients), maxMel))

	// fmt.Printf("1 to mel and back: %f\n", melToIndex(IndToMel(1, samplingRate, len(coefficients)), samplingRate, len(coefficients)))

	return banks
}
