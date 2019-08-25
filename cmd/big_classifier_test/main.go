package main

import (
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func main() {
	if len(os.Args) < 3 {
		panic("go run main.go <gmm-dir> <input-file>\n<input-file>:<emotion>	<csv-file>")
	}

	gmmDir := os.Args[1]
	speechFiles, eegFiles, _, err := emotions.ParseArgumentsFromFile(os.Args[2], true)
	if err != nil {
		panic(err)
	}

	emotions.ClassifyGMMConcat(gmmDir, speechFiles, eegFiles)
}
