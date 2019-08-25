package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func getEegFeaturesForFiles(bucketSize int, files []string) [][]float64 {
	frameLen := 200
	frameStep := 150

	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		average := emotions.GetAverage(bucketSize, frameLen, len(current))
		averaged := emotions.AverageSlice(current, average)
		data = append(data, averaged...)
	}

	return data
}

func main() {
	if len(os.Args) < 4 {

		panic("go run main.go <k> <dir-template> <input_file>\n<input_file>: <emotion>	<wav_file>")
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}
	k = 3

	outputDir := os.Args[2]
	speechFiles, eegFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[3], true)

	if err != nil || len(speechFiles) != len(eegFiles) {
		panic(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s_k%d", outputDir, k)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s_k%d", outputDir, k), 0775)
	}

	emotionTypes := make([]string, 0, len(speechFiles))
	for e := range speechFiles {
		emotionTypes = append(emotionTypes, e)
	}

	sort.Strings(emotionTypes)
	var gmm emotions.EmotionGausianMixure

	for _, emotion := range emotionTypes {
		currentSpeechFiles := speechFiles[emotion]
		currentEegFiles := eegFiles[emotion]

		currentEegFeatures := getEegFeaturesForFiles(0, currentEegFiles)
		allFeatures := emotions.ReadSpeechFeatures(currentSpeechFiles)
		averaged := emotions.AverageSlice(allFeatures, len(allFeatures)/len(currentEegFeatures))

		currentSpeechFeatures := averaged[0 : len(averaged)-(len(averaged)-len(currentEegFeatures))]

		fmt.Printf("Emotion: %s eeg len: %d  all %d averaged: %d  speech len: %d\n", emotion, len(currentEegFeatures), len(allFeatures), len(averaged), len(currentSpeechFeatures))

		gmm = emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      emotions.GMM(emotions.Concat(currentSpeechFeatures, currentEegFeatures), k),
		}
		bytes, err := json.Marshal(gmm)
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s - %s", emotion, err.Error()))
		}
		filename := path.Join(fmt.Sprintf("%s_k%d", outputDir, k), fmt.Sprintf("%s.gmm", emotion))
		ioutil.WriteFile(filename, bytes, 0644)
	}
}
