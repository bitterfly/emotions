#!/usr/bin/env python3
import numpy as np
import scipy.stats as stats
def partition(data):
    d = dict()
    for label, items in data:
        if label not in d:
            d[label] = []
        d[label].append(items)
    return [np.array(v) for _, v in sorted(d.items())]
def initial_variances(data):
    return [np.apply_along_axis(np.var, 0, table) for table in data]
def initial_expectations(data):
    return [np.apply_along_axis(np.mean, 0, table) for table in data]
def initial_pis(data):
    all = sum(len(table) for table in data)
    return [len(table) / all for table in data]
def getN(x, expectation, variance):
    pos = np.array([x])
    bell = stats.multivariate_normal(expectation, np.diag(variance))
    return bell.pdf(pos)
def initialise(data):
    part = partition(data)
    X = np.array([x[1] for x in data])
    return (X, (initial_expectations(part), initial_variances(part), initial_pis(part)))
def semi_flatten(arrs):
    foo = []
    for arr in arrs:
        foo += list(arr)
    return np.array(foo)
def em(things):
    (X, (expectations, variances, pis)) = things
    w = np.zeros((len(X), len(pis)))
    for n in range(len(X)):
        SUM = sum( ( pis[k] * getN(X[n], expectations[k], variances[k]) for k in range(len(pis)) ) )
        for k in range(len(pis)):
            w[n][k] = pis[k] * getN(X[n], expectations[k], variances[k]) / SUM
    N = np.array( [ sum( (w[n][k] for n in range(len(X))) ) for k in range(len(pis)) ] )

    n_expectations = np.array(
        [ np.sum( (w[n][k] * X[n] for n in range(len(X))), axis=0 )
                / N[k]
                            for k in range(len(pis)) ] )
    n_variances = np.array(
        [ np.sum( (w[n][k] * ((X[n] - n_expectations[k]) ** 2) for n in range(len(X))), axis=0 )
                / N[k]
                            for k in range(len(pis)) ] )
    n_pis = np.array( [N[k] / len(X) for k in range(len(pis))] )

   
    return (X, (n_expectations, n_variances, n_pis))

def stab_em(things, epsilon=0.00001, steps=100):
    for step in range(steps):
        print("step ", step)
        old = logLikelihood(things)
        things = em(things)
        new = logLikelihood(things)
        if new - old < epsilon:
            print("break on step {}".format(step))
            break
    return things
def logLikelihood(things):
    (X, (expectations, variances, pis)) = things
    return sum (
                (np.log(sum(
                            (pis[k] * getN(X[n], expectations[k], variances[k])
                            for k in range(len(pis)) ))
                )
           for n in range(len(X)) ) )
def show(things):
    (X, (expectations, variances, pis)) = things
    print("expectations: \n", expectations, "\n")
    print("variances: \n", variances, "\n")
    print("pis: \n", pis, "\n")
    print("log likelihood: ", logLikelihood(things))
