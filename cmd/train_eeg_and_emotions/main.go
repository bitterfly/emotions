package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(eegTrainSet []emotions.Tagged, eegTrainVars []float64, bucketSize int, frameLen int, frameStep int, emotion string, filenames []string) [][]float64 {
	// mfccs := make([][]float64, 0, len(filenames)*100)
	for _, f := range filenames {
		// wf, _ := emotions.Read(f, 0.01, 0.97)
		e := emotions.KNNOne(eegTrainSet, eegTrainVars, bucketSize, frameLen, frameStep, f)
		fmt.Printf("%s %f\n", emotion, e)
		// mfcc := emotions.MFCCs(wf, 13, 23)

		// mfccs = append(mfccs, mfcc...)
	}

	// return mfccs
	return nil
}

func getGMM(mfccs [][]float64, k int) emotions.GaussianMixture {
	return emotions.GMM(mfccs, k)
}

// func getGMMfromEmotion(filenames []string, k int) emotions.GaussianMixture {
// 	return getGMM(readEmotion(filenames), k)
// }

func main() {
	if len(os.Args) < 5 {
		panic("go run main.go <k> <bucket-size> <dir-template> <eeg-model> <--emotion1 emotion1.wav [emotion1.wav... --emotion2]>")
	}

	k, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic("Must provide k")
	}
	bucketSize, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic("Must provide bucketSize")
	}
	eegTrainSetFilename := os.Args[3]

	outputDir := os.Args[5]

	emotionFiles := emotions.ParseArguments(os.Args[6:])

	fmt.Printf("k: %d\nbucket: %d\noutput: %s\n", k, bucketSize, outputDir)

	eegTrainSet, err := emotions.UnmarshallEeg(eegTrainSetFilename)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshall eeg train set file: %s", eegTrainSetFilename))
	}
	_, eegTrainVars := emotions.GetμAndσTagged(eegTrainSet)

	for emotion, files := range emotionFiles {
		readEmotion(eegTrainSet, eegTrainVars, bucketSize, 200, 150, emotion, files)
		// 	for j := k; j <= maxK; j++ {
		// 		fmt.Fprintf(os.Stderr, "%s %d\n", emotion, j)
		// 		egm := emotions.EmotionGausianMixure{
		// 			Emotion: emotion,
		// 			GM:      getGMM(mfccs, j),
		// 		}

		// 		bytes, err := json.Marshal(egm)
		// 		if err != nil {
		// 			panic(fmt.Sprintf("Error when marshaling %s with k %d\n", emotion, j))
		// 		}

		// 		filename := path.Join(fmt.Sprintf("%s_k%d", outputDir, j), fmt.Sprintf("%s.gmm", emotion))
		// 		fmt.Fprintf(os.Stderr, filename)
		// 		ioutil.WriteFile(filename, bytes, 0644)
		// 	}
	}
}
