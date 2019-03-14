package emotions

import (
	"fmt"
	"math"
)

type Tagged struct {
	Tag  string
	Data [][]float64
}

func findDistance(x [][]float64, y [][]float64) float64 {
	sum := 0.0
	for i := 0; i < len(x); i++ {
		sum += euclidianDistance(x[i], y[i], nil)
	}

	return sum
}

func KNN(trainSetFilename string, testVectorFile string) {
	trainSet := getEegTrainingSet(trainSetFilename)
	testVector := GetFeatureVector(testVectorFile, 19)

	min := math.Inf(42)
	var argmin string
	for _, ts := range trainSet {
		dist := findDistance(ts.Data, testVector)
		fmt.Printf("%s %f\n", ts.Class, dist)
		if dist < min {
			min = dist
			argmin = ts.Class
		}
	}

	fmt.Printf("Best: %s\n", argmin)

}
