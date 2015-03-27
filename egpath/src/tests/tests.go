package tests

/*
var V = [][]int{{11: 3}, {13: Int()}}

func Int() int {
	return 1
}

type S struct {
	a int
	b byte
	s string
}

var (
	S1 = &S{1, 2, "foo"}
	S2 = S{b: 3, s: "bar"}
	S3 = S{a: Int()}
	S4 = &S{a: Int()}
)

var A = []int{1, 2, 30: 3, 4, 20: 5, 6}
*/

func f(v interface{}) (int, bool) {
	i, ok := v.(int)
	return i, ok
}