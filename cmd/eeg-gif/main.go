package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"os/exec"
	"path"

	"github.com/bitterfly/emotions/emotions"
	"github.com/lucasb-eyer/go-colorful"
)

type GradientTable []struct {
	Col colorful.Color
	Pos float64
}

type EEGImage struct {
	colours   []float64
	size      image.Rectangle
	positions [][2]float64
}

func (e *EEGImage) At(x, y int) color.Color {
	return getColour(float64(x)/float64(e.size.Max.X), float64(y)/float64(e.size.Max.Y), e.colours, e.positions[:])
}

func (e *EEGImage) Bounds() image.Rectangle {
	return e.size
}

func (e *EEGImage) ColorModel() color.Model {
	return color.RGBAModel
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getData(bucketSize int, frameLen int, frameStep int, files []string) [][]float64 {
	numElectrodes := 19
	numWaves := 4
	data := make([][]float64, 0, numElectrodes*numWaves)

	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		average := emotions.GetAverage(bucketSize, frameLen, len(current))

		data = append(data, emotions.AverageSlice(current, average)...)
	}

	return data
}

func getMinMax(data []float64, minimums []float64, maximums []float64) {
	for i := 0; i < len(data); i++ {
		if data[i] < minimums[i] {
			minimums[i] = data[i]
		}

		if data[i] > maximums[i] {
			maximums[i] = data[i]
		}
	}
}

func MustParseHex(s string) colorful.Color {
	c, err := colorful.Hex(s)
	if err != nil {
		panic("MustParseHex: " + err.Error())
	}
	return c
}

func (self GradientTable) GetInterpolatedColorFor(t float64) colorful.Color {
	for i := 0; i < len(self)-1; i++ {
		c1 := self[i]
		c2 := self[i+1]
		if c1.Pos <= t && t <= c2.Pos {
			// We are in between c1 and c2. Go blend them!
			t := (t - c1.Pos) / (c2.Pos - c1.Pos)
			return c1.Col.BlendHcl(c2.Col, t).Clamped()
		}
	}

	// Nothing found? Means we're at (or past) the last gradient keypoint.
	return self[len(self)-1].Col
}

func ramp(x float64) color.Color {
	// gr := GradientTable{
	// 	{MustParseHex("#FF1E38"), 0.0},
	// 	{MustParseHex("#5CC549"), 0.5},
	// 	{MustParseHex("#29357B"), 1.0},
	// }
	gr := GradientTable{
		{MustParseHex("#FFFFFF"), 0.0},
		{MustParseHex("#000000"), 1.0},
	}

	return &image.Uniform{gr.GetInterpolatedColorFor(x)}
}

func dist(x1, y1, x2, y2 float64) (float64, bool) {
	d := (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
	return 1 / math.Pow(d, 2), d < 0.00001
}

func distToPositions(x, y float64, positions [][2]float64) (int, []float64, float64) {
	distances := make([]float64, len(positions), len(positions))
	sum := 0.0
	for i := 0; i < len(positions); i++ {
		d, underflow := dist(x, y, positions[i][0], positions[i][1])
		sum += d
		distances[i] = d
		if underflow {
			return i, []float64{}, 0
		}
	}
	return -1, distances, sum
}

func getColour(x float64, y float64, colours []float64, positions [][2]float64) color.Color {
	colour := 0.0
	underflow, distances, sum := distToPositions(x, y, positions)
	if underflow != -1 {
		return ramp(colours[underflow])
	}
	if (x-0.5)*(x-0.5)+(y-0.5)*(y-0.5) > 0.45*0.45 {
		return color.RGBA{R: 0, G: 0, B: 0, A: 0}
	}

	for i := 0; i < len(distances); i++ {
		colour += colours[i] * distances[i] / sum
	}

	return ramp(colour)
}

func drawEeg(colours []float64, outputDir string) {
	positions := [19][2]float64{
		[2]float64{199.0 / 490.0, 119.0 / 490.0},
		[2]float64{297.0 / 490.0, 119.0 / 490.0},
		[2]float64{187.0 / 490.0, 183.0 / 490.0},
		[2]float64{303.0 / 490.0, 183.0 / 490.0},
		[2]float64{245.0 / 490.0, 183.0 / 490.0},
		[2]float64{173.0 / 490.0, 256.0 / 490.0},
		[2]float64{317.0 / 490.0, 256.0 / 490.0},
		[2]float64{245.0 / 490.0, 256.0 / 490.0},
		[2]float64{187.0 / 490.0, 328.0 / 490.0},
		[2]float64{303.0 / 490.0, 328.0 / 490.0},
		[2]float64{245.0 / 490.0, 328.0 / 490.0},
		[2]float64{200.0 / 490.0, 393.0 / 490.0},
		[2]float64{129.0 / 490.0, 172.0 / 490.0},
		[2]float64{290.0 / 490.0, 393.0 / 490.0},
		[2]float64{361.0 / 490.0, 172.0 / 490.0},
		[2]float64{101.0 / 490.0, 256.0 / 490.0},
		[2]float64{389.0 / 490.0, 256.0 / 490.0},
		[2]float64{129.0 / 490.0, 340.0 / 490.0},
		[2]float64{361.0 / 490.0, 340.0 / 490.0},
	}

	size := image.Rect(0, 0, 1024, 1024)
	canvas := &EEGImage{colours: colours, positions: positions[:], size: size}

	fd, err := os.Create(outputDir)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	err = png.Encode(fd, canvas)
	if err != nil {
		log.Fatal(err)
	}
}

func normaliseByWaves(data []float64) {
	waveMin := []float64{data[0], data[0], data[0], data[0]}
	waveMax := []float64{data[0], data[0], data[0], data[0]}

	updateMinMax(data, waveMin, waveMax)

	for i := 0; i < len(data); i++ {
		data[i] = (data[i] - waveMin[i%4]) / (waveMax[i%4] - waveMin[i%4])
	}
}

func min(a, b, c float64) float64 {
	if a < b {
		if c < a {
			return c
		}
		return a
	} else {
		if c < b {
			return c
		}
		return b
	}
}

func minMax(d ...float64) (float64, float64) {
	min := d[0]
	max := d[0]
	for i := 1; i < len(d); i++ {
		if d[i] < min {
			min = d[i]
		}
		if d[i] > max {
			max = d[i]
		}
	}
	return min, max
}

func updateMinMax(data []float64, min []float64, max []float64) {
	for i := 0; i < len(data); i++ {
		if data[i] < min[i%4] {
			min[i%4] = data[i]
		}
		if data[i] > max[i%4] {
			max[i%4] = data[i]
		}
	}
}

func drawEmotion(emotion string, data []float64, outputDir string, i int) {
	fmt.Printf("Emotion: %s\n", emotion)

	j := 0
	waves := make([][]float64, 4, 4)
	for i := 0; i < 4; i++ {
		waves[i] = make([]float64, 19, 19)
	}

	for i := 0; i < len(data); i += 4 {
		waves[0][j] = data[i]
		waves[1][j] = data[i+1]
		waves[2][j] = data[i+2]
		waves[3][j] = data[i+3]
		j++
	}

	drawEeg(waves[0], path.Join(outputDir, fmt.Sprintf("%s_%05d.png", "δ", i)))
	drawEeg(waves[1], path.Join(outputDir, fmt.Sprintf("%s_%05d.png", "α", i)))
	drawEeg(waves[2], path.Join(outputDir, fmt.Sprintf("%s_%05d.png", "β", i)))
	drawEeg(waves[3], path.Join(outputDir, fmt.Sprintf("%s_%05d.png", "γ", i)))
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <ouput-file> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	outputDir := os.Args[1]
	arguments, _, _, err := emotions.ParseArgumentsFromFile(os.Args[2], false)

	if err != nil {
		panic(err)
	}

	var emotion string
	for k := range arguments {
		emotion = k
		break
	}

	fmt.Printf("%s\n%s\n", emotion, outputDir)
	frameLen := 200
	frameStep := 150

	data := getData(0, frameLen, frameStep, arguments[emotion])

	for i := 0; i < len(data); i++ {
		normaliseByWaves(data[i])
		drawEmotion(emotion, data[i], outputDir, i+1)
	}
	err = exec.Command("./blendit.sh", outputDir).Run()
	if err != nil {
		panic(err)
	}
}
