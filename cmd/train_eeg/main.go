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

func getData(bucketSize int, frameLen int, frameStep int, files []string, tag string) emotions.Tagged {
	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		// fmt.Printf("File: %s current: %v\n", files[i], current)
		average := emotions.GetAverage(bucketSize, frameLen, len(current))
		data = append(data, emotions.AverageSlice(current, average)...)
	}

	return emotions.Tagged{
		Tag:  tag,
		Data: data,
	}
}

func marshallToFile(bucketSize int, frameLen int, frameStep int, filename string, positiveFiles []string, negativeFiles []string, neutralFiles []string) error {
	bytes, err := json.Marshal([]emotions.Tagged{
		getData(bucketSize, frameLen, frameStep, positiveFiles, "eeg-positive"),
		getData(bucketSize, frameLen, frameStep, negativeFiles, "eeg-negative"),
		getData(bucketSize, frameLen, frameStep, neutralFiles, "eeg-neutral"),
	})
	if err != nil {
		return fmt.Errorf("could not marshal eeg data: %s", err.Error())
	}

	err = ioutil.WriteFile(filename, bytes, 0664)
	return err
}

func KNN(eegPositive []string, eegNegative []string, eegNeutral []string, outputFile string, bucketSize int, frameLen int, frameStep int) error {
	if fileExists(outputFile) {
		os.Remove(outputFile)
	}

	err := marshallToFile(bucketSize, frameLen, frameStep, outputFile, eegPositive, eegNegative, eegNeutral)
	return err
}

func marshallGMM(outputFilename string, emotion string, data [][]float64, k int) error {
	egm := emotions.EmotionGausianMixure{
		Emotion: emotion,
		GM:      emotions.GMM(data, k),
	}

	bytes, err := json.Marshal(egm)
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, k))
	}

	return ioutil.WriteFile(outputFilename, bytes, 0644)
}

func GMM(eegPositive []string, eegNegative []string, eegNeutral []string, outputFile string, k int, bucketSize int, frameLen int, frameStep int) error {
	positive := getData(bucketSize, frameLen, frameStep, eegPositive, "eeg-positive")
	negative := getData(bucketSize, frameLen, frameStep, eegNegative, "eeg-negative")
	neutral := getData(bucketSize, frameLen, frameStep, eegNeutral, "eeg-neutral")

	if !fileExists(outputFile) {
		os.Mkdir(outputFile, 0775)
	}

	err := marshallGMM(path.Join(outputFile, "positive.gmm"), "eeg-positive", positive.Data, k)
	if err != nil {
		return err
	}
	err = marshallGMM(path.Join(outputFile, "negative.gmm"), "eeg-negative", negative.Data, k)
	if err != nil {
		return err
	}
	return marshallGMM(path.Join(outputFile, "neutral.gmm"), "eeg-neutral", neutral.Data, k)
}

func main() {

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	if len(os.Args) < 3 {
		panic("go run main.go <type> <bucket-size> <output-file> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	classifierType := os.Args[1]

	bucketSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[1]))
	}

	outputFile := os.Args[3]
	arguments, _, err := emotions.ParseArgumentsFromFile(os.Args[4], false)

	if err != nil {
		panic(err)
	}

	frameLen := 200
	frameStep := 200
	k := 3

	eegPositive, ok := arguments["eeg-positive"]
	if !ok {
		panic("No eeg positive files were provided")
	}
	eegNegative, ok := arguments["eeg-negative"]
	if !ok {
		panic("No eeg positive files were provided")
	}
	eegNeutral, ok := arguments["eeg-neutral"]
	if !ok {
		panic("No eeg positive files were provided")
	}

	switch classifierType {
	case "knn":
		err = KNN(eegPositive, eegNegative, eegNeutral, outputFile, bucketSize, frameLen, frameStep)
		if err != nil {
			panic(err)
		}
	case "gmm":
		err = GMM(eegPositive, eegNegative, eegNeutral, outputFile, k, bucketSize, frameLen, frameStep)
	default:
		panic(fmt.Sprintf("Unknown classifier %s", classifierType))
	}

}
