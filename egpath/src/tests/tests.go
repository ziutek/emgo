package tests

type Int int

func (i Int) Get() int {
	return int(i)
}

func (i *Int) Set(v int) {
	*i = Int(v)
}

type T struct {
	Int
}
