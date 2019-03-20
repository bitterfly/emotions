package emotions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
)

type Tagged struct {
	Tag  string
	Data [][]float64
}

func UnmarshallEeg(filename string) ([]Tagged, error) {
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

func GetμAndσTagged(tagged []Tagged) ([]float64, []float64) {
	featureVectorSize := len(tagged[0].Data[0])
	fmt.Fprintf(os.Stderr, "FeatureVectorSize: %d\n", featureVectorSize)

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

	fmt.Fprintf(os.Stderr, "numVectors: %d\n", numVectors)
	for j := 0; j < featureVectorSize; j++ {
		expectation[j] /= float64(numVectors)
		expectationSquared[j] /= float64(numVectors)

		variances[j] = expectationSquared[j] - expectation[j]*expectation[j]
	}

	fmt.Fprintf(os.Stderr, "Expectations: %d\n", len(expectation))
	fmt.Fprintf(os.Stderr, "Variances: %d\n", len(variances))
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

func findMostCommonTag(vectors [][]float64, trainSet []Tagged, trainVar []float64) (map[string]int, string) {
	freqMap := make(map[string]int)
	for t := range trainSet {
		freqMap[trainSet[t].Tag] = 0
	}

	for v := range vectors {
		freqMap[findClosest(vectors[v], trainSet, trainVar)]++
	}

	max := 0
	maxe := ""
	for e, f := range freqMap {
		if f > max {
			max = f
			maxe = e
		}

	}
	return freqMap, maxe
}

func KNNOne(trainSet []Tagged, trainVar []float64, bucketSize int, frameLen int, frameStep int, filename string) float64 {
	vec := GetFourierForFile(filename, 19, frameLen, frameStep)
	average := GetAverage(bucketSize, frameLen, len(vec))
	averaged := AverageSlice(vec, average)
	_, mc := findMostCommonTag(averaged, trainSet, trainVar)
	switch mc {
	case "eeg-neutral":
		return 0
	case "eeg-positive":
		return 1
	case "eeg-negative":
		return -1
	}
	return 0
}

func KNN(bucketSize int, frameLen int, frameStep int, trainSetFilename string, emotionFiles map[string][]string) error {
	trainSet, err := UnmarshallEeg(trainSetFilename)
	if err != nil {
		return err
	}

	for i := range trainSet {
		fmt.Fprintf(os.Stderr, "Tag: %s, len: %d x %d\n", trainSet[i].Tag, len(trainSet[i].Data), len(trainSet[i].Data[0]))
	}

	_, trainVar := GetμAndσTagged(trainSet)

	fileKeys := make([]string, 0, len(emotionFiles))
	for k := range emotionFiles {
		fileKeys = append(fileKeys, k)
	}
	sort.Strings(fileKeys)

	counts := make(map[string]int)

	for _, emotion := range fileKeys {
		counts[emotion] = 0
		for _, f := range emotionFiles[emotion] {
			fmt.Printf("%s\t", emotion)
			vec := GetFourierForFile(f, 19, frameLen, frameStep)
			average := GetAverage(bucketSize, frameLen, len(vec))
			averaged := AverageSlice(vec, average)

			dict, _ := findMostCommonTag(averaged, trainSet, trainVar)
			keys := make([]string, 0, len(dict))
			for k := range dict {
				keys = append(keys, k)
			}

			sort.Strings(keys)

			for _, k := range keys {
				fmt.Printf("%d\t", dict[k])
			}
			fmt.Printf("\n")
		}
	}

	return nil
}
