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

func averageSlice(data [][]float64) []float64 {
	result := make([]float64, len(data[0]), len(data[0]))
	for j := 0; j < len(data[0]); j++ {
		for i := 0; i < len(data); i++ {
			result[j] += data[i][j]
		}
		result[j] /= float64(len(data))
	}
	return result
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
		if data[i] < minimums[i%4] {
			minimums[i%4] = data[i]
		}

		if data[i] > maximums[i%4] {
			maximums[i%4] = data[i]
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

func drawEeg(colours []float64, outputFile string) {
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

	fd, err := os.Create(outputFile)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	err = png.Encode(fd, canvas)
	if err != nil {
		log.Fatal(err)
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

// ei = 19x4
func globalMinMax(e1 []float64, e2 []float64, e3 []float64) ([]float64, []float64) {
	waveMin := []float64{e1[0], e1[1], e1[2], e1[3]}
	waveMax := []float64{e1[0], e1[1], e1[2], e1[3]}

	updateMinMax(e1, waveMin, waveMax)
	updateMinMax(e2, waveMin, waveMax)
	updateMinMax(e3, waveMin, waveMax)

	return waveMin, waveMax
}

func drawEmotion(emotion string, data []float64, outputDir string) {
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

	drawEeg(waves[0], path.Join(outputDir, fmt.Sprintf("%s_%s", emotion, "δ.png")))
	drawEeg(waves[1], path.Join(outputDir, fmt.Sprintf("%s_%s", emotion, "α.png")))
	drawEeg(waves[2], path.Join(outputDir, fmt.Sprintf("%s_%s", emotion, "β.png")))
	drawEeg(waves[3], path.Join(outputDir, fmt.Sprintf("%s_%s", emotion, "γ.png")))
}

func normalize(data, min, max []float64) {
	for i := 0; i < len(data); i++ {
		data[i] = (data[i] - min[i%4]) / (max[i%4] - min[i%4])
	}
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <ouput-dir> [<single-file>] <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	outputDir := os.Args[1]
	file := os.Args[2]

	arguments := make(map[string][]string)
	var err error
	if len(os.Args) > 3 {
		_, arguments,_ , err = emotions.ParseArgumentsFromFile(os.Args[3], true)
		if err != nil {
			panic(err)
		}
	}
	frameLen := 200
	frameStep := 150

	eegPositive, _ := arguments["eeg-positive"]
	eegNegative, _ := arguments["eeg-negative"]
	eegNeutral, ok := arguments["eeg-neutral"]
	var min, max []float64
	data := getData(1, frameLen, frameStep, []string{file})[0]

	min = []float64{data[0], data[1], data[2], data[3]}
	max = []float64{data[0], data[1], data[2], data[3]}
	if !ok {
		getMinMax(data, min, max)
	} else {
		globalData := make([][]float64, 0, 100)

		globalData = append(globalData, getData(1, frameLen, frameStep, eegPositive)...)
		globalData = append(globalData, getData(1, frameLen, frameStep, eegNegative)...)
		globalData = append(globalData, getData(1, frameLen, frameStep, eegNeutral)...)

		for i := 0; i < len(globalData); i++ {
			updateMinMax(globalData[i], min, max)
		}
	}
	normalize(data, min, max)
	drawEmotion("data", data, outputDir)

	err = exec.Command("./blendit.sh", outputDir).Run()
	if err != nil {
		panic(err)
	}
}
