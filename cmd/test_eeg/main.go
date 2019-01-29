package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <eeg-train-file> <emotion1.csv [emotion2.csv...]>")
	}

	trainSetFile := os.Args[1]
	// testVectorFile := os.Args[2]
	// emotions.KNN(trainSetFile, testVectorFile)
	emotions.PlotEeg(trainSetFile, "/tmp/ts.png")
}
