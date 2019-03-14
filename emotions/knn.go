package emotions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
)

type Tagged struct {
	Tag  string
	Data [][]float64
}

func unmarshallEeg(filename string) ([]Tagged, error) {
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

func getμAndσTagged(tagged []Tagged) ([]float64, []float64) {
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

func findMostCommonTag(vectors [][]float64, trainSet []Tagged, trainVar []float64) string {
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
	return maxe
}

func KNN(trainSetFilename string, emotionFiles map[string][]string) error {
	trainSet, err := unmarshallEeg(trainSetFilename)
	if err != nil {
		return err
	}

	for i := range trainSet {
		fmt.Fprintf(os.Stderr, "Tag: %s, len: %d x %d\n", trainSet[i].Tag, len(trainSet[i].Data), len(trainSet[i].Data[0]))
	}

	_, trainVar := getμAndσTagged(trainSet)

	counts := make(map[string]int)

	for emotion, files := range emotionFiles {
		counts[emotion] = 0
		for f := range files {
			vec := GetFourierForFile(files[f], 19)
			mostCommonTag := findMostCommonTag(vec, trainSet, trainVar)
			if mostCommonTag == emotion {
				counts[emotion]++
			}
		}
	}

	sum := 0
	files := 0
	for e, c := range counts {
		fmt.Printf("%s: %.3f%%\n", e, float64(c)*100/float64(len(emotionFiles[e])))
		sum += c
		files += len(emotionFiles[e])
	}
	fmt.Printf("All: %.3f%%\n", float64(sum)*100/float64(files))

	return nil
}
