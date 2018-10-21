package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/bitterfly/emotions/fourier"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {

	// s := []fourier.Complex{fourier.Complex{Re: 2, Im: 3}, fourier.Complex{Re: 1, Im: -1}}
	// fourier.PrintSignal(s)
	// fmt.Printf("<e1, e2>: %s\n", fourier.Edot(1, 2, 3))
	// fmt.Printf("<e2, e2>: %s\n", fourier.Edot(2, 2, 3))

	// n := 16000
	// old := make([]fourier.Complex, n, n)
	// for i := 0; i < n; i++ {
	// 	// T = 1 sec
	// 	// n =  16000
	// 	// 800Hz

	// 	old[i] = fourier.Complex{
	// 		// Re: math.Cos(3*math.Pi/2.0 + math.Pi*float64(2*i*800)/float64(n)),
	// 		// Re: math.Cos(math.Pi/2 + math.Pi*float64(2*i*800)/float64(n)),
	// 		Im: 0.0,
	// 	}
	// }

	wf, err := fourier.Read(os.Args[1])
	if err != nil {
		panic(err)
	}

	signal := wf.GetData()

	coefficients, _ := fourier.Dft(signal)
	panic(coefficients[1000])

	inverseSignal := fourier.Idft(coefficients)

	v := make(plotter.Values, len(coefficients))
	for i := range v {
		if fourier.Magnitude(coefficients[i]) > fourier.EPS {
			fmt.Printf("%d %.3f\n", i, fourier.Magnitude(coefficients[i]))
		}
		v[i] = fourier.Magnitude(coefficients[i])
	}

	f := make(plotter.XYs, 50)
	s := make(plotter.XYs, 50)
	for i := 0; i < 50; i++ {
		fmt.Printf("%d %s\n", i, inverseSignal[i])

		f[i].X = float64(i)
		f[i].Y = signal[i].Re
		// f[i].Y = fourier.Magnitude(signal[i])

		s[i].X = float64(i)
		s[i].Y = inverseSignal[i].Re
		// s[i].Y = fourier.Magnitude(inverseSignal[i])

	}

	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	// p.X.Min = 0
	// p.X.Max = float64(len(coefficients))
	// p.X.Label.Text = "Frequency"
	// p.Y.Label.Text = "Energy"
	// bars, err := plotter.NewBarChart(v, 2)
	// bars.Color = color.RGBA{10, 120, 120, 1}
	// p.Add(bars)

	fl, _ := plotter.NewLine(f)
	fl.Color = color.RGBA{254, 1, 2, 1}

	sl, _ := plotter.NewLine(s)
	p.Add(fl)
	p.Add(sl)

	if err := p.Save(4*vg.Inch, 4*vg.Inch, "after.png"); err != nil {
		panic(err)
	}
}
