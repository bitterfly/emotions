package emotions

import (
	"fmt"
	"math"
)

func Zero(x *[]float64) {
	zero(x)
}

func zero(x *[]float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] = 0.0
	}
}

func Divide(x *[]float64, n float64) {
	divide(x, n)
}

func divide(x *[]float64, n float64) {
	for i := 0; i < len(*x); i++ {
		(*x)[i] /= n
	}
}

func Add(x *[]float64, y []float64) {
	add(x, y)
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

func combineSlices(a, b []string) []string {
	c := make([]string, len(a)+len(b), len(a)+len(b))
	for i := 0; i < len(a); i++ {
		c[i] = a[i]
	}

	for j := 0; j < len(b); j++ {
		c[j+len(a)] = b[j]
	}

	return c
}

// AverageSlice accumulates every average elements of the array x
func AverageSlice(x [][]float64, average int) [][]float64 {
	averagedSlice := make([][]float64, 0, len(x)/average)
	currentSlice := make([]float64, len(x[0]), len(x[0]))
	for xi := range x {
		if xi != 0 && xi%average == 0 {

			averagedSlice = append(averagedSlice, multiplied(currentSlice, 1/float64(average)))
			zero(&currentSlice)
		}

		add(&currentSlice, x[xi])
	}

	averagedSlice = append(averagedSlice, multiplied(currentSlice, 1/float64(average)))

	if len(averagedSlice) != (len(x)+average-1)/average {
		panic(fmt.Sprintf("Len: %d shoudl be: %d", len(averagedSlice), (len(x)+average)/average))
	}

	return averagedSlice
}
