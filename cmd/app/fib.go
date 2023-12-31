package app

import "fmt"

func Fib(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
	}

	if n > 93 {
		return 0, fmt.Errorf(
			"unsupported Fibonacci number %d: too large", n,
		)
	}

	var (
		n2 uint64 = 0
		n1 uint64 = 1
	)

	for i := uint(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1, nil
}
