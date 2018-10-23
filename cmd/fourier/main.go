package main

import (
	"fmt"
	"math"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	n := 4
	old := make([]fourier.Complex, n, n)
	for i := 0; i < n; i++ {
		// T = 1 sec
		// n =  16000
		// 800Hz

		old[i] = fourier.Complex{
			Re: 2 * math.Cos(2-math.Pi*float64(2*i*1)/float64(n)),
			// Re: math.Cos(math.Pi/2 + math.Pi*float64(2*i*800)/float64(n)),
			Im: 0.0,
		}
	}

	// wf, err := fourier.Read(os.Args[1])
	// if err != nil {
	// 	panic(err)
	// }

	// signal := wf.GetData()

	coefficients, _ := fourier.Dft(old)
	for i, c := range coefficients {
		fmt.Printf("%d %s\n", i, c)
	}

	// inverseSignal := fourier.Idft(coefficients)

	// // v := make(plotter.Values, len(coefficients))
	// // for i := range v {
	// // 	if fourier.Magnitude(coefficients[i]) > fourier.EPS {
	// // 		fmt.Printf("%d %s\n", i, coefficients[i])
	// // 		fmt.Printf("%d %.3f\n", i, fourier.Magnitude(coefficients[i]))
	// // 	}
	// // 	v[i] = fourier.Magnitude(coefficients[i])
	// // }

	// f := make(plotter.XYs, 50)
	// s := make(plotter.XYs, 50)
	// for i := 0; i < 50; i++ {
	// 	fmt.Printf("%d %s\n", i, inverseSignal[i])

	// 	f[i].X = float64(i)
	// 	f[i].Y = signal[i].Re
	// 	// f[i].Y = fourier.Magnitude(signal[i])

	// 	s[i].X = float64(i)
	// 	s[i].Y = inverseSignal[i].Re
	// 	// s[i].Y = fourier.Magnitude(inverseSignal[i])

	// }

	// p, err := plot.New()
	// if err != nil {
	// 	panic(err)
	// }

	// // p.X.Min = 0
	// // p.X.Max = float64(len(coefficients))
	// // p.X.Label.Text = "Frequency"
	// // p.Y.Label.Text = "Energy"
	// // bars, err := plotter.NewBarChart(v, 2)
	// // bars.Color = color.RGBA{10, 120, 120, 1}
	// // p.Add(bars)

	// fl, _ := plotter.NewLine(f)
	// fl.Color = color.RGBA{254, 1, 2, 1}

	// sl, _ := plotter.NewLine(s)
	// p.Add(fl)
	// p.Add(sl)

	// if err := p.Save(4*vg.Inch, 4*vg.Inch, "forward.png"); err != nil {
	// 	panic(err)
	// }
}
