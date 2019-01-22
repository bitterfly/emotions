package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bitterfly/emotions/emotions"
)

func readEmotion(filename string, k int) ([][]float64, emotions.GaussianMixture) {
	wf, _ := emotions.Read(filename, 0, 0.97)

	allMfccs := emotions.MFCCs(wf, 13, 23)

	testingIndices := make(map[int]struct{})
	rand.Seed(time.Now().UTC().UnixNano())
	for len(testingIndices) < len(allMfccs)*20.0/100.0 {
		ind := rand.Intn(len(allMfccs))
		if _, ok := testingIndices[ind]; !ok {
			testingIndices[ind] = struct{}{}
		}
	}

	test := make([][]float64, 0, len(allMfccs))
	train := make([][]float64, 0, len(allMfccs)-len(testingIndices))

	for i, m := range allMfccs {
		if _, ok := testingIndices[i]; ok {
			test = append(test, m)
		} else {
			train = append(train, m)
		}
	}

	gMixture := emotions.GMM(train, k)
	return test, gMixture
}

func main() {
	k := 5

	happyTest, happyGM := readEmotion(os.Args[1], k)
	sadTest, sadGM := readEmotion(os.Args[2], k)

	happy := 0
	sad := 0
	for _, m := range happyTest {
		if emotions.EvaluateVector(m, k, happyGM) > emotions.EvaluateVector(m, k, sadGM) {
			happy++
		} else {
			sad++
		}
	}

	fmt.Printf("In happy(%d): Happy: %d, Fear: %d\n", len(happyTest), happy, sad)

	happy = 0
	sad = 0
	for _, m := range sadTest {
		if emotions.EvaluateVector(m, k, happyGM) > emotions.EvaluateVector(m, k, sadGM) {
			happy++
		} else {
			sad++
		}
	}

	fmt.Printf("In sad(%d): Happy: %d, Fear: %d\n", len(sadTest), happy, sad)
}
