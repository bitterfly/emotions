package fourier

import (
	"fmt"
	"math"
)

const EPS = 0.00001

// Complex is a representation of complex numbers, because I didn't want to use the built in one
type Complex struct {
	Re float64
	Im float64
}

func zeroComplex() Complex {
	return Complex{
		Re: 0.0,
		Im: 0.0,
	}
}

func (c Complex) String() string {
	// a
	if c.Im < EPS && c.Im > 0-EPS {
		if c.Re < EPS && c.Re > 0-EPS {
			return fmt.Sprintf("0")
		}
		return fmt.Sprintf("%.3f", c.Re)
	}

	//(a - bi)
	if c.Im < 0.0 {
		if c.Re < EPS && c.Re > 0-EPS {
			return fmt.Sprintf("0 - i(%.3f)", math.Abs(c.Im))
		}
		return fmt.Sprintf("%.3f - i(%.3f)", c.Re, math.Abs(c.Im))
	}
	// (a + bi)
	if c.Re < EPS && c.Re > 0-EPS {
		return fmt.Sprintf("0 + i(%.3f)", c.Im)
	}
	return fmt.Sprintf("%.3f + i(%.3f)", c.Re, c.Im)
}

func (c *Complex) swapped() {
	temp := c.Re
	c.Re = c.Im
	c.Im = temp
}

func (c Complex) divide(x float64) Complex {
	return Complex{
		Re: c.Re / x,
		Im: c.Im / x,
	}
}

func (c *Complex) divided(x float64) {
	c.Re = c.Re / x
	c.Im = c.Im / x
}

func (c Complex) conjugate() Complex {
	return Complex{Re: c.Re, Im: -c.Im}
}

func (c *Complex) added(o Complex) {
	c.Re = c.Re + o.Re
	c.Im = c.Im + o.Im
}

func (c Complex) add(o Complex) Complex {
	return Complex{
		Re: c.Re + o.Re,
		Im: c.Im + o.Im,
	}
}

//PrintSignal brints the given signal for test purposes in the format a + bi\n...
func PrintSignal(x []Complex) {
	for _, xi := range x {
		fmt.Printf("%s\n", xi)
	}
}

func e(k, j int, n int) Complex {
	return Complex{
		Re: math.Cos(2 * math.Pi * float64(k) * float64(j) / float64(n)),
		Im: math.Sin(2 * math.Pi * float64(k) * float64(j) / float64(n)),
	}
}

func dot(c1, c2 Complex) Complex {
	// c1 = (a + bi)
	// c2 = (c + di)
	// ac - bd + i(bc + ad)

	return Complex{
		Re: c1.Re*c2.Re - c1.Im*c2.Im,
		Im: c1.Im*c2.Re + c1.Re*c2.Im,
	}
}

//Edot returns the dot product of the k-th and the s-th roots of unity in the
// n dimentional space
func Edot(k, s, n int) Complex {
	sum := zeroComplex()
	for i := 0; i < n; i++ {
		sum.added(dot(e(k, i, n), e(s, i, n)))
	}
	sum.divided(float64(n))
	return sum
}

func a(k int, b []Complex, n int) Complex {
	sum := zeroComplex()
	// b[1000] = 0 + i0.5

	for j := 0; j <= n/2-1; j++ {
		d := dot(b[j], e(k, j, n))
		sum.added(d)

		d = dot(b[j+1].conjugate(), e(k, n-j-1, n))
		sum.added(d)
	}

	// sum.divided(float64(len(c)))
	return sum
}

func b_real(k int, x []float64) Complex {
	sum := zeroComplex()

	for j := 0; j < len(x); j++ {
		d := Complex{}
		d.Re = x[j] * math.Cos(2*math.Pi*float64(k)*float64(j)/float64(len(x)))
		d.Im = -x[j] * math.Sin(2*math.Pi*float64(k)*float64(j)/float64(len(x)))
		sum.added(d)
	}

	sum.divided(float64(len(x)))
	return sum
}

func b(k int, x []Complex) Complex {
	sum := zeroComplex()
	for j := 0; j < len(x); j++ {
		d := dot(x[j], e(k, j, len(x)).conjugate())
		sum.added(d)
	}

	sum.divided(float64(len(x)))
	return sum
}

func Power(c Complex) float64 {
	return (c.Re*c.Re + c.Im*c.Im)
}

func Magnitude(c Complex) float64 {
	return math.Sqrt(c.Re*c.Re + c.Im*c.Im)
}

func Idft(c []Complex) []Complex {
	n := (len(c) - 1) * 2
	x := make([]Complex, n, n)

	for k := 0; k < n; k++ {
		x[k] = a(k, c, n)
		x[k].divided(2)
	}

	return x
}

func b_k_fast(x []Complex, W []Complex, depth int, first int, step int, len int) Complex {
	if len == 1 {
		return x[first]
	}

	return b_k_fast(x, W, depth+1, first, step*2, len/2).add(dot(W[depth], b_k_fast(x, W, depth+1, first+step, step*2, len/2))).divide(2)

}

func fft(x []Complex, W [][]Complex) []Complex {
	n := len(x)
	coefficients := make([]Complex, n, n)
	for k := 0; k < n; k++ {
		coefficients[k] = b_k_fast(x, W[k], 0, 0, 1, len(x))
	}

	return coefficients
}

func Fft(x []Complex) []Complex {
	n := len(x)

	W := make([][]Complex, n, n)
	for k := 0; k < n; k++ {
		W[k] = make([]Complex, n, n)
		j := 0
		for m := n; m != 0; m /= 2 {
			W[k][j] = e(1, k, m).conjugate()
			j++
		}
	}
	return fft(x, W)
}

func Ifft(x []Complex) []Complex {
	n := len(x)

	W := make([][]Complex, n, n)
	for k := 0; k < n; k++ {
		W[k] = make([]Complex, n, n)
		j := 0
		for m := n; m != 0; m /= 2 {
			W[k][j] = e(1, k, m).conjugate()
			j++
		}
	}

	for i := 0; i < n; i++ {
		x[i].swapped()
	}

	coefficients := fft(x, W)
	for i := 0; i < n; i++ {
		coefficients[i].swapped()
	}

	return coefficients
}

func FftReal(x []float64) []Complex {
	if !IsPowerOfTwo(len(x)) {
		panic("FFT expects the len of the data to be a power of 2")
	}

	n := len(x)
	even := make([]float64, n/2, n/2)
	odd := make([]float64, n/2, n/2)

	X := make([]Complex, n/2+1, n/2+1)

	for k := 0; k < n/2; k++ {
		X[n/2].Re += x[2*k] - x[2*k+1]
		even[k] = x[2*k]
		odd[k] = x[2*k+1]
	}
	Even, Odd := DoubleReal(even, odd)

	for k := 0; k < n/2; k++ {
		X[k] = Even[k].add(dot(e(k, 1, n).conjugate(), Odd[k]))
	}

	X[0].divided(2.0)
	X[n/2].divided(float64(n))

	return X
}

func FftWav(f WavFile) []Complex {
	if !IsPowerOfTwo(len(f.data)) {
		panic("FFT expects the len of the data to be a power of 2")
	}

	return FftReal(f.data)
}

func DoubleReal(x, y []float64) ([]Complex, []Complex) {
	n := len(x)

	xCoefficients := make([]Complex, n, n)
	yCoefficients := make([]Complex, n, n)

	//TODO: check this
	z := make([]Complex, n, n)
	for i := 0; i < n; i++ {
		xCoefficients[0].Re += x[i]
		yCoefficients[0].Re += y[i]

		z[i].Re = x[i]
		z[i].Im = y[i]
	}

	xCoefficients[0].Re = xCoefficients[0].Re / float64(n)
	yCoefficients[0].Re = yCoefficients[0].Re / float64(n)

	zCoefficients := Fft(z)
	for k := 1; k < n; k++ {
		xCoefficients[k].Re = (zCoefficients[k].Re + zCoefficients[n-k].Re) / 2.0
		xCoefficients[k].Im = (zCoefficients[k].Im - zCoefficients[n-k].Im) / 2.0

		yCoefficients[k].Re = (zCoefficients[k].Im + zCoefficients[n-k].Im) / 2.0
		yCoefficients[k].Im = -(zCoefficients[k].Re - zCoefficients[n-k].Re) / 2.0
	}

	return xCoefficients, yCoefficients
}

//Dft returns the discrete fourier transform
func Dft(x []Complex) []Complex {
	coefficients := make([]Complex, len(x), len(x))
	for k := 0; k < len(x); k++ {
		coefficients[k] = b(k, x)
	}
	return coefficients
}

//Dft returns the discrete fourier transform
func DftReal(x []float64) []Complex {
	coefficients := make([]Complex, len(x)/2+1, len(x)/2+1)
	for k := 0; k < len(x)/2+1; k++ {
		coefficients[k] = b_real(k, x)

		if k > 0 && k < len(x)/2 {
			coefficients[k].divided(0.5)
		}

	}
	return coefficients
}
