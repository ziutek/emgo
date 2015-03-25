package tests

var v = [][]int{{11: 3}, {13: Int()}}

func Int() int {
	return 1
}

type S struct {
	a int
	b byte
	s string
}

var (
	s1 = &S{1, 2, "foo"}
	s2 = S{b: 3, s: "bar"}
	s3 = S{a: Int()}
	s4 = &S{a: Int()}
)
