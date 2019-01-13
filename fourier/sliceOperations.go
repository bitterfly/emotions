package fourier

import "math"

func zero(x *[]float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = 0.0
	}
}

func divide(x *[]float64, n float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] /= n
	}
}

func add(x *[]float64, y []float64) {
	for i := 0; i < len(y); i++ {
		(*x)[i] += y[i]
	}
}

func multiply(x *[]float64, y float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] *= y
	}
}

func multiplied(x []float64, y float64) []float64 {
	z := make([]float64, len(x), len(x))
	for i := 0; i < len(z); i++ {
		z[i] = x[i] * y
	}

	return z
}

func inverse(x *[]float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = float64(1) / (*x)[i]
	}
}

func square(x *[]float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = (*x)[i] * (*x)[i]
	}
}

func eps(x *[]float64, epsilon float64) {
	for i := 0; i < len(*x); i++ {
		if (*x)[i] < epsilon {
			(*x)[i] = epsilon
		}
	}
}

func minused(x []float64, y []float64) []float64 {
	z := make([]float64, len(x), len(x))
	for i := 0; i < len(x); i++ {
		z[i] = x[i] - y[i]
	}
	return z
}

func getSqrt(x *[]float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = math.Sqrt((*x)[i])
	}
}
