package fourier

import (
	"fmt"
	"image/color"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

// FindClosestPower finds the closest power of two above the given number
func FindClosestPower(x int) int {
	p := 1

	for p < x {
		p *= 2
	}

	return p
}

//IsPowerOfTwo returns whether a number is a power of two
func IsPowerOfTwo(x int) bool {
	return (x & (x - 1)) == 0
}

// PrintCoefficients prints only those fourier coefficients that are greater than ɛ
func PrintCoefficients(coefficients []Complex) {
	fmt.Printf("Number of coefficients: %d\n", len(coefficients))
	for i, c := range coefficients {
		if Magnitude(c) > EPS {
			fmt.Printf("%d: %s\n", i, c)
		}
	}
}

// PlotSignal draws a graph of the given signal and saves it into a file
func PlotSignal(data []float64, file string) {
	plots, err := plot.New()
	if err != nil {
		panic(err)
	}

	for _, b := range data {
		if math.IsInf(b, 0) {
			fmt.Printf("inf\n")
		} else {
			fmt.Printf("%f\n", b)
		}
	}

	s := make(plotter.XYs, len(data))
	for i := 0; i < len(data); i++ {
		s[i].X = float64(i)
		if math.IsInf(data[i], 0) {
			s[i].Y = -10
		} else {
			s[i].Y = data[i]
		}
	}

	line, _ := plotter.NewLine(s)
	line.Color = color.RGBA{0, 232, 88, 255}

	plots.Add(line)

	if err := plots.Save(32*vg.Inch, 16*vg.Inch, file); err != nil {
		panic(err)
	}
}

// func PlotSignals(data [][]float64, offsets []int, file string) {
// 	plots, err := plot.New()
// 	if err != nil {
// 		panic(err)
// 	}

// 	for j, d := range data {
// 		s := make(plotter.XYs, len(d))
// 		for i := 0; i < len(d); i++ {
// 			s[i].X = float64(i + offsets[j])
// 			s[i].Y = d[i]
// 		}

// 		line, _ := plotter.NewLine(s)
// 		line.Color = color.RGBA{0, 232, 88, 255}

// 		plots.Add(line)
// 	}

// 	if err := plots.Save(32*vg.Inch, 16*vg.Inch, file); err != nil {
// 		panic(err)
// 	}
// }

// PlotCoefficients draws a bar plot of the fourier coefficients and saves it into a file
func PlotCoefficients(coefficients []Complex, file string) {
	v := make(plotter.Values, len(coefficients))
	max := 0.0
	j := 0
	for i := range v {
		v[i] = Magnitude(coefficients[i])
		if v[i]-max > EPS {
			max = v[i]
			j = i
		}
	}

	fmt.Printf("%d: %f\n", j, max)

	plotc, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotc.X.Min = 0
	plotc.X.Max = float64(len(coefficients))
	plotc.X.Label.Text = "Frequency"
	plotc.Y.Label.Text = "Energy"

	bars, err := plotter.NewBarChart(v, 2)
	plotc.Add(bars)

	if err := plotc.Save(16*vg.Inch, 16*vg.Inch, file); err != nil {
		panic(err)
	}

}

func PrintFrameSlice(frames [][]float64) {
	for i, frame := range frames {
		fmt.Printf("%d\n", i)
		PlotSignal(frame, fmt.Sprintf("signal/signal%d.png", i))
		PlotCoefficients(FftReal(frame), fmt.Sprintf("spectrum/spectrum%d.png", i))

	}
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sliceCopy(first []float64, from, to, length int) []float64 {
	second := make([]float64, length, length)
	copy(second, first[from:Min(to, len(first))])
	hanningWindow(second[0 : Min(to, len(first))-from])
	return second
}
