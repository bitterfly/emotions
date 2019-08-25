package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}

func getData(featureType string, bucketSize int, frameLen int, frameStep int, files []string, tag string) emotions.Tagged {
	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		average := emotions.GetAverage(bucketSize, frameLen, len(current))
		if featureType == "de" {
			data = append(data, emotions.GetDE(emotions.AverageSlice(current, average))...)
		} else {
			data = append(data, emotions.AverageSlice(current, average)...)
		}
	}

	return emotions.Tagged{
		Tag:  tag,
		Data: data,
	}
}

func KNN(arguments map[string][]string, featureType string, outputDirname string, bucketSize int, frameLen int, frameStep int) error {
	if fileExists(outputDirname) {
		os.Remove(outputDirname)
	}

	data := make([]emotions.Tagged, 0, 10)
	for e, v := range arguments {
		data = append(data, getData(featureType, bucketSize, frameLen, frameStep, v, e))
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("could not marshal eeg data: %s", err.Error())
	}
	err = ioutil.WriteFile(path.Join(outputDirname, "emotion.knn"), bytes, 0664)

	return err
}

func marshallGMM(outputDirname string, emotion string, data [][]float64, k int) error {
	egm := emotions.EmotionGausianMixure{
		Emotion: emotion,
		GM:      emotions.GMM(data, k),
	}

	bytes, err := json.Marshal(egm)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, k))
	}

	return ioutil.WriteFile(outputDirname, bytes, 0644)
}

func GMM(arguments map[string][]string, featureType string, outputDir string, k int, bucketSize int, frameLen int, frameStep int) error {
	data := make(map[string]emotions.Tagged)
	for e, v := range arguments {
		data[e] = getData(featureType, bucketSize, frameLen, frameStep, v, e)
	}

	if !fileExists(outputDir) {
		os.Mkdir(outputDir, 0775)
	}

	for e, v := range data {
		err := marshallGMM(path.Join(outputDir, fmt.Sprintf("%s.gmm", e)), e, v.Data, k)
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	if len(os.Args) < 5 {
		panic("go run main.go <classifier-type> <feature-type> <bucket-size> <output-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	classifierType := os.Args[1]

	featureType := os.Args[2]

	bucketSize, err := strconv.Atoi(os.Args[3])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[3]))
	}

	outputDir := os.Args[4]
	_, arguments, _, err := emotions.ParseArgumentsFromFile(os.Args[5], false)

	if err != nil {
		panic(err)
	}

	frameLen := 200
	frameStep := 150
	k := 3

	switch classifierType {
	case "knn":
		err = KNN(arguments, featureType, outputDir, bucketSize, frameLen, frameStep)
		if err != nil {
			panic(err)
		}
	case "gmm":
		err = GMM(arguments, featureType, outputDir, k, bucketSize, frameLen, frameStep)
		if err != nil {
			panic(err)
		}
	default:
		panic(fmt.Sprintf("Unknown classifier %s", classifierType))
	}

}
