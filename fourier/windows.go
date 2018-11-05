package fourier

import "math"

const WINDOW int = 1024
const INNER_WINDOW = 768

func f(x float64) float64 {
	return 1 / (1 + math.Exp(-x))
}

func rectangular(x float64) float64 {
	return 1
}

func hanning(x float64) float64 {
	return 0.5 - 0.5*math.Cos(2*math.Pi*x/INNER_WINDOW)
}

func hamming(x float64) float64 {
	return 0.54 - 0.46*math.Cos(2*math.Pi*x/INNER_WINDOW)
}

func window(content []float64, windowFunction func(float64) float64) {
	for i, x := range content {
		content[i] = x * (windowFunction(float64(i)-float64(uttBeg)) * windowFunction(float64(uttEnd)-float64(i)))
	}
}

func hammingWindow(content []float64, uttBeg, uttEnd uint64) {
	window(content, uttBeg, uttEnd, hamming)
}

func hanningWindow(content []float64, uttBeg, uttEnd uint64) {
	window(content, uttBeg, uttEnd, hanning)
}

func rectangularWindow(content []float64, uttBeg, uttEnd uint64) {
	window(content, uttBeg, uttEnd, rectangular)
}
