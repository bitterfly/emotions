package fourier

import (
	"math"
)

func melToIndex(M int, m float64, sr int, n int, maxMel float64) float64 {
	return melToFreq(maxMel*float64(m)/float64(M+1)) * float64(2*(n-1)) / float64(sr)
}

func indToMel(M int, i float64, sr int, n int, maxMel float64) float64 {
	return freqToMel(float64(sr)*i/float64(2*(n-1))) * float64(M+1) / maxMel
}

func melToFreq(mel float64) float64 {
	return (math.Pow(10, mel/2595.0) - 1) * 700
}

func freqToMel(freq float64) float64 {
	return 2595 * math.Log10(1+freq/700.0)
}

func triangleBank(coefficients []Complex, s, e, center float64) float64 {
	sum := 0.0

	var power float64
	for i := int(math.Ceil(s)); i <= int(math.Floor(e)); i++ {
		power = Power(coefficients[i])
		if float64(i) < center {
			sum += power * (float64(i) - s) / (center - s)
		} else {
			sum += power * (e - float64(i)) / (e - center)
		}
	}

	return math.Log(sum)
}

// MFCCS returns the mfcc coefficients for each bank
// so an array of size (len(banks), C)
// where C is the number of mffcs
func MFCCS(banks [][]float64, C int) [][]float64 {
	M := len(banks[0])
	cosines := make([][]float64, C, C)
	for c := 0; c < C; c++ {
		cosines[c] = make([]float64, M, M)
		for m := 0; m < M; m++ {
			cosines[c][m] = math.Cos(math.Pi * float64(c) * (float64(m) + 0.5) / float64(M))
		}
	}

	mfccs := make([][]float64, len(banks), len(banks))

	for i, bank := range banks {
		mfccs[i] = make([]float64, C, C)
		for c := 0; c < C; c++ {
			for m := 0; m < M; m++ {
				mfccs[i][c] += bank[m] * cosines[c][m]
			}
		}
	}

	return mfccs
}

// MFCC returns the mfcc coefficients for the given bank
// so it's an array of size C
// where C is the number of mffcs
func MFCC(bank []float64, C int) []float64 {
	M := len(bank)
	cosines := make([][]float64, C, C)
	for n := 0; n < C; n++ {
		cosines[n] = make([]float64, M, M)
		for m := 0; m < M; m++ {
			cosines[n][m] = math.Cos(math.Pi * float64(n) * (float64(m) + 0.5) / float64(M))
		}
	}

	mfcc := make([]float64, C, C)
	for n := 0; n < C; n++ {
		for m := 0; m < M; m++ {
			mfcc[n] += bank[m] * cosines[n][m]
		}
	}

	return mfcc
}

func Cepstrum(coefficients []Complex, samplerate int) []float64 {
	return nil
}

// Bank takes the fourier coefficient for one frame
// and puts it on the ~mel scale (actually puts it on a logarithmic scale with M banks)
func Bank(coefficients []Complex, sampleRate int, M int) []float64 {
	maxMel := freqToMel(float64(sampleRate) / 2.0)

	banks := make([]float64, M, M)
	for m := 0; m < M; m++ {
		s := melToIndex(M, float64(m), sampleRate, len(coefficients), maxMel)
		center := melToIndex(M, float64(m+1), sampleRate, len(coefficients), maxMel)
		e := melToIndex(M, float64(m+2), sampleRate, len(coefficients), maxMel)

		banks[m] = triangleBank(coefficients, s, e, center)
	}

	return banks
}
