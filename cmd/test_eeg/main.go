package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitterfly/emotions/emotions"
)

func UnmarshallEeg(filename string) (error, []emotions.Tagged) {
	var tagged []emotions.Tagged
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err, nil
	}

	err = json.Unmarshal(bytes, &tagged)
	if err != nil {
		return err, nil
	}

	return nil, tagged
}

func main() {
	if len(os.Args) < 6 {
		panic("go run main.go <eeg-train-file> --eeg-positive eeg_pos1.txt [eeg_pos2.txt...] --eeg-negative eeg_neg1.txt [eeg_neg2.txt...] --eeg-neutral eeg_neu1.txt [eeg_neu2.txt...]")
	}

	train_file := os.Args[1]
	err, trainSet := UnmarshallEeg(train_file)
	if err != nil {
		panic(err.Error())
	}

	for i := range trainSet {
		fmt.Printf("Tag: %s, len: %d\n", trainSet[i].Tag, len(trainSet[i].Data))
	}
	// emotions.KNN(os.Args[1], os.Args[2])
}
