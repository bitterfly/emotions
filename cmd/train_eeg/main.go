package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
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

func getData(bucketSize int, frameLen int, frameStep int, files []string, tag string) emotions.Tagged {
	data := make([][]float64, 0, 100)
	for i := range files {
		data = append(data, emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)...)
	}

	average := emotions.GetAverage(bucketSize, frameLen, len(data))

	return emotions.Tagged{

		Tag:  tag,
		Data: emotions.AverageSlice(data, average),
	}
}

func main() {

	//if bucket size is:
	// 0 take feature vectors every 200ms
	// 1 take only one vector for the whole file
	// n, n ≥ 2 take feature vector every n ms

	if len(os.Args) < 5 {
		panic("go run main.go bucket-size output-file --eeg-positive eeg-pos1.csv [eeg-pos2.csv...]  --eeg-negative eeg-neg1.csv [eeg-neg2.csv..] --eeg_neutral eeg-neu1.csv [eeg-neu2.csv...] ")
	}

	bucketSize, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[1]))
	}

	outputFile := os.Args[2]
	arguments := emotions.ParseArguments(os.Args[3:])

	if fileExists(outputFile) {
		os.Remove(outputFile)
	}

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

	fmt.Printf("OutputFile: %s\n", outputFile)
	fmt.Printf("eeg_pos: %s\n", eegPositive)
	fmt.Printf("eeg_neg: %s\n", eegNegative)
	fmt.Printf("eeg_neu: %s\n", eegNeutral)

	err = marshallToFile(bucketSize, 200, 150, outputFile, eegPositive, eegNegative, eegNeutral)
	if err != nil {
		panic(err.Error())
	}
}
