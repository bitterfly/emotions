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

func getEegFeaturesForFiles(emotion string, files []string, eegFeatures *[]emotions.EmotionFeature) [][]float64 {
	frameLen := 200
	frameStep := 150

	data := make([][]float64, 0, 100)
	for i := range files {
		current := emotions.GetFourierForFile(files[i], 19, frameLen, frameStep)
		data = append(data, current...)
		for _, c := range current {
			*eegFeatures = append(*eegFeatures, emotions.EmotionFeature{Emotion: emotion, Feature: c})
		}
	}

	return data
}

func getIncorrect(features []emotions.EmotionFeature, gmms []emotions.EmotionGausianMixure) []int {
	failedAll := 0
	incorrect := make([]int, len(features), len(features))

	for i := 0; i < len(features); i++ {
		best, failed := emotions.FindBestGaussian(features[i].Feature, len(gmms[0].GM), gmms)
		if failed {
			failedAll++
			incorrect[i] = 1
			continue
		}
		if best != features[i].Emotion {
			incorrect[i] = 1
		}
	}

	if failedAll != 0 {
		fmt.Fprintf(os.Stderr, "Failed: %d\n", failedAll)
	}
	return incorrect
}

func gradientDescent(emotionTypes []string,
	speechGMMs []emotions.EmotionGausianMixure,
	speechFeatures []emotions.EmotionFeature,
	eegGMMs []emotions.EmotionGausianMixure,
	eegFeatures []emotions.EmotionFeature,
) (float64, float64) {

	if len(eegFeatures) != len(speechFeatures) {
		panic(fmt.Sprintf("Features numbers must be equal to the number of files. EegFeatures: %d\tSpeechFeatures: %d\n", len(eegFeatures), len(speechFeatures)))
	}

	// trainingSize := len(speechFeatures)
	fmt.Printf("Summing corpus for speech\n")
	speechSumOverCorpus, failed := emotions.SumGmmOverCorpus(speechFeatures, speechGMMs)
	if failed {
		fmt.Errorf("Some vectors have failed")
	}

	fmt.Printf("Summing corpus for eeg\n")
	eegSumOverCorpus, failed := emotions.SumGmmOverCorpus(eegFeatures, eegGMMs)
	if failed {
		fmt.Errorf("Some vectors have failed")
	}

	fmt.Printf("eeg: %f speech %f\n", eegSumOverCorpus, speechSumOverCorpus)
	speechAlpha := 0.5
	eegAlpha := 0.5

	L := speechAlpha*speechSumOverCorpus + eegAlpha*eegSumOverCorpus
	for i := 0; i < len(speechFeatures); i++ {
		sum := 0.0
		for j := 0; j < len(speechGMMs); j++ {
			curSpeech, _ := emotions.GmmsProbabilityGivenClass(speechFeatures[i].Feature, speechGMMs[j].Emotion, speechGMMs)
			curEeg, _ := emotions.GmmsProbabilityGivenClass(eegFeatures[i].Feature, speechGMMs[j].Emotion, eegGMMs)
			sum += math.Exp(speechAlpha*curSpeech + eegAlpha*curEeg)
		}
		L -= math.Log(sum)
	}
	fmt.Printf("k: 0, L: %f\n", L)
	var snom, enom, den float64
	prevL := L

	for k := 0; k < 200; k++ {
		newSpeechAlpha := speechSumOverCorpus
		newEegAlpha := eegSumOverCorpus
		snom = 0.0
		enom = 0.0
		den = 0.0
		for i := 0; i < len(speechFeatures); i++ {
			for j := 0; j < len(speechGMMs); j++ {
				curSpeech, _ := emotions.GmmsProbabilityGivenClass(speechFeatures[i].Feature, speechGMMs[j].Emotion, speechGMMs)
				curEeg, _ := emotions.GmmsProbabilityGivenClass(eegFeatures[i].Feature, speechGMMs[j].Emotion, eegGMMs)
				exp := math.Exp(speechAlpha*curSpeech + eegAlpha*curEeg)
				snom += exp * curSpeech
				enom += exp * curEeg
				den += exp
			}
		}
		newSpeechAlpha -= snom / den
		newEegAlpha -= enom / den

		// L
		L := speechAlpha*speechSumOverCorpus + eegAlpha*eegSumOverCorpus
		for i := 0; i < len(speechFeatures); i++ {
			sum := 0.0
			for j := 0; j < len(speechGMMs); j++ {
				curSpeech, _ := emotions.GmmsProbabilityGivenClass(speechFeatures[i].Feature, speechGMMs[j].Emotion, speechGMMs)
				curEeg, _ := emotions.GmmsProbabilityGivenClass(eegFeatures[i].Feature, speechGMMs[j].Emotion, eegGMMs)
				sum += math.Exp(speechAlpha*curSpeech + eegAlpha*curEeg)
			}
			L -= math.Log(sum)
		}
		fmt.Printf("k: %d, L: %f\n", k, L)
		if L < prevL {
			fmt.Printf("Breaking on step %d\n", k)
			break
		}
		prevL = L
		speechAlpha = speechAlpha + 0.0001*newSpeechAlpha
		eegAlpha = eegAlpha + 0.0001*newEegAlpha
		fmt.Printf("Speech: %f, eeg: %f\n", speechAlpha, eegAlpha)
	}

	return speechAlpha, eegAlpha
}

func tag(emotion string, data [][]float64, speechFeatures *[]emotions.EmotionFeature) {
	for _, m := range data {
		*speechFeatures = append(*speechFeatures, emotions.EmotionFeature{Emotion: emotion, Feature: m})
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

	speechFeatures := make([]emotions.EmotionFeature, 0, 1024)
	eegFeatures := make([]emotions.EmotionFeature, 0, 1024)

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

	fmt.Printf("Getting weights\n")
	speechAlpha, eegAlpha := gradientDescent(emotionTypes, speechGMMs, speechFeatures, eegGMMs, eegFeatures)
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
