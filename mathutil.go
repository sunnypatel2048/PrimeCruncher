package main

import "math/big"

// isPrime returns true if n is a prime number, and false otherwise
func isPrime(n uint64) bool {
	if n < 2 {
		return false
	}

	bigN := big.NewInt(int64(n))

	return bigN.ProbablyPrime(20)
}
