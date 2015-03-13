package tests

type Int int

func (i Int) Get() int {
	return int(i)
}

func (i *Int) Set(v1, v2 int) {
	*i = Int(v1 + v2)
}

type T struct {
	Int
}

type Terr struct {
	T
	error
}

func f(t *Terr) string {
	s := t.error.Error()
	s = t.Error()
	return s
}
