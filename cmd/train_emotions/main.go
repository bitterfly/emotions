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

func readEmotion(filenames []string) [][]float64 {
	mfccs := make([][]float64, 0, len(filenames)*100)
	for _, f := range filenames {
		wf, _ := emotions.Read(f, 0.01, 0.97)

		mfcc := emotions.MFCCs(wf, 13, 23)

		mfccs = append(mfccs, mfcc...)
	}

	return mfccs
}

func getGMM(mfccs [][]float64, k int) emotions.GaussianMixture {
	return emotions.GMM(mfccs, k)
}

func getGMMfromEmotion(filenames []string, k int) emotions.GaussianMixture {
	return getGMM(readEmotion(filenames), k)
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <k> <dir-template> <--emotion1 emotion1.wav [emotion1.wav... --emotion2]>")
	}

	outputDirIndex := 2
	maxK := -1

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}
	if _, err := strconv.Atoi(os.Args[2]); err != nil {
		outputDirIndex = 2
		maxK = k
	} else {
		maxK, _ = strconv.Atoi(os.Args[2])
		outputDirIndex = 3
	}

	outputDir := os.Args[outputDirIndex]
	emotionFiles := emotions.ParseArguments(os.Args[outputDirIndex+1:])

	for j := k; j <= maxK; j++ {

		if _, err := os.Stat(fmt.Sprintf("%s_k%d", outputDir, j)); os.IsNotExist(err) {
			os.Mkdir(fmt.Sprintf("%s_k%d", outputDir, j), 0775)
		}
	}

	for emotion, files := range emotionFiles {
		mfccs := readEmotion(files)
		for j := k; j <= maxK; j++ {
			fmt.Fprintf(os.Stderr, "%s %d\n", emotion, j)
			egm := emotions.EmotionGausianMixure{
				Emotion: emotion,
				GM:      getGMM(mfccs, j),
			}

			bytes, err := json.Marshal(egm)
			if err != nil {
				panic(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, j))
			}

			filename := path.Join(fmt.Sprintf("%s_k%d", outputDir, j), fmt.Sprintf("%s.gmm", emotion))
			fmt.Fprintf(os.Stderr, filename)
			ioutil.WriteFile(filename, bytes, 0644)
		}
	}
}
