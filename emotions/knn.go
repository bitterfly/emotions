package emotions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
)

type Tagged struct {
	Tag  string
	Data [][]float64
}

func UnmarshallKNNEeg(filename string) ([]Tagged, error) {
	var tagged []Tagged
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(bytes, &tagged)
	if err != nil {
		return nil, err
	}

	return tagged, nil
}

func UnmarshallGMMEeg(filename string) ([]EmotionGausianMixure, error) {
	return nil, nil
}

func GetμAndσTagged(tagged []Tagged) ([]float64, []float64) {
	featureVectorSize := len(tagged[0].Data[0])
	// fmt.Fprintf(os.Stderr, "FeatureVectorSize: %d\n", featureVectorSize)

	variances := make([]float64, featureVectorSize, featureVectorSize)

	expectation := make([]float64, featureVectorSize, featureVectorSize)
	expectationSquared := make([]float64, featureVectorSize, featureVectorSize)

	for t := range tagged {
		for i := range tagged[t].Data {
			for j := range tagged[t].Data[i] {
				expectation[j] += tagged[t].Data[i][j]
				expectationSquared[j] += expectation[j] * expectation[j]
			}
		}
	}

	numVectors := 0
	for t := range tagged {
		numVectors += len(tagged[t].Data)
	}

	// fmt.Fprintf(os.Stderr, "numVectors: %d\n", numVectors)
	for j := 0; j < featureVectorSize; j++ {
		expectation[j] /= float64(numVectors)
		expectationSquared[j] /= float64(numVectors)

		variances[j] = expectationSquared[j] - expectation[j]*expectation[j]
	}

	// fmt.Fprintf(os.Stderr, "Expectations: %d\n", len(expectation))
	// fmt.Fprintf(os.Stderr, "Variances: %d\n", len(variances))
	return expectation, variances
}

func findClosest(v []float64, trainSet []Tagged, trainVar []float64) string {
	minDist := math.Inf(42)
	minDistTag := ""

	for t := range trainSet {
		for tv := range trainSet[t].Data {
			curDist := mahalanobisDistance(v, trainSet[t].Data[tv], trainVar)
			if curDist < minDist {
				minDist = curDist
				minDistTag = trainSet[t].Tag
			}
		}
	}

	return minDistTag
}

func testKNN(emotion string, emotions []string, vectors [][]float64, trainSet []Tagged, trainVar []float64) (int, int, int) {
	fmt.Printf("%s\t", emotion)

	counters := make(map[string]int)
	for _, e := range emotions {
		counters[e] = 0
	}

	for v := range vectors {
		counters[findClosest(vectors[v], trainSet, trainVar)]++
	}

	sum := 0
	for _, e := range emotions {
		fmt.Printf("%d\t", counters[e])
		sum += counters[e]
	}
	fmt.Printf("\n")

	return correct(emotion, counters), counters[emotion], sum
}

// func ClassifyOne(trainSet []Tagged, trainVar []float64, bucketSize int, frameLen int, frameStep int, filename string) float64 {
// 	vec := GetFourierForFile(filename, 19, frameLen, frameStep)
// 	average := GetAverage(bucketSize, frameStep, len(vec))
// 	averaged := AverageSlice(vec, average)
// 	_, mc := testKNN(averaged, trainSet, trainVar)
// 	switch mc {
// 	case "eeg-neutral":
// 		return 0
// 	case "eeg-positive":
// 		return 1
// 	case "eeg-negative":
// 		return -1
// 	}
// 	return 0
// }

func ClassifyKNN(trainSetFilename string, bucketSize int, frameLen int, frameStep int, emotionFiles map[string][]string) error {
	trainSet, err := UnmarshallKNNEeg(trainSetFilename)
	if err != nil {
		return err
	}
	_, trainVar := GetμAndσTagged(trainSet)

	fileKeys := make([]string, 0, len(emotionFiles))
	for k := range emotionFiles {
		fileKeys = append(fileKeys, k)
	}
	sort.Strings(fileKeys)

	correctFiles := make(map[string]int, len(fileKeys))
	correctVectors := make(map[string]int, len(fileKeys))
	sumVectors := make(map[string]int, len(fileKeys))

	for _, emotion := range fileKeys {
		for _, f := range emotionFiles[emotion] {
			fmt.Printf("%s\t", emotion)
			vec := GetFourierForFile(f, 19, frameLen, frameStep)
			average := GetAverage(bucketSize, frameStep, len(vec))
			averaged := AverageSlice(vec, average)

			boolCorrect, correctVector, sumVector := testKNN(emotion, fileKeys, averaged, trainSet, trainVar)
			correctFiles[emotion] += boolCorrect
			correctVectors[emotion] += correctVector
			sumVectors[emotion] += sumVector
		}
	}
	sort.Strings(fileKeys)
	fmt.Printf("\tCorrectFiles\tCorrectVectors\n")
	for _, emotion := range fileKeys {
		fmt.Printf("%s\t%f\t%f\n", emotion, float64(correctFiles[emotion])/float64(len(emotionFiles[emotion])), float64(correctVectors[emotion])/float64(sumVectors[emotion]))
	}

	return nil
}

func ClassifyGMM(trainSetFilename string, bucketSize int, frameLen int, frameStep int, emotionFiles map[string][]string) error {
	trainSet, err := GetEGMs(trainSetFilename)
	if err != nil {
		return err
	}

	fileKeys := make([]string, 0, len(emotionFiles))
	for k := range emotionFiles {
		fileKeys = append(fileKeys, k)
	}

	correctFiles := make(map[string]int, len(fileKeys))
	correctVectors := make(map[string]int, len(fileKeys))
	sumVectors := make(map[string]int, len(fileKeys))

	for _, emotion := range fileKeys {
		for _, f := range emotionFiles[emotion] {
			vec := GetFourierForFile(f, 19, frameLen, frameStep)
			average := GetAverage(bucketSize, frameStep, len(vec))
			averaged := AverageSlice(vec, average)

			boolCorrect, correctVector, sumVector := TestGMM(emotion, fileKeys, averaged, trainSet)
			correctFiles[emotion] += boolCorrect
			correctVectors[emotion] += correctVector
			sumVectors[emotion] += sumVector
		}
	}
	sort.Strings(fileKeys)
	fmt.Printf("\tCorrectFiles\tCorrectVectors\n")
	for _, emotion := range fileKeys {
		fmt.Printf("%s\t%f\t%f\n", emotion, float64(correctFiles[emotion])/float64(len(emotionFiles[emotion])), float64(correctVectors[emotion])/float64(sumVectors[emotion]))
	}

	return nil
}
