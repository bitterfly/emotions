package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func getEegFeaturesForFiles(bucketSize int, files []string, eegFeatures *[]([][]float64)) [][]float64 {
	frameLen := 200
	frameStep := 150

	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		average := emotions.GetAverage(bucketSize, frameLen, len(current))
		averaged := emotions.AverageSlice(current, average)
		data = append(data, averaged...)

		*eegFeatures = append(*eegFeatures, averaged)
	}

	return data
}

func getError(emotionTypes []string, filesTags []string,
	gmms []emotions.EmotionGausianMixure,
	features []([][]float64),
	weights []float64,
) (float64, []int) {
	incorrect := make([]int, len(filesTags), len(filesTags))

	err := 0.0
	sum := 0.0

	for i := 0; i < len(features); i++ {
		sum += weights[i]

		correct, _, _ := emotions.TestGMM(filesTags[i], emotionTypes, features[i], gmms, false)
		if correct == 0 {
			err += weights[i]
			incorrect[i] = 1
		}
	}

	return err / sum, incorrect
}

func getWeightsAndAlpha(emotionTypes []string, filesTags []string,
	gmms []emotions.EmotionGausianMixure,
	features []([][]float64),
	weights []float64,
) (float64, []float64) {
	err, incorrectFiles := getError(emotionTypes, filesTags, gmms, features, weights)
	if err < emotions.EPS && err > -emotions.EPS {
		err = emotions.EPS
	}

	k := len(emotionTypes)
	// fmt.Printf("Err: %f\tIncorrect: %v\n", err, incorrectFiles)
	alpha := math.Log((1-err)/err) + math.Log(float64(k-1))
	// fmt.Printf("Alpha: %f\n", alpha)
	newWeights := make([]float64, len(incorrectFiles), len(incorrectFiles))
	sum := 0.0
	for i := 0; i < len(incorrectFiles); i++ {
		newWeights[i] = weights[i] * math.Exp(alpha*float64(incorrectFiles[i]))
		sum += newWeights[i]
	}

	for i := 0; i < len(incorrectFiles); i++ {
		newWeights[i] /= sum
	}

	return alpha, newWeights
}

func getWeights(emotionTypes []string, filesTags []string,
	speechGMMs []emotions.EmotionGausianMixure,
	speechFeatures []([][]float64),
	eegGMMs []emotions.EmotionGausianMixure,
	eegFeatures []([][]float64),
) (float64, float64) {

	if len(eegFeatures) != len(speechFeatures) || len(eegFeatures) != len(filesTags) {
		panic(fmt.Sprintf("Features numbers must be equal to the number of files. EegFeatures: %d\tSpeechFeatures: %d\tFiles: %d\n", len(eegFeatures), len(speechFeatures), len(filesTags)))
	}

	trainingSize := len(filesTags)

	weights := make([]float64, trainingSize, trainingSize)
	for i := 0; i < trainingSize; i++ {
		weights[i] = 1.0 / float64(trainingSize)
	}

	var alphaSpeech, alphaEeg float64
	var newWeights []float64

	errSpeech, _ := getError(emotionTypes, filesTags, speechGMMs, speechFeatures, weights)
	errEeg, _ := getError(emotionTypes, filesTags, eegGMMs, eegFeatures, weights)

	for i := 0; i < 5; i++ {
		if errSpeech <= errEeg {
			fmt.Printf("Choosing speech\n")
			alphaSpeech, newWeights = getWeightsAndAlpha(emotionTypes, filesTags, speechGMMs, speechFeatures, weights)
		} else {
			fmt.Printf("Choosing eeg\n")
			alphaEeg, newWeights = getWeightsAndAlpha(emotionTypes, filesTags, eegGMMs, eegFeatures, weights)
		}
		errEeg, _ = getError(emotionTypes, filesTags, eegGMMs, eegFeatures, newWeights)
		errSpeech, _ = getError(emotionTypes, filesTags, speechGMMs, speechFeatures, newWeights)

		weights = append([]float64{}, newWeights...)
		// fmt.Printf("AlphaSpeech: %f\nAlphaEEG: %f\n NewWeights: %v\n", alphaSpeech, alphaEeg, weights)
	}

	return alphaSpeech, alphaEeg
}

func main() {
	if len(os.Args) < 4 {

		panic("go run main.go <k> <bucket-size> <dir-template> <input_file>\n<input_file>: <emotion>	<wav_file>")
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}

	bucketSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(fmt.Sprintf("could not parse bucket-size argument: %s", os.Args[2]))
	}

	outputDir := os.Args[3]
	speechFiles, eegFiles, err := emotions.ParseArgumentsFromFile(os.Args[4], true)

	if err != nil || len(speechFiles) != len(eegFiles) {
		panic(err)
	}

	if _, err := os.Stat(fmt.Sprintf("%s_speech", outputDir)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s_speech", outputDir), 0775)
	}

	if _, err := os.Stat(fmt.Sprintf("%s_eeg", outputDir)); os.IsNotExist(err) {
		os.Mkdir(fmt.Sprintf("%s_eeg", outputDir), 0775)
	}

	emotionTypes := make([]string, 0, len(speechFiles))
	for e := range speechFiles {
		emotionTypes = append(emotionTypes, e)
	}

	speechGMMs := make([]emotions.EmotionGausianMixure, len(emotionTypes), len(emotionTypes))
	eegGMMs := make([]emotions.EmotionGausianMixure, len(emotionTypes), len(emotionTypes))

	speechFeatures := make([]([][]float64), 0, 1024)
	eegFeatures := make([]([][]float64), 0, 1024)

	sort.Strings(emotionTypes)

	wLen := 0

	eegFilesSorted := make([]string, 0, 1024)
	filesTags := make([]string, 0, 1024)

	for i, emotion := range emotionTypes {
		currentSpeechFiles := speechFiles[emotion]
		currentEegFiles := eegFiles[emotion]
		wLen += len(currentSpeechFiles)
		for j := 0; j < len(currentSpeechFiles); j++ {
			filesTags = append(filesTags, emotion)
		}
		eegFilesSorted = append(eegFilesSorted, currentEegFiles...)
		currentSpeechFeatures := emotions.ReadSpeechFeaturesAppend(currentSpeechFiles, &speechFeatures)

		speechGMMs[i] = emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      emotions.GMM(currentSpeechFeatures, k),
		}

		currentEegFeatures := getEegFeaturesForFiles(bucketSize, currentEegFiles, &eegFeatures)

		eegGMMs[i] = emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      emotions.GMM(currentEegFeatures, 3),
		}
	}
	speechAlpha, eegAlpha := getWeights(emotionTypes, filesTags, speechGMMs, speechFeatures, eegGMMs, eegFeatures)
	if speechAlpha < emotions.EPS && speechAlpha > -emotions.EPS {
		speechAlpha = emotions.EPS
	}

	if eegAlpha < emotions.EPS && eegAlpha > -emotions.EPS {
		eegAlpha = emotions.EPS
	}
	// getWeights(emotionTypes, wLen, filesTags, speechFiles, speechGMMs, speechFeatures, eegFiles, eegGMMs, eegFeatures)

	for i := 0; i < len(speechGMMs); i++ {
		emotion := speechGMMs[i].Emotion
		bytes, err := json.Marshal(emotions.AlphaEGM{Alpha: speechAlpha, EGM: speechGMMs[i]})
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s - %s", emotion, err.Error()))
		}

		filename := path.Join(fmt.Sprintf("%s_speech", outputDir), fmt.Sprintf("%s.gmm", emotion))
		fmt.Fprintf(os.Stderr, filename)
		ioutil.WriteFile(filename, bytes, 0644)
	}

	for i := 0; i < len(eegGMMs); i++ {
		emotion := eegGMMs[i].Emotion

		bytes, err := json.Marshal(emotions.AlphaEGM{Alpha: eegAlpha, EGM: eegGMMs[i]})
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s - %s", emotion, err.Error()))
		}

		filename := path.Join(fmt.Sprintf("%s_eeg", outputDir), fmt.Sprintf("%s.gmm", emotion))
		fmt.Fprintf(os.Stderr, filename)
		ioutil.WriteFile(filename, bytes, 0644)
	}
}
