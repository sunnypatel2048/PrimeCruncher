package main

import "math/big"

func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}

	// convert n to big.Int for use with math/big package
	bigN := big.NewInt(int64(n))

	return bigN.ProbablyPrime(20)
}
