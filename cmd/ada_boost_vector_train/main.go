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

type EmotionFeature struct {
	emotion string
	feature []float64
}

func getEegFeaturesForFiles(emotion string, files []string, eegFeatures *[]EmotionFeature) [][]float64 {
	frameLen := 200
	frameStep := 150

	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		data = append(data, current...)
		for _, c := range current {
			*eegFeatures = append(*eegFeatures, EmotionFeature{emotion: emotion, feature: c})
		}
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

func getIncorrect(features []EmotionFeature, gmms []emotions.EmotionGausianMixure) []int {
	failedAll := 0
	incorrect := make([]int, len(features), len(features))

	for i := 0; i < len(features); i++ {
		best, failed := emotions.FindBestGaussian(features[i].feature, len(gmms[0].GM), gmms)
		if failed {
			failedAll++
			incorrect[i] = 1
			continue
		}
		if best != features[i].emotion {
			incorrect[i] = 1
		}
	}

	if failedAll != 0 {
		fmt.Fprintf(os.Stderr, "Failed: %d\n", failedAll)
	}
	return incorrect
}

func getFinalAlpha(emotionTypes []string,
	speechGMMs []emotions.EmotionGausianMixure,
	speechFeatures []EmotionFeature,
	eegGMMs []emotions.EmotionGausianMixure,
	eegFeatures []EmotionFeature,
) (float64, float64) {

	if len(eegFeatures) != len(speechFeatures) {
		panic(fmt.Sprintf("Features numbers must be equal to the number of files. EegFeatures: %d\tSpeechFeatures: %d\n", len(eegFeatures), len(speechFeatures)))
	}

	trainingSize := len(speechFeatures)

	weights := make([]float64, trainingSize, trainingSize)
	for i := 0; i < trainingSize; i++ {
		weights[i] = 1.0 / float64(trainingSize)
	}
	var alphaSpeech, alphaEeg float64
	var newWeights []float64

	speechIncorrect := getIncorrect(speechFeatures, speechGMMs)
	eegIncorrect := getIncorrect(eegFeatures, eegGMMs)

	errSpeech := getError(weights, speechIncorrect)
	errEeg := getError(weights, eegIncorrect)

	fmt.Printf("Err speech: %f, err eeg: %f\n", errSpeech, errEeg)

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

func tag(emotion string, data [][]float64, speechFeatures *[]EmotionFeature) {
	for _, m := range data {
		*speechFeatures = append(*speechFeatures, EmotionFeature{emotion: emotion, feature: m})
	}
}

func main() {
	if len(os.Args) < 4 {

		panic("go run main.go <speech-gmm-dir> <eeg-gmm-dir> <output-dir> <input_file>\n<input_file>: <emotion>	<wav_file>")
	}

	speechGmmDir := os.Args[1]
	eegGmmDir := os.Args[2]

	outputDir := os.Args[3]
	speechFiles, eegFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[4], true)

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

	speechGMMs, err := emotions.GetEGMs(speechGmmDir)
	if err != nil {
		panic(err)
	}

	eegGMMs, err := emotions.GetEGMs(eegGmmDir)
	if err != nil {
		panic(err)
	}

	speechFeatures := make([]EmotionFeature, 0, 1024)
	eegFeatures := make([]EmotionFeature, 0, 1024)

	sort.Strings(emotionTypes)

	for _, emotion := range emotionTypes {
		fmt.Printf("Reading vectors for %s\n", emotion)
		currentSpeechFiles := speechFiles[emotion]
		currentEegFiles := eegFiles[emotion]

		currentEegFeatures := getEegFeaturesForFiles(emotion, currentEegFiles, &eegFeatures)

		allFeatures := emotions.ReadSpeechFeatures(currentSpeechFiles)
		averaged := emotions.AverageSlice(allFeatures, len(allFeatures)/len(currentEegFeatures))
		currentSpeechFeatures := averaged[0 : len(averaged)-(len(averaged)-len(currentEegFeatures))]
		tag(emotion, currentSpeechFeatures, &speechFeatures)
	}
	speechAlpha, eegAlpha := getFinalAlpha(emotionTypes, speechGMMs, speechFeatures, eegGMMs, eegFeatures)
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
