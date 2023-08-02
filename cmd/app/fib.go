package app

func Fib(n uint) (uint64, error) {
	if n <= 1 {
		return uint64(n), nil
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