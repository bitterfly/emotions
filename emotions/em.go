package emotions

import (
	"fmt"
	"math"
	"os"
)

// Gaussian represent a single gaussian
type Gaussian struct {
	Phi          float64
	Expectations []float64
	Variances    []float64
}

// GaussianMixture represent a mixture of gaussians
type GaussianMixture = []Gaussian

func zeroMixture(g GaussianMixture, K int) {
	for k := 0; k < K; k++ {
		zero(&g[k].Expectations)
		zero(&g[k].Variances)
		g[k].Phi = 0.0
	}
}

type EmotionGausianMixure struct {
	Emotion string
	GM      GaussianMixture
}

// GMM returns the k gaussian mixures for the given data
func GMM(mfccsFloats [][]float64, k int) GaussianMixture {
	X, expectations, variances, numInCluster := KMeans(mfccsFloats, k)
	fmt.Fprintf(os.Stderr, "\n==============EM================\n")

	gmixture := make(GaussianMixture, k, k)

	for i := 0; i < k; i++ {
		gmixture[i] = Gaussian{
			Phi:          float64(numInCluster[i]) / float64(len(X)),
			Expectations: expectations[i],
			Variances:    variances[i],
		}
	}

	return em(X, k, gmixture)
}

func em(X []MfccClusterisable, k int, gMixture GaussianMixture) GaussianMixture {
	prevLikelihood := 0.0
	likelihood := 0.0
	step := 0
	for step < 200 {
		w := make([][]float64, len(X), len(X))
		var sum float64
		maximums := make([]float64, len(X), len(X))

		for i := 0; i < len(maximums); i++ {
			maximums[i] = math.Inf(-1)
		}

		for i := 0; i < len(X); i++ {
			w[i] = make([]float64, k, k)

			for j := 0; j < k; j++ {
				w[i][j] = math.Log(gMixture[j].Phi) + N(X[i].coefficients, gMixture[j].Expectations, gMixture[j].Variances)

				if maximums[i] < w[i][j] {
					maximums[i] = w[i][j]
				}
			}

			sum = 0
			for j := 0; j < k; j++ {
				if i == 10 {
				}
				if w[i][j] < maximums[i]-10 {
					w[i][j] = 0
				} else {
					w[i][j] = math.Exp(w[i][j] - maximums[i])
					sum += w[i][j]
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

		zeroMixture(gMixture, k)

		// Ecpectations
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				add(&gMixture[j].Expectations, multiplied(X[i].coefficients, w[i][j]))
			}
		}
		for j := 0; j < k; j++ {
			divide(&(gMixture[j].Expectations), N[j])
		}

		// Variances
		for i := 0; i < len(X); i++ {
			for j := 0; j < k; j++ {
				diagonal := minused(X[i].coefficients, gMixture[j].Expectations)
				square(&diagonal)

				add(&gMixture[j].Variances, multiplied(diagonal, w[i][j]))
			}
		}

		// Phi and 1/Nk
		for j := 0; j < k; j++ {
			divide(&(gMixture[j].Variances), N[j])
			gMixture[j].Phi = N[j] / float64(len(X))
		}

		likelihood = logLikelihood(X, k, gMixture)

		if math.IsNaN(likelihood) {
			panic("AAAA")
		}

		if epsDistance(likelihood, prevLikelihood, 0.00001) {
			break
		}

		prevLikelihood = likelihood
		step++
	}

	fmt.Fprintf(os.Stderr, "EM: Break on step: %d with likelihood: %f\n===================================================\n", step, likelihood)
	return gMixture
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

// EvaluateVector returns the likelihood a given vector
func EvaluateVector(X []float64, k int, g GaussianMixture) float64 {
	return logLikelihoodFloat(X, k, g)
}

func N(xi []float64, expectation []float64, variance []float64) float64 {
	var exp float64
	for i := 0; i < len(xi); i++ {
		exp += (xi[i] - expectation[i]) * (xi[i] - expectation[i]) / variance[i]
	}

	return -0.5 * (exp + float64(len(xi))*math.Log(2.0*math.Pi) + math.Log(getDeterminant(variance)))
	// return log of this
	// return math.Exp(-0.5*exp) / math.Sqrt(math.Pow(2*math.Pi, float64(len(xi)))*getDeterminant(variance))
}

func logLikelihoodFloat(X []float64, k int, g GaussianMixture) float64 {
	sum := 0.0
	for j := 0; j < k; j++ {
		sum += g[j].Phi * math.Exp(N(X, g[j].Expectations, g[j].Variances))
	}
	return math.Log(sum)
}

// sum_i log(sum_j phi_j * N(x[i], m[k], s[k]))

func logLikelihood(X []MfccClusterisable, k int, g GaussianMixture) float64 {
	sum := 0.0
	for i := 0; i < len(X); i++ {
		sum += logLikelihoodFloat(X[i].coefficients, k, g)
	}
	return sum
}
