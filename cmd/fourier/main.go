package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	// dirname := os.Args[1]
	// files, _ := ioutil.ReadDir(dirname)
	wf1, _ := fourier.Read(os.Args[1], 0, 0.97)
	wf2, _ := fourier.Read(os.Args[2], 0, 0.97)

	mfcc_happy := fourier.MFCCs(wf1, 13, 23)
	mfcc_sad := fourier.MFCCs(wf2, 13, 23)

	k := 5

	happy_indices := make(map[int]struct{})
	sad_indices := make(map[int]struct{})

	rand.Seed(time.Now().UTC().UnixNano())
	for len(happy_indices) < len(mfcc_happy)*20.0/100.0 {
		ind := rand.Intn(len(mfcc_happy))
		if _, ok := happy_indices[ind]; !ok {
			happy_indices[ind] = struct{}{}
		}
	}

	rand.Seed(time.Now().UTC().UnixNano())
	for len(sad_indices) < len(mfcc_sad)*20.0/100.0 {
		ind := rand.Intn(len(mfcc_happy))
		if _, ok := sad_indices[ind]; !ok {
			sad_indices[ind] = struct{}{}
		}
	}

	happy_test := make([][]float64, 0, len(happy_indices))
	happy_train := make([][]float64, 0, len(mfcc_happy)-len(happy_indices))

	sad_test := make([][]float64, 0, len(sad_indices))
	sad_train := make([][]float64, 0, len(mfcc_sad)-len(sad_indices))

	for i, m := range mfcc_happy {
		if _, ok := happy_indices[i]; ok {
			happy_test = append(happy_test, m)
		} else {
			happy_train = append(happy_train, m)
		}
	}

	for i, m := range mfcc_sad {
		if _, ok := sad_indices[i]; ok {
			sad_test = append(sad_test, m)
		} else {
			sad_train = append(sad_train, m)
		}
	}

	fmt.Printf("Happy_indices: %d\nHappy_test: %d\n", len(happy_indices), len(happy_test))
	fmt.Printf("Sad_indices: %d\nSad_test: %d\n", len(sad_indices), len(sad_test))

	// _, _, _ = fourier.GMM(mfcc_happy, k)
	fmt.Printf("Happy GMM train with %d\n", len(happy_train))
	happy_phi, happy_exp, happy_var := fourier.GMM(happy_train, k)

	fmt.Printf("Sad GMM train with %d\n", len(sad_train))
	sad_phi, sad_exp, sad_var := fourier.GMM(sad_train, k)

	happy := 0
	sad := 0
	for _, m := range happy_test {
		// fmt.Printf("With happy: %f\n", fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k))
		// fmt.Printf("With sad: %f\n", fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k))
		if fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k) > fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k) {
			happy++
		} else {
			sad++
		}
	}

	fmt.Printf("In happy(%d): Happy: %d, Fear: %d\n", len(happy_indices), happy, sad)

	happy = 0
	sad = 0
	for _, m := range sad_test {
		// fmt.Printf("With happy: %f\n", fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k))
		// fmt.Printf("With sad: %f\n", fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k))
		if fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k) > fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k) {
			happy++
		} else {
			sad++
		}
	}

	fmt.Printf("In sad(%d): Happy: %d, Fear: %d\n", len(sad_indices), happy, sad)
}
