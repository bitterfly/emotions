package emotions

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type EegClusterable struct {
	Class string      `json:"class"`
	Data  [][]float64 `json:"data"` //19x4
}

var waveRanges = [][2]float64{
	[2]float64{4.0, 8.0},   //Θ
	[2]float64{8.0, 12.0},  // α
	[2]float64{12.0, 30.0}, // β
	[2]float64{30.0, 50.0}, // γ
}

func getRange(n float64) int {
	if n > waveRanges[len(waveRanges)-1][1] {
		return -1
	}

	if n < waveRanges[0][0] {
		return -1
	}

	for i, w := range waveRanges {
		if n > w[0] && n < w[1] {
			return i
		}
	}
	return -1
}

func getVector(line []string) []float64 {
	floatValues := make([]float64, 0, len(line))

	for _, s := range line {
		if strings.TrimSpace(s) == "" {
			continue
		}

		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			panic(err)
		}

		floatValues = append(floatValues, v)
	}
	return floatValues
}

// ReadXML takes an xml file with eeg readings and returns a vector for each electrode in time
// where the first coordinate is the data from the first electrode and so on
func ReadXML(filename string, elNum int) [][]float64 {
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}

	electrodes := make([][]float64, elNum, elNum)
	scanner := csv.NewReader(file)

	for {
		line, err := scanner.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				panic(err)
			}
		}

		for i, value := range getVector(line) {
			electrodes[i] = append(electrodes[i], value)
		}
	}

	return electrodes
}

func cutElectrodeIntoFrames(electrode []float64) [][]float64 {
	return CutSliceIntoFrames(electrode, 500, 200, 150)
}

func fourierElectrode(frames [][]float64) [][]Complex {
	fouriers := make([][]Complex, len(frames), len(frames))
	for i := 0; i < len(frames); i++ {
		fouriers[i], _ = FftReal(frames[i])
	}

	return fouriers
}

func getWavesMean(coefficients [][]Complex) []float64 {
	means := make([]float64, len(waveRanges), len(waveRanges))
	for i := 0; i < len(coefficients); i++ {
		for j := 0; j < len(coefficients[0]); j++ {
			magnitude := Magnitude(coefficients[i][j])
			w := getRange(magnitude)
			if w == -1 {
				break
			}
			means[w] += magnitude
		}
	}

	divide(&means, float64(len(coefficients)))
	return means
}

func getElectrodeWavesDistribution(electrodeData []float64) []float64 {
	frames := cutElectrodeIntoFrames(electrodeData)
	fouriers := fourierElectrode(frames)
	return getWavesMean(fouriers)
}

// GetFeatureVector returns the mean of Θ, α, β and γ waves for each of the given elNum electrodes
// returns a vector 19x4
func GetFeatureVector(filename string, elNum int) [][]float64 {
	data := ReadXML(filename, elNum)
	features := make([][]float64, len(data), len(data))
	for i, d := range data {
		features[i] = getElectrodeWavesDistribution(d)
	}

	return features
}

func getFeaturesFromFiles(filenames []string) []EegClusterable {
	trainingSet := make([]EegClusterable, len(filenames), len(filenames))

	for i, file := range filenames {
		filename := filepath.Base(file)
		name := filename[0 : len(filename)-len(filepath.Ext(filename))]
		newFeatures := GetFeatureVector(file, 19)
		trainingSet[i] = EegClusterable{
			Class: name,
			Data:  newFeatures,
		}
	}

	return trainingSet
}

func SaveEegTrainingSet(filenames []string, outputFilename string) {
	features := getFeaturesFromFiles(filenames)

	bytes, err := json.Marshal(features)
	if err != nil {
		panic(err)
	}

	ioutil.WriteFile(outputFilename, bytes, 0644)
}

func getEegTrainingSet(filename string) []EegClusterable {
	var clusterables []EegClusterable
	bytes, _ := ioutil.ReadFile(filename)
	err := json.Unmarshal(bytes, &clusterables)
	if err != nil {
		panic(err)
	}

	return clusterables
}
