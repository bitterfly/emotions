package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"sort"

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

func getError(weights []float64, incorrect []int) float64 {
	err := 0.0
	sum := 0.0
	for i := 0; i < len(incorrect); i++ {
		sum += weights[i]
		err += weights[i] * float64(incorrect[i])
	}

	if err <= emotions.EPS && err >= -emotions.EPS {
		err = emotions.EPS
	}

	return err / sum
}

func getIncorrect(emotionTypes []string, filesTags []string, features []([][]float64), gmms []emotions.EmotionGausianMixure) []int {
	incorrect := make([]int, len(filesTags), len(filesTags))

	for i := 0; i < len(features); i++ {
		correct, _, _ := emotions.TestGMM(filesTags[i], emotionTypes, features[i], gmms, false)
		if correct == 0 {
			incorrect[i] = 1
		}
	}

	return incorrect
}

func getAlpha(err float64, k int) float64 {
	return math.Log((1-err)/err) + math.Log(float64(k-1))
}

func getWeights(alpha float64, weights []float64, incorrect []int) []float64 {
	newWeights := make([]float64, len(incorrect), len(incorrect))
	sum := 0.0
	for i := 0; i < len(incorrect); i++ {
		newWeights[i] = weights[i] * math.Exp(alpha*float64(incorrect[i]))
		sum += newWeights[i]
	}

	for i := 0; i < len(incorrect); i++ {
		newWeights[i] /= sum
	}

	return newWeights
}

func getFinalAlpha(emotionTypes []string, filesTags []string,
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

	speechIncorrect := getIncorrect(emotionTypes, filesTags, speechFeatures, speechGMMs)
	eegIncorrect := getIncorrect(emotionTypes, filesTags, eegFeatures, eegGMMs)

	errSpeech := getError(weights, speechIncorrect)
	fmt.Printf("ErrSpeech: %f\n", errSpeech)
	errEeg := getError(weights, eegIncorrect)
	fmt.Printf("ErrEeg: %f\n", errEeg)

	k := len(emotionTypes)

	for i := 0; i < 10; i++ {
		if errSpeech <= errEeg {
			newalpha := getAlpha(errSpeech, k)
			newWeights = getWeights(newalpha, weights, speechIncorrect)
			alphaSpeech += newalpha
			fmt.Printf("Choosing speech with alpha %f and err %f\n", newalpha, errSpeech)

			if newalpha <= emotions.EPS && newalpha >= -emotions.EPS {
				break
			}
		} else {
			newalpha := getAlpha(errEeg, k)
			newWeights = getWeights(newalpha, weights, eegIncorrect)
			alphaEeg += newalpha
			fmt.Printf("Choosing eeg with alpha %f and err %f\n", newalpha, errEeg)

			if newalpha <= emotions.EPS && newalpha >= -emotions.EPS {
				break
			}
		}

		errSpeech = getError(newWeights, speechIncorrect)
		errEeg = getError(newWeights, eegIncorrect)

		weights = append([]float64{}, newWeights...)
	}
	fmt.Printf("Alpha speech: %f, alpha eeg: %f\n", alphaSpeech, alphaEeg)

	return alphaSpeech, alphaEeg
}

func main() {
	if len(os.Args) < 4 {

		panic("go run main.go <speech-gmm-dir> <eeg-gmm-dir> <output-dir> <input_file>\n<input_file>: <emotion>	<wav_file>")
	}

	speechGmmDir := os.Args[1]
	eegGmmDir := os.Args[2]

	outputDir := os.Args[3]
	speechFiles, eegFiles, filesTags, err := emotions.ParseArgumentsFromFile(os.Args[4], true)

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
	sort.Strings(emotionTypes)

	speechFeatures := make([]([][]float64), 0, 1024)
	eegFeatures := make([]([][]float64), 0, 1024)

	speechGMMs, err := emotions.GetEGMs(speechGmmDir)
	if err != nil {
		panic(err)
	}

	eegGMMs, err := emotions.GetEGMs(eegGmmDir)
	if err != nil {
		panic(err)
	}

	for _, emotion := range emotionTypes {
		fmt.Printf("Processing files for %s\n", emotion)
		currentSpeechFiles := speechFiles[emotion]
		currentEegFiles := eegFiles[emotion]
		emotions.ReadSpeechFeaturesAppend(currentSpeechFiles, &speechFeatures)
		getEegFeaturesForFiles(0, currentEegFiles, &eegFeatures)
	}

	speechAlpha, eegAlpha := getFinalAlpha(emotionTypes, filesTags, speechGMMs, speechFeatures, eegGMMs, eegFeatures)
	if speechAlpha < emotions.EPS && speechAlpha > -emotions.EPS {
		speechAlpha = emotions.EPS
	}

	if eegAlpha < emotions.EPS && eegAlpha > -emotions.EPS {
		eegAlpha = emotions.EPS
	}

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
