package testy

func F(a, b int) int {
	return func(x int) int {
		return x + b
	}(a)
}
