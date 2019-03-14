package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func MarshallToFile(filename string, positiveFiles []string, negativeFiles []string, neutralFiles []string) error {
	bytes, err := json.Marshal([]emotions.Tagged{
		getData(positiveFiles, "positive"),
		getData(negativeFiles, "negative"),
		getData(neutralFiles, "neutral"),
	})
	if err != nil {
		return fmt.Errorf("Could not marshal eeg data: %s\n", err.Error())
	}

	err = ioutil.WriteFile(filename, bytes, 0664)
	return err
}

func getData(files []string, tag string) emotions.Tagged {
	data := make([][]float64, 0, 100)
	for i := range files {
		data = append(data, emotions.GetFourierForFile(files[i], 19)...)
	}

	return emotions.Tagged{
		Tag:  tag,
		Data: data,
	}
}

func main() {
	if len(os.Args) < 6 {
		panic("go run main.go output_file --eeg-positive eeg_pos1.csv [eeg_pos2.csv...]  --eeg-negative eeg_neg1.csv [eeg_neg2.csv..] --eeg_neutral eeg_neu1.csv [eeg_neu2.csv...] ")
	}

	output_file := os.Args[1]
	arguments := emotions.ParseArguments(os.Args[2:])

	if fileExists(output_file) {
		os.Remove(output_file)
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

	fmt.Printf("Output_file: %s\n", output_file)
	fmt.Printf("eeg_pos: %s\n", eegPositive)
	fmt.Printf("eeg_neg: %s\n", eegNegative)
	fmt.Printf("eeg_neu: %s\n", eegNeutral)

	err := MarshallToFile(output_file, eegPositive, eegNegative, eegNeutral)
	if err != nil {
		panic(err.Error())
	}
}
