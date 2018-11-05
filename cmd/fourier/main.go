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
	filename := os.Args[1]
	wf, err := fourier.Read(filename)
	if err != nil {
		panic(err.Error)
	}

	coefficients := fourier.FftWav(wf)

	v := make(plotter.Values, len(coefficients))
	for i := range v {
		if fourier.Magnitude(coefficients[i]) > fourier.EPS {
			fmt.Printf("%d %s\n", i, coefficients[i])
		}
		v[i] = fourier.Magnitude(coefficients[i])
	}

	plotc, err := plot.New()
	if err != nil {
		panic(err)
	}
	plots, err := plot.New()
	if err != nil {
		panic(err)
	}

	plotc.X.Min = 0
	plotc.X.Max = float64(len(coefficients))
	plotc.X.Label.Text = "Frequency"
	plotc.Y.Label.Text = "Energy"
	bars, err := plotter.NewBarChart(v, 2)
	bars.Color = color.RGBA{10, 120, 120, 1}
	plotc.Add(bars)

	if err := plotc.Save(16*vg.Inch, 16*vg.Inch, "spectre.png"); err != nil {
		panic(err)
	}

	data := wf.GetData()
	s := make(plotter.XYs, len(data))
	for i := 0; i < len(data); i++ {
		s[i].X = float64(i)
		s[i].Y = data[i]
	}

	line, _ := plotter.NewLine(s)
	line.Color = color.RGBA{254, 1, 2, 1}

	plots.Add(line)

	if err := plots.Save(4*vg.Inch, 4*vg.Inch, "signal.png"); err != nil {
		panic(err)
	}
}
