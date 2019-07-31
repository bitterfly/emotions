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

func readSpeechFeaturesForFiles(filenames []string, speechFeatures *[]([][]float64)) [][]float64 {
	mfccs := make([][]float64, 0, len(filenames)*100)
	for _, f := range filenames {
		wf, _ := emotions.Read(f, 0.01, 0.97)

		mfcc := emotions.MFCCs(wf, 13, 23)
		mfccs = append(mfccs, mfcc...)

		*speechFeatures = append(*speechFeatures, mfcc)
	}
	return mfccs
}

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

func getGMM(mfccs [][]float64, k int) emotions.GaussianMixture {
	return emotions.GMM(mfccs, k)
}
func getError(emotionTypes []string, trainingSize int,
	files map[string][]string,
	gmms []emotions.EmotionGausianMixure,
	features []([][]float64),
	tags []string,
	weights []float64,
) (float64, []int) {
	incorrect := make([]int, trainingSize, trainingSize)

	err := 0.0
	sum := 0.0
	for i := 0; i < len(features); i++ {
		sum += weights[i]
		incorrect[i] = 0
		correct, _, _ := emotions.TestGMM(tags[i], emotionTypes, features[i], gmms)
		if correct == 0 {
			err += weights[i]
			incorrect[i] = 1
		}
	}
	return err / sum, incorrect
}

func getWeightsAndAlpha(emotionTypes []string, trainingSize int,
	files map[string][]string,
	gmms []emotions.EmotionGausianMixure,
	features []([][]float64),
	tags []string,
	weights []float64,
) (float64, []float64) {
	err, incorrectFiles := getError(emotionTypes, trainingSize, files, gmms, features, tags, weights)
	k := len(emotionTypes)

	alpha := math.Log((1-err)/err) + math.Log(float64(k-1))
	newWeights := make([]float64, len(incorrectFiles), len(incorrectFiles))
	sum := 0.0
	for i := 0; i < len(incorrectFiles); i++ {
		newWeights[i] = weights[i] * math.Exp(alpha*float64(incorrectFiles[i]))
		sum += newWeights[i]
	}
	fmt.Printf("\n")

	for i := 0; i < len(incorrectFiles); i++ {
		newWeights[i] /= sum
	}

	return alpha, newWeights
}

func getWeights(emotionTypes []string, trainingSize int, vectorTags []string,
	speechFiles map[string][]string,
	speechGMMs []emotions.EmotionGausianMixure,
	speechFeatures []([][]float64),
	eegFiles map[string][]string,
	eegGMMs []emotions.EmotionGausianMixure,
	eegFeatures []([][]float64),
) (float64, float64) {
	weights := make([]float64, trainingSize, trainingSize)
	for i := 0; i < trainingSize; i++ {
		weights[i] = 1.0 / float64(trainingSize)
	}

	var alphaSpeech, alphaEeg float64
	var newWeights []float64

	errSpeech, _ := getError(emotionTypes, trainingSize, speechFiles, speechGMMs, speechFeatures, vectorTags, weights)
	errEeg, _ := getError(emotionTypes, trainingSize, eegFiles, eegGMMs, eegFeatures, vectorTags, weights)

	for i := 0; i < 10; i++ {
		if errSpeech <= errEeg {
			alphaSpeech, newWeights = getWeightsAndAlpha(emotionTypes, trainingSize, speechFiles, speechGMMs, speechFeatures, vectorTags, weights)
		} else {
			alphaEeg, newWeights = getWeightsAndAlpha(emotionTypes, trainingSize, eegFiles, eegGMMs, eegFeatures, vectorTags, weights)
		}
		errEeg, _ = getError(emotionTypes, trainingSize, eegFiles, eegGMMs, eegFeatures, vectorTags, newWeights)
		errSpeech, _ = getError(emotionTypes, trainingSize, speechFiles, speechGMMs, speechFeatures, vectorTags, newWeights)

		weights = append([]float64{}, newWeights...)
	}

	return alphaSpeech, alphaEeg
}

func main() {
	if len(os.Args) < 3 {

		panic("go run main.go <k> <dir-template> <input_file>\n<input_file>: <emotion>	<wav_file>")
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

	vectorTags := make([]string, 0, 1024)
	for i, emotion := range emotionTypes {
		currentSpeechFiles := speechFiles[emotion]
		currentEegFiles := eegFiles[emotion]
		wLen += len(currentSpeechFiles)
		for j := 0; j < len(currentSpeechFiles); j++ {
			vectorTags = append(vectorTags, emotion)
		}

		currentSpeechFeatures := readSpeechFeaturesForFiles(currentSpeechFiles, &speechFeatures)
		currentEegFeatures := getEegFeaturesForFiles(bucketSize, currentEegFiles, &eegFeatures)

		fmt.Printf("len speech: %d\nlen eeg %d\n", len(speechFeatures), len(eegFeatures))

		speechGMMs[i] = emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      getGMM(currentSpeechFeatures, k),
		}

		eegGMMs[i] = emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      getGMM(currentEegFeatures, 3),
		}
	}

	speechAlpha, eegAlpha := getWeights(emotionTypes, wLen, vectorTags, speechFiles, speechGMMs, speechFeatures, eegFiles, eegGMMs, eegFeatures)

	for i := 0; i < len(speechGMMs); i++ {
		emotion := emotionTypes[i]
		bytes, err := json.Marshal(emotions.AlphaEGM{Alpha: speechAlpha, EGM: eegGMMs[i]})
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s", emotion))
		}

		filename := path.Join(fmt.Sprintf("%s_speech", outputDir), fmt.Sprintf("%s.gmm", emotion))
		fmt.Fprintf(os.Stderr, filename)
		ioutil.WriteFile(filename, bytes, 0644)
	}

	for i := 0; i < len(eegGMMs); i++ {
		emotion := emotionTypes[i]
		bytes, err := json.Marshal(emotions.AlphaEGM{Alpha: eegAlpha, EGM: eegGMMs[i]})
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s", emotion))
		}

		filename := path.Join(fmt.Sprintf("%s_eeg", outputDir), fmt.Sprintf("%s.gmm", emotion))
		fmt.Fprintf(os.Stderr, filename)
		ioutil.WriteFile(filename, bytes, 0644)
	}

}
