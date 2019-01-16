package main

import (
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
	// test_happy := make([][]float64, 40, 40)
	// // test_sad := make([][]float64, 40, 40)
	// for i := 0; i < 40; i++ {
	// 	test_happy[i] = mfcc_happy[i][0:6]
	// 	// test_sad[i] = mfcc_sad[i][0:6]
	// }

	_, _, _ = fourier.GMM(mfcc_happy, k)
	// fmt.Printf("Happy GMM\n")
	// happy_phi, happy_exp, happy_var := fourier.GMM(test_happy, k)

	// fmt.Printf("Sad GMM\n")
	// sad_phi, sad_exp, sad_var := fourier.GMM(test_sad, k)

	// happy := 0
	// sad := 0
	// for _, m := range test_sad {
	// 	fmt.Printf("With happy: %f\n", fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k))
	// 	fmt.Printf("With sad: %f\n", fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k))
	// 	if fourier.EvaluateVector(m, happy_phi, happy_exp, happy_var, k) > fourier.EvaluateVector(m, sad_phi, sad_exp, sad_var, k) {
	// 		happy++
	// 	} else {
	// 		sad++
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
