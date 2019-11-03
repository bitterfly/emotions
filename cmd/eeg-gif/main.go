package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
	"path"

	"github.com/bitterfly/emotions/emotions"
	"github.com/lucasb-eyer/go-colorful"
)

var googleColours [256][3]uint8 = [256][3]uint8{
	[3]uint8{48, 18, 59},
	[3]uint8{50, 21, 67},
	[3]uint8{51, 24, 74},
	[3]uint8{52, 27, 81},
	[3]uint8{53, 30, 88},
	[3]uint8{54, 33, 95},
	[3]uint8{55, 36, 102},
	[3]uint8{56, 39, 109},
	[3]uint8{57, 42, 115},
	[3]uint8{58, 45, 121},
	[3]uint8{59, 47, 128},
	[3]uint8{60, 50, 134},
	[3]uint8{61, 53, 139},
	[3]uint8{62, 56, 145},
	[3]uint8{63, 59, 151},
	[3]uint8{63, 62, 156},
	[3]uint8{64, 64, 162},
	[3]uint8{65, 67, 167},
	[3]uint8{65, 70, 172},
	[3]uint8{66, 73, 177},
	[3]uint8{66, 75, 181},
	[3]uint8{67, 78, 186},
	[3]uint8{68, 81, 191},
	[3]uint8{68, 84, 195},
	[3]uint8{68, 86, 199},
	[3]uint8{69, 89, 203},
	[3]uint8{69, 92, 207},
	[3]uint8{69, 94, 211},
	[3]uint8{70, 97, 214},
	[3]uint8{70, 100, 218},
	[3]uint8{70, 102, 221},
	[3]uint8{70, 105, 224},
	[3]uint8{70, 107, 227},
	[3]uint8{71, 110, 230},
	[3]uint8{71, 113, 233},
	[3]uint8{71, 115, 235},
	[3]uint8{71, 118, 238},
	[3]uint8{71, 120, 240},
	[3]uint8{71, 123, 242},
	[3]uint8{70, 125, 244},
	[3]uint8{70, 128, 246},
	[3]uint8{70, 130, 248},
	[3]uint8{70, 133, 250},
	[3]uint8{70, 135, 251},
	[3]uint8{69, 138, 252},
	[3]uint8{69, 140, 253},
	[3]uint8{68, 143, 254},
	[3]uint8{67, 145, 254},
	[3]uint8{66, 148, 255},
	[3]uint8{65, 150, 255},
	[3]uint8{64, 153, 255},
	[3]uint8{62, 155, 254},
	[3]uint8{61, 158, 254},
	[3]uint8{59, 160, 253},
	[3]uint8{58, 163, 252},
	[3]uint8{56, 165, 251},
	[3]uint8{55, 168, 250},
	[3]uint8{53, 171, 248},
	[3]uint8{51, 173, 247},
	[3]uint8{49, 175, 245},
	[3]uint8{47, 178, 244},
	[3]uint8{46, 180, 242},
	[3]uint8{44, 183, 240},
	[3]uint8{42, 185, 238},
	[3]uint8{40, 188, 235},
	[3]uint8{39, 190, 233},
	[3]uint8{37, 192, 231},
	[3]uint8{35, 195, 228},
	[3]uint8{34, 197, 226},
	[3]uint8{32, 199, 223},
	[3]uint8{31, 201, 221},
	[3]uint8{30, 203, 218},
	[3]uint8{28, 205, 216},
	[3]uint8{27, 208, 213},
	[3]uint8{26, 210, 210},
	[3]uint8{26, 212, 208},
	[3]uint8{25, 213, 205},
	[3]uint8{24, 215, 202},
	[3]uint8{24, 217, 200},
	[3]uint8{24, 219, 197},
	[3]uint8{24, 221, 194},
	[3]uint8{24, 222, 192},
	[3]uint8{24, 224, 189},
	[3]uint8{25, 226, 187},
	[3]uint8{25, 227, 185},
	[3]uint8{26, 228, 182},
	[3]uint8{28, 230, 180},
	[3]uint8{29, 231, 178},
	[3]uint8{31, 233, 175},
	[3]uint8{32, 234, 172},
	[3]uint8{34, 235, 170},
	[3]uint8{37, 236, 167},
	[3]uint8{39, 238, 164},
	[3]uint8{42, 239, 161},
	[3]uint8{44, 240, 158},
	[3]uint8{47, 241, 155},
	[3]uint8{50, 242, 152},
	[3]uint8{53, 243, 148},
	[3]uint8{56, 244, 145},
	[3]uint8{60, 245, 142},
	[3]uint8{63, 246, 138},
	[3]uint8{67, 247, 135},
	[3]uint8{70, 248, 132},
	[3]uint8{74, 248, 128},
	[3]uint8{78, 249, 125},
	[3]uint8{82, 250, 122},
	[3]uint8{85, 250, 118},
	[3]uint8{89, 251, 115},
	[3]uint8{93, 252, 111},
	[3]uint8{97, 252, 108},
	[3]uint8{101, 253, 105},
	[3]uint8{105, 253, 102},
	[3]uint8{109, 254, 98},
	[3]uint8{113, 254, 95},
	[3]uint8{117, 254, 92},
	[3]uint8{121, 254, 89},
	[3]uint8{125, 255, 86},
	[3]uint8{128, 255, 83},
	[3]uint8{132, 255, 81},
	[3]uint8{136, 255, 78},
	[3]uint8{139, 255, 75},
	[3]uint8{143, 255, 73},
	[3]uint8{146, 255, 71},
	[3]uint8{150, 254, 68},
	[3]uint8{153, 254, 66},
	[3]uint8{156, 254, 64},
	[3]uint8{159, 253, 63},
	[3]uint8{161, 253, 61},
	[3]uint8{164, 252, 60},
	[3]uint8{167, 252, 58},
	[3]uint8{169, 251, 57},
	[3]uint8{172, 251, 56},
	[3]uint8{175, 250, 55},
	[3]uint8{177, 249, 54},
	[3]uint8{180, 248, 54},
	[3]uint8{183, 247, 53},
	[3]uint8{185, 246, 53},
	[3]uint8{188, 245, 52},
	[3]uint8{190, 244, 52},
	[3]uint8{193, 243, 52},
	[3]uint8{195, 241, 52},
	[3]uint8{198, 240, 52},
	[3]uint8{200, 239, 52},
	[3]uint8{203, 237, 52},
	[3]uint8{205, 236, 52},
	[3]uint8{208, 234, 52},
	[3]uint8{210, 233, 53},
	[3]uint8{212, 231, 53},
	[3]uint8{215, 229, 53},
	[3]uint8{217, 228, 54},
	[3]uint8{219, 226, 54},
	[3]uint8{221, 224, 55},
	[3]uint8{223, 223, 55},
	[3]uint8{225, 221, 55},
	[3]uint8{227, 219, 56},
	[3]uint8{229, 217, 56},
	[3]uint8{231, 215, 57},
	[3]uint8{233, 213, 57},
	[3]uint8{235, 211, 57},
	[3]uint8{236, 209, 58},
	[3]uint8{238, 207, 58},
	[3]uint8{239, 205, 58},
	[3]uint8{241, 203, 58},
	[3]uint8{242, 201, 58},
	[3]uint8{244, 199, 58},
	[3]uint8{245, 197, 58},
	[3]uint8{246, 195, 58},
	[3]uint8{247, 193, 58},
	[3]uint8{248, 190, 57},
	[3]uint8{249, 188, 57},
	[3]uint8{250, 186, 57},
	[3]uint8{251, 184, 56},
	[3]uint8{251, 182, 55},
	[3]uint8{252, 179, 54},
	[3]uint8{252, 177, 54},
	[3]uint8{253, 174, 53},
	[3]uint8{253, 172, 52},
	[3]uint8{254, 169, 51},
	[3]uint8{254, 167, 50},
	[3]uint8{254, 164, 49},
	[3]uint8{254, 161, 48},
	[3]uint8{254, 158, 47},
	[3]uint8{254, 155, 45},
	[3]uint8{254, 153, 44},
	[3]uint8{254, 150, 43},
	[3]uint8{254, 147, 42},
	[3]uint8{254, 144, 41},
	[3]uint8{253, 141, 39},
	[3]uint8{253, 138, 38},
	[3]uint8{252, 135, 37},
	[3]uint8{252, 132, 35},
	[3]uint8{251, 129, 34},
	[3]uint8{251, 126, 33},
	[3]uint8{250, 123, 31},
	[3]uint8{249, 120, 30},
	[3]uint8{249, 117, 29},
	[3]uint8{248, 114, 28},
	[3]uint8{247, 111, 26},
	[3]uint8{246, 108, 25},
	[3]uint8{245, 105, 24},
	[3]uint8{244, 102, 23},
	[3]uint8{243, 99, 21},
	[3]uint8{242, 96, 20},
	[3]uint8{241, 93, 19},
	[3]uint8{240, 91, 18},
	[3]uint8{239, 88, 17},
	[3]uint8{237, 85, 16},
	[3]uint8{236, 83, 15},
	[3]uint8{235, 80, 14},
	[3]uint8{234, 78, 13},
	[3]uint8{232, 75, 12},
	[3]uint8{231, 73, 12},
	[3]uint8{229, 71, 11},
	[3]uint8{228, 69, 10},
	[3]uint8{226, 67, 10},
	[3]uint8{225, 65, 9},
	[3]uint8{223, 63, 8},
	[3]uint8{221, 61, 8},
	[3]uint8{220, 59, 7},
	[3]uint8{218, 57, 7},
	[3]uint8{216, 55, 6},
	[3]uint8{214, 53, 6},
	[3]uint8{212, 51, 5},
	[3]uint8{210, 49, 5},
	[3]uint8{208, 47, 5},
	[3]uint8{206, 45, 4},
	[3]uint8{204, 43, 4},
	[3]uint8{202, 42, 4},
	[3]uint8{200, 40, 3},
	[3]uint8{197, 38, 3},
	[3]uint8{195, 37, 3},
	[3]uint8{193, 35, 2},
	[3]uint8{190, 33, 2},
	[3]uint8{188, 32, 2},
	[3]uint8{185, 30, 2},
	[3]uint8{183, 29, 2},
	[3]uint8{180, 27, 1},
	[3]uint8{178, 26, 1},
	[3]uint8{175, 24, 1},
	[3]uint8{172, 23, 1},
	[3]uint8{169, 22, 1},
	[3]uint8{167, 20, 1},
	[3]uint8{164, 19, 1},
	[3]uint8{161, 18, 1},
	[3]uint8{158, 16, 1},
	[3]uint8{155, 15, 1},
	[3]uint8{152, 14, 1},
	[3]uint8{149, 13, 1},
	[3]uint8{146, 11, 1},
	[3]uint8{142, 10, 1},
	[3]uint8{139, 9, 2},
	[3]uint8{136, 8, 2},
	[3]uint8{133, 7, 2},
	[3]uint8{129, 6, 2},
	[3]uint8{126, 5, 2},
	[3]uint8{122, 4, 3},
}

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
	c := googleColours[int((x-0.0001)*256.0)]
	return &color.RGBA{R: c[0], G: c[1], B: c[2], A: 255}
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
	foo, arguments, _, err := emotions.ParseArgumentsFromFile(os.Args[2], true)

	fmt.Printf("%v\n", foo)

	if err != nil {
		panic(err)
	}

	var emotion string
	for k := range arguments {
		emotion = k
		break
	}

	frameLen := 200
	frameStep := 150

	data := getData(0, frameLen, frameStep, arguments[emotion])
	for i := 0; i < len(data); i++ {
		normaliseByWaves(data[i])
		drawEmotion(emotion, data[i], outputDir, i+1)
	}
	// err = exec.Command("./blendit.sh", outputDir).Run()
	// if err != nil {
	// 	panic(err)
	// }
}
