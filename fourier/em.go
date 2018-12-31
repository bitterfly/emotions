package fourier

import "math"

func EM(k int, expectations [][]float64, variances [][]float64, numInCluster []int, X []MfccClusterisable) {
	// start with k-mean clusters expectations and variances
	// and choose

	// ϕ is 1/k

	phi := make([]float64, k, k)
	for i := 0; i < k; i++ {
		phi[i] = float64(numInCluster[i]) / float64(len(X))
	}

	//Ns := getNs(k,expectations, variances, X)
	newExpectations := make([][]float64, k, k)
	newVariances := make([][]float64, k, k)

	for i := 0; i < k; i++ {
		copy(newExpectations[i], expectations[i])
		copy(newVariances[i], variances[i])
	}

	w := make([][]float64, len(X), len(X))

	//E step
	//compute the weights

	for i := 0; i < len(X); i++ {
		sum := 0.0

		for j := 0; j < k; j++ {
			w[i] = make([]float64, k, k)
			w[i][j] = phi[j] * N(X[i].coefficients, expectations[j], variances[j])
		}
		divide(&w[i], sum)
	}

	sumOfWeights := make([]float64, k, k)
	for i := 0; i < len(X); i++ {
		for j := 0; k < k; j++ {
			sumOfWeights[j] += w[i][j]
		}
	}

	// compute the μ Σ ɸ

	for j := 0; j < k; j++ {
		zero(&(expectations[j]))
		zero(&(variances[j]))
	}
	zero(&phi)

	for i := 0; i < len(X); i++ {
		for j := 0; j < k; j++ {
			add(&expectations[j], multiplied(X[i].coefficients, w[i][j]))
		}
	}

	for i := 0; i < len(X); i++ {
		for j := 0; j < k; j++ {

		}
	}

	for j := 0; j < k; j++ {
		divide(&(expectations[j]), sumOfWeights[j])
		divide(&(variances[j]), sumOfWeights[j])
		phi[j] = sumOfWeights[j] / float64(len(X))
	}

}

func getDeterminant(variance []float64) float64 {
	det := 1.0
	for i := 0; i < len(variance); i++ {
		det *= variance[i]
	}
	return det
}

func N(xi []float64, expectation []float64, variance []float64) float64 {
	exp := 0.0
	for i := 0; i < len(xi); i++ {
		exp += (xi[i] - expectation[i]) * (xi[i] - expectation[i]) / variance[i]
	}

	return float64(1) / math.Sqrt(math.Pow(2.0*math.Pi, float64(len(xi)))*getDeterminant(variance)) * math.Exp(-0.5*exp)
}

// func getNs(k int, expectations [][]float64, variances [][]float64, X []MfccClusterisable) [][]float64 {
// 	Ns := make([][]float64, k, k)

// 	for

// 	return nil
// }
