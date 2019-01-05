package fourier

import (
	"fmt"
	"math"
)

func GMM(mfccsFloats [][]float64, k int) {
	X, expectations, variances, numInCluster := KMeans(mfccsFloats, k)
	em(expectations, variances, numInCluster, X, k)
}

func em(expectations [][]float64, variances [][]float64, numInCluster []int, X []MfccClusterisable, k int) {
	// start with k-mean clusters expectations, variances and we use the number of points in a cluster for phi
	// and choose

	phi := make([]float64, k, k)
	for i := 0; i < k; i++ {
		phi[i] = float64(numInCluster[i]) / float64(len(X))
	}

	// weights

	fmt.Printf("X: %v\n", X)

	prev_likelihood := 0.0
	for step := 0; step < 100; step++ {
		// fmt.Printf("phi: %v\n", phi)
		// fmt.Printf("expectations: %v\n", expectations)
		// fmt.Printf("variances: %v\n", variances)

		// fmt.Printf("len(X): %d\n", len(X))

		w := make([][]float64, len(X), len(X))
		maximums := make([]float64, len(X), len(X))
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				w[i] = make([]float64, k, k)
				w[i][j] = math.Log(phi[j]) + N(X[i].coefficients, expectations[j], variances[j])
				if maximums[i] < w[i][j] {
					maximums[i] = w[i][j]
				}
			}

			var sum float64
			for j := 0; j < k; j++ {
				if w[i][j] < maximums[i]-10 {
					w[i][j] = 0
				} else {
					sum += math.Exp(w[i][j] - maximums[i])
					w[i][j] = math.Exp(w[i][j] - maximums[i])
				}
			}
			divide(&w[i], sum)
		}

		N := make([]float64, k, k)
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				N[j] += w[i][j]
			}
		}

		// fmt.Printf("weights: %v\n", w)
		// fmt.Printf("Ns: %v\n", N)
		// compute the μ Σ ɸ
		for j := 0; j < k; j++ {
			zero(&(expectations[j]))
			zero(&(variances[j]))
		}
		zero(&phi)

		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				add(&expectations[j], multiplied(X[i].coefficients, w[i][j]))

				// fmt.Printf("Expectation[%d]: %f\n", j, expectations[j])
				diagonal := minused(X[i].coefficients, expectations[j])
				square(&diagonal)

				add(&variances[j], multiplied(diagonal, w[i][j]))
			}
		}

		for j := 0; j < k; j++ {
			divide(&(expectations[j]), N[j])
			divide(&(variances[j]), N[j])
			phi[j] = N[j] / float64(len(X))
		}

		likelihood := logLikelihood(X, phi, expectations, variances, k)
		fmt.Printf("Step: %d, Log likelihood: %f\n", step, likelihood)

		if epsDistance(likelihood, prev_likelihood, 0.00001) {
			fmt.Printf("Break on step: %d\n", step)
			break
		}
		prev_likelihood = likelihood
	}

}

func epsDistance(a, b, e float64) bool {
	return (a-b < e && a-b > -e)
}

func getDeterminant(variance []float64) float64 {
	det := 1.0
	for i := 0; i < len(variance); i++ {
		det *= variance[i]
	}
	return det
}

func N(xi []float64, expectation []float64, variance []float64) float64 {
	var exp float64
	for i := 0; i < len(xi); i++ {
		exp += (xi[i] - expectation[i]) * (xi[i] - expectation[i]) / variance[i]
	}

	// return math.Exp(-0.5*exp) / math.Sqrt(math.Pow(2.0*math.Pi, float64(len(xi)))*getDeterminant(variance))
	// return log of this
	return -0.5 * (exp - float64(len(xi))*math.Log(2*math.Pi) - math.Log(getDeterminant(variance)))
}

func logLikelihood(X []MfccClusterisable, phi []float64, expectations [][]float64, variances [][]float64, k int) float64 {
	sum := 0.0
	for i := 0; i < len(X); i++ {
		for j := 0; j < k; j++ {
			sum += math.Log(phi[j] * math.Exp(N(X[i].coefficients, expectations[j], variances[j])))
		}
	}
	return sum
}
