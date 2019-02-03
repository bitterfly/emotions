package emotions

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/palette"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
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

	s := make(plotter.XYs, len(data))
	for i := 0; i < len(data); i++ {
		s[i].X = float64(i)
		s[i].Y = data[i]
	}

	line, _ := plotter.NewLine(s)
	line.Color = color.RGBA{0, 232, 88, 255}

	plots.Add(line)

	if err := plots.Save(32*vg.Inch, 16*vg.Inch, file); err != nil {
		panic(err)
	}
}

func PlotClusters(data []MfccClusterisable, k int, file string) {
	plots, err := plot.New()
	if err != nil {
		panic(err)
	}

	s := make(plotter.XYs, len(data))
	clusters := make([]string, len(data))

	for i, d := range data {
		s[i].X = d.coefficients[0]
		s[i].Y = d.coefficients[1]
		clusters[i] = fmt.Sprintf("%d", d.clusterID)
	}

	scatter, _ := plotter.NewScatter(s)
	palette := palette.Heat(k, 1)

	scatter.GlyphStyleFunc = func(i int) draw.GlyphStyle {
		return draw.GlyphStyle{Color: palette.Colors()[data[i].clusterID], Radius: vg.Points(3), Shape: draw.CircleGlyph{}}
	}

	plots.Add(scatter)

	if err := plots.Save(32*vg.Inch, 16*vg.Inch, file); err != nil {
		panic(err)
	}
}

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
		c, _ := FftReal(frame)
		PlotCoefficients(c, fmt.Sprintf("spectrum/spectrum%d.png", i))

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

func PlotEeg(filename string, output string) {
	ts := getEegTrainingSet(filename)
	plotEeg(ts, output)
}

func plotEeg(data []EegClusterable, file string) {
	fmt.Printf("data: %d %d %d\n", len(data), len(data[0].Data), len(data[0].Data[0]))

	plots, err := plot.New()
	if err != nil {
		panic(err)
	}

	maximums := make([]float64, len(data[0].Data), len(data[0].Data))
	for i := 0; i < len(data); i++ {
		for j := 0; j < len(data[i].Data); j++ {
			if data[i].Data[j][0] > maximums[0] {
				maximums[0] = data[i].Data[j][0]
			}

			if data[i].Data[j][1] > maximums[1] {
				maximums[1] = data[i].Data[j][1]
			}

			if data[i].Data[j][2] > maximums[2] {
				maximums[2] = data[i].Data[j][2]
			}

			if data[i].Data[j][3] > maximums[3] {
				maximums[3] = data[i].Data[j][3]
			}
		}
	}

	fmt.Printf("len: %d\n", len(data))

	for i := 0; i < len(data); i++ {
		s := make(plotter.XYs, len(data[i].Data), len(data[i].Data))
		for j := 0; j < len(data[i].Data); j++ {
			s[j].X = float64(j)
			s[j].Y = float64(i)
		}

		scatter, _ := plotter.NewScatter(s)
		bla := data[i]
		scatter.GlyphStyleFunc = func(k int) draw.GlyphStyle {
			return draw.GlyphStyle{
				Color:  getColour(bla.Data[k], maximums),
				Radius: vg.Points(50),
				Shape:  draw.BoxGlyph{},
			}
		}
		plots.Add(scatter)
	}

	err = plots.Save(vg.Length(2.0*50*float64(len(data[0].Data))), vg.Length(2*50*float64(len(data))), file)
	if err != nil {
		panic(err)
	}
}

func PlotEmotion(filename string, output string) {
	data := ReadXML(filename, 19)

	var cl []EegClusterable

	for i, d := range data {
		var features [][]float64
		frames := cutElectrodeIntoFrames(d)
		fouriers := fourierElectrode(frames)
		for _, f := range fouriers {
			v := make([]float64, 4, 4)
			for _, ff := range f {
				magnitude := Magnitude(ff)
				w := getRange(magnitude)
				if w == -1 {
					break
				}
				v[w] = magnitude
			}

			features = append(features, v)
		}

		for j := 0; j < len(features); j++ {
			if j > len(cl)-1 {
				cl = append(cl, EegClusterable{
					Data: make([][]float64, len(data), len(data)),
				})
			}

			cl[j].Data[i] = features[j]
		}
	}

	plotEeg(cl, output)
}

func getColour(x []float64, maximums []float64) color.RGBA {
	r := uint8(math.Min(math.Floor(x[0]*float64(255)/maximums[0]), 255))
	g := uint8(math.Min(math.Floor(x[1]*float64(255)/maximums[1]), 255))
	b := uint8(math.Min(math.Floor(x[2]*float64(255)/maximums[2]), 255))
	a := uint8(math.Min(math.Floor(x[3]*float64(155)/maximums[3]+100), 255))

	// if r == 255 || g == 255 || b == 255 || a == 255 {
	// fmt.Printf("x: %v, r: %d g: %d b: %d a: %d\n", x, r, g, b, a)
	// }

	return color.RGBA{
		R: r,
		G: g,
		B: b,
		A: a,
	}
}

func getName(s string) string {
	if rune(s[1]) == '-' {
		return s[2:]
	}

	switch rune(s[1]) {
	case 'h':
		return "happiness"
	case 's':
		return "sadness"
	case 'a':
		return "anger"
	case 'n':
		return "neutral"
	default:
		panic(s)
	}
}

func GetValence(s string, dither float64) float64 {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	d := (r1.Float64() - 0.5) * dither

	switch s {
	case "happiness":
		return 1 + d
	case "sadness":
		return -1 + d
	case "anger":
		return -1 + d
	case "neutral":
		return 0 + d
	default:
		panic(s)
	}
}

//ParseArguments receives command arguments and separates them on spaces
func ParseArguments(args []string) map[string][]string {
	var arguments []string
	var emotion string
	emotions := make(map[string][]string)

	for i := 0; i < len(args); i++ {
		if rune(args[i][0]) == '-' {
			if len(arguments) != 0 {
				emotions[emotion] = arguments
				arguments = []string{}
			}
			emotion = getName(args[i])
		} else {
			arguments = append(arguments, args[i])
		}
	}

	emotions[emotion] = arguments

	return emotions
}
