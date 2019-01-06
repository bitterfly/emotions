package main

import (
	"fmt"
	"os"

	"github.com/bitterfly/emotions/fourier"
)

func main() {
	// dirname := os.Args[1]
	// files, _ := ioutil.ReadDir(dirname)
	wf1, _ := fourier.Read(os.Args[1], 0, 0.97)
	// wf2, _ := fourier.Read(os.Args[2], 0, 0.97)

	mfcc_happy := fourier.MFCCs(wf1, 13, 23)
	// mfcc_sad := fourier.MFCCs(wf2, 13, 23)

	k := 3

	test_mfccs := make([][]float64, 40, 40)
	for i := 0; i < 40; i++ {
		test_mfccs[i] = make([]float64, 7, 7)
		test_mfccs[i][0] = mfcc_happy[i][0]
		test_mfccs[i][1] = mfcc_happy[i][1]
		test_mfccs[i][2] = mfcc_happy[i][2]
		test_mfccs[i][3] = mfcc_happy[i][3]
		test_mfccs[i][4] = mfcc_happy[i][4]
		test_mfccs[i][5] = mfcc_happy[i][5]
		test_mfccs[i][6] = mfcc_happy[i][6]
	}
	fmt.Printf("Happy GMM\n")
	_, _, _ = fourier.GMM(test_mfccs, k)

	// happy_phi, happy_exp, happy_var := fourier.GMM(mfcc_happy, k)

	// fmt.Printf("Sad GMM\n")
	// sad_phi, sad_exp, sad_var := fourier.GMM(mfcc_sad, k)

	// happy := 0
	// sad := 0
	// for _, m := range mfcc_happy {
	// 	if fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k) > fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k) {
	// 		happy += 1
	// 	} else {
	// 		sad += 1
	// 	}
	// }

	// fmt.Printf("Happy: %d, Sad: %d\n", happy, sad)
	// indices := make([]int, len(files), len(files))
	// mfccs := make([][]float64, 0, len(files)*1000)
	// names := make([]string, len(files), len(files))
	// for i, f := range files {
	// 	names[i] = f.Name()[0 : len(f.Name())-4]
	// 	wf, err := fourier.Read(path.Join(dirname, f.Name()), 0, 0.97)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	mfccs = append(mfccs, fourier.MFCCs(wf, 13, 23)...)
	// 	indices[i] = len(mfccs) - 1
	// }

	// fmt.Printf("%d\n", len(mfccs))
	// fmt.Printf("Kmeans: \n")
	// points := [][]float64{
	// 	[]float64{1, 0, 1.002},
	// 	[]float64{0, 1, 1.02},
	// 	[]float64{5.23, 5.22, 5.10},
	// 	[]float64{6.23, 6.4, 6.33},
	// 	[]float64{7.11, 7.22, 7.13},
	// }

	// pf, _ := os.Create("/tmp/points.csv")
	// cf, _ := os.Create("/tmp/centroids.csv")
	// defer pf.Close()
	// defer cf.Close()

	// fmt.Fprintf(pf, "X, Y, Z\n")
	// fmt.Fprintf(cf, "X, Y, Z\n")

	// for i, ms := range ms {
	// 	fmt.Fprintf(pf, "%f, %f, %d\n", points[i][0], points[i][1], ms.GetCluster())
	// }

	// for i, cc := range c {
	// 	fmt.Fprintf(cf, "%f, %f, %d\n", cc[0], cc[1], i)
	// }
}
