package fourier

import (
	"fmt"
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

	fmt.Printf("%d %d %d\n", s, center, e)

	var power float64
	for i := s; i <= e; i++ {
		power = Power(coefficients[i])

		if i < center {
			fmt.Printf("%d %f %f %f\n", i, power, float64(i-s)/float64(center-s), power*float64(i-s)/float64(center-e))

			sum += power * float64(i-s) / float64(center-s)
		} else {
			sum += power * float64(e-i) / float64(e-center)
			fmt.Printf("%d %f %f %f\n", i, power, float64(e-i)/float64(e-center), power*float64(e-i)/float64(e-center))

		}
	}

	return math.Log(sum)
}

func MFCCS(banks [][]float64) [][]float64 {
	M := len(banks[0])
	cosines := make([][]float64, M, M)
	for n := 0; n < M; n++ {
		cosines[n] = make([]float64, M, M)
		for m := 0; m < M; m++ {
			cosines[n][m] = math.Cos(math.Pi * float64(n) * (float64(m) + 0.5) / float64(M))
		}
	}

	mfccs := make([][]float64, len(banks), len(banks))

	for i, bank := range banks {
		mfccs[i] = make([]float64, M, M)
		for n := 0; n < M; n++ {
			for m := 0; m < M; m++ {
				mfccs[i][n] += bank[m] * cosines[n][m]
			}
		}
	}

	return mfccs
}

func MFCC(bank []float64) []float64 {
	M := len(bank)
	cosines := make([][]float64, M, M)
	for n := 0; n < M; n++ {
		cosines[n] = make([]float64, M, M)
		for m := 0; m < M; m++ {
			cosines[n][m] = math.Cos(math.Pi * float64(n) * (float64(m) + 0.5) / float64(M))
		}
	}

	mfcc := make([]float64, M, M)
	for n := 0; n < M; n++ {
		for m := 0; m < M; m++ {
			mfcc[n] += bank[m] * cosines[n][m]
		}
	}

	return mfcc
}

func Cepstrum(coefficients []Complex, samplerate int) []float64 {
	return nil
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

	return banks
}
