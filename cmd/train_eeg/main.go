package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func getCluster(filename string, k int, elNum int) emotions.GaussianMixture {
	fv := emotions.GetFeatureVector(filename, elNum)
	emotions.KMeans(fv, k)
	return nil
}

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <output-file> <emotion1.csv [emotion2.csv...]>")
	}

	outputFile := os.Args[1]
	emotions.SaveEegTrainingSet(os.Args[2:len(os.Args)], outputFile)
}
