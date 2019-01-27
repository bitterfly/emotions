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
		panic("go run main.go <k> <outut_dir> <emotion1.csv [emotion2.csv...]>")
	}

	// k, err := strconv.Atoi(os.Args[1])
	// if err != nil {
	// 	panic(err)
	// }

	// outputDir := fmt.Sprintf("%s_k%d", os.Args[2], k)
	// if _, err := os.Stat(outputDir); os.IsNotExist(err) {
	// 	os.Mkdir(outputDir, 0775)
	// }

	// filenames := make([]string, len(os.Args)-2, len(os.Args)-2)
	// egms := make([]emotions.EmotionGausianMixure, len(os.Args)-2, len(os.Args)-2)

	// for i := 3; i < len(os.Args); i++ {

	// 	filename := filepath.Base(os.Args[i])
	// 	name := filename[0 : len(filename)-len(filepath.Ext(filename))]
	// 	filenames[i-2] = filepath.Join(outputDir, name+".gmm")

	// 	egms[i-2] = emotions.EmotionGausianMixure{
	// 		Emotion: name,
	// 		GM:      getGaussians(os.Args[i], k, 19),
	// 	}
	// }

	// for i := 0; i < len(filenames); i++ {
	// 	bytes, err := json.Marshal(egms[i])
	// 	if err != nil {
	// 		panic(fmt.Sprintf("Error when marshaling %s\n", filenames[i]))
	// 	}
	// 	ioutil.WriteFile(filenames[i], bytes, 0644)
	// }
}
