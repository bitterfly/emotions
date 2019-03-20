package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(eegTrainSet []emotions.Tagged, eegTrainVars []float64, bucketSize int, frameLen int, frameStep int, emotion string, wav_filenames []string, eeg_filenames []string) [][]float64 {
	sort.Strings(wav_filenames)
	sort.Strings(eeg_filenames)

	// mfccs := make([][]float64, 0, len(filenames)*100)
	for i := range wav_filenames {
		// wf, _ := emotions.Read(f, 0.01, 0.97)
		fmt.Printf("wf: %s\nef: %s\n\n", wav_filenames[i], eeg_filenames[i])

		e := emotions.KNNOne(eegTrainSet, eegTrainVars, bucketSize, frameLen, frameStep, eeg_filenames[i])
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
		panic("go run main.go <k> <bucket-size> <dir-template> <eeg-model> <input-file>")
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

	outputDir := os.Args[4]

	wav_files, eeg_files, err := emotions.ParseArgumentsFromFile(os.Args[5], true)
	if err != nil {
		panic(err)
	}

	fmt.Printf("k: %d\nbucket: %d\noutput: %s\n", k, bucketSize, outputDir)

	eegTrainSet, err := emotions.UnmarshallEeg(eegTrainSetFilename)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshall eeg train set file: %s", eegTrainSetFilename))
	}
	_, eegTrainVars := emotions.GetμAndσTagged(eegTrainSet)

	for emotion, wf := range wav_files {
		readEmotion(eegTrainSet, eegTrainVars, bucketSize, 200, 150, emotion, wf, eeg_files[emotion])
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
