package testy

func f(a, b int) int {
	return func(x int) int {
		return x + b
	}(a)
}
