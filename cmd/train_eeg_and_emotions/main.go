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

func combine(bucketSize int, mfccs [][]float64, eegVectors [][]float64) [][]float64 {
	if bucketSize < emotions.STEP_IN_MS {
		panic(fmt.Sprintf("Eeg data should not be sampled in less than %dms", emotions.FRAME_IN_MS))
	}

	rat := bucketSize / emotions.STEP_IN_MS
	featureVectors := make([][]float64, len(eegVectors), len(eegVectors))

	ev := 0
	currentMfcc := make([]float64, len(mfccs[0]), len(mfccs[0]))
	num := 0
	for i := range mfccs {
		if i%rat == 0 && i != 0 {
			emotions.Divide(&currentMfcc, float64(num))
			featureVectors[ev] = append(featureVectors[ev], currentMfcc...)
			featureVectors[ev] = append(featureVectors[ev], eegVectors[ev]...)
			ev++
			num = 0
			emotions.Zero(&currentMfcc)
		}
		emotions.Add(&currentMfcc, mfccs[i])
		num++
	}
	for i := range featureVectors {
		fmt.Printf("i: %d, len: %d\n", i, len(featureVectors[i]))
	}

	fmt.Printf("Ev: %d\n", ev)
	if ev < len(eegVectors) {
		emotions.Divide(&currentMfcc, float64(num))
		featureVectors[ev] = append(featureVectors[ev], currentMfcc...)
		featureVectors[ev] = append(featureVectors[ev], eegVectors[ev]...)
	}

	for i := range featureVectors {
		fmt.Printf("i: %d, len: %d\n", i, len(featureVectors[i]))
		if len(featureVectors[i]) != len(mfccs[0])+len(eegVectors[0]) {
			panic(fmt.Sprintf("len: %d i: %d, vec: %v\n", len(featureVectors), i, featureVectors[i]))
		}
	}
	return featureVectors
}

func getEegVectors(bucketSize int, frameLen int, frameStep int, filename string) [][]float64 {
	current := emotions.GetFourierForFile(filename, 19, frameLen, frameStep)
	fmt.Printf("Current gas: %d\n", len(current))

	average := emotions.GetAverage(bucketSize, frameStep, len(current))
	fmt.Printf("Average is: %d\n", average)
	return emotions.AverageSlice(current, average)
}

func readEmotion(bucketSize int, frameLen int, frameStep int, wavFilenames []string, eegFilenames []string) [][]float64 {
	sort.Strings(wavFilenames)
	sort.Strings(eegFilenames)

	featureVectors := make([][]float64, 0, len(wavFilenames)*100)
	for i := range wavFilenames {
		wf, _ := emotions.Read(wavFilenames[i], 0.01, 0.97)
		fmt.Printf("wf: %s\nef: %s\n\n", wavFilenames[i], eegFilenames[i])

		mfcc := emotions.MFCCs(wf, 13, 23)
		fmt.Fprintf(os.Stderr, "Got %d mfccs for %s\n", len(mfcc), wavFilenames[i])
		eegVectors := getEegVectors(bucketSize, frameLen, frameStep, eegFilenames[i])
		fmt.Fprintf(os.Stderr, "Got %d eeg vectors for %s\n", len(eegVectors), eegFilenames[i])

		featureVectors = append(featureVectors, combine(bucketSize, mfcc, eegVectors)...)
	}

	f, _ := os.Create("/tmp/foo")
	defer f.Close()
	for _, v := range featureVectors {
		for i, vv := range v {
			fmt.Fprintf(f, "%d: %f ", i, vv)
		}
		fmt.Fprintf(f, "\n")
	}
	return featureVectors
}

func getGMM(mfccs [][]float64, k int) emotions.GaussianMixture {
	return emotions.GMM(mfccs, k)
}

func main() {
	if len(os.Args) < 4 {
		panic("go run main.go <k> <bucket-size> <dir-template> <input-file>")
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}
	bucketSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Must provide bucketSize")
	}
	outputDir := os.Args[3]

	wavFiles, eegFiles, err := emotions.ParseArgumentsFromFile(os.Args[4], true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("k: %d\nbucket: %d\noutput: %s\nwf: %d\neegf: %d\n", k, bucketSize, outputDir, len(wavFiles), len(eegFiles))

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		os.Mkdir(outputDir, 0775)
	}

	emotionsNames := emotions.SortKeysS(wavFiles)
	for _, emotion := range emotionsNames {
		featureVectors := readEmotion(bucketSize, 200, 150, wavFiles[emotion], eegFiles[emotion])
		egm := emotions.EmotionGausianMixure{
			Emotion: emotion,
			GM:      getGMM(featureVectors, k),
		}

		bytes, err := json.Marshal(egm)
		if err != nil {
			panic(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, k))
		}

		filename := path.Join(outputDir, fmt.Sprintf("%s.gmm", emotion))
		fmt.Fprintf(os.Stderr, filename)
		ioutil.WriteFile(filename, bytes, 0644)
	}
}
