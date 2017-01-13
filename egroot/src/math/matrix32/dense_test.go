package matrix32

import (
	"fmt"
	"math/rand"
	"testing"
)

func TestDenseI(t *testing.T) {
	m := MakeDense(11, 11)
	m.SetIdentity()
	numrow, numcol := m.Size()
	for i := 0; i < numrow; i++ {
		for k := 0; k < numcol; k++ {
			a := m.Get(i, k)
			if i == k {
				if a != 1 {
					t.Fatalf("(%d, %d) == %f != 1", i, k, a)
				}
			} else {
				if a != 0 {
					t.Fatalf("(%d, %d) == %f != 0", i, k, a)
				}
			}
		}
	}
}

func TestDenseSlice(t *testing.T) {
	m := MakeDense(5, 6)
	m.setRandom(-1, 1)
	start, stop := 1, 4

	s := m.Rows(start, stop)
	if s.NumRow() != stop-start || s.NumCol() != m.NumCol() {
		t.Fatal("Rows: wrong size")
	}
	for i := 0; i < s.NumRow(); i++ {
		for k := 0; k < s.NumCol(); k++ {
			if s.Get(i, k) != m.Get(i+start, k) {
				t.Fatal("Rows: elements do not match")
			}
		}
	}

	s = m.Cols(start, stop)
	if s.NumRow() != m.NumRow() || s.NumCol() != stop-start {
		t.Fatal("Cols: wrong size")
	}
	for i := 0; i < s.NumRow(); i++ {
		for k := 0; k < s.NumCol(); k++ {
			if s.Get(i, k) != m.Get(i, k+start) {
				t.Fatal("Cols: elements don't match")
			}
		}
	}

}

func TestDenseSetAll(t *testing.T) {
	m := MakeDense(5, 6)
	start, stop := 1, 4
	m.Cols(start, stop).SetAll(1)
	for i := 0; i < m.NumRow(); i++ {
		for k := 0; k < m.NumCol(); k++ {
			v := m.Get(i, k)
			if k >= start && k < stop && v == 1 {
				continue
			}
			if (k < start || k >= stop) && v == 0 {
				continue
			}
			m.wrongValue(t, i, k)
		}
	}
}

func TestDenseEqual(t *testing.T) {
	a := MakeDense(6, 5)
	a.Cols(1, 2).SetAll(1)
	a.Cols(2, 4).SetAll(2)
	a.Cols(4, 5).SetAll(3)
	a.Rows(2, 3).SetAll(4)
	a.Rows(3, 4).SetAll(5)

	if !a.Equal(a) {
		t.Fatal("a != a")
	}

	b := MakeDense(4, 3)
	c := a.Cols(1, 4).Rows(1, 5)
	b.Copy(c)
	if !b.Equal(c) {
		t.Fatal("b != c")
	}

	d := a.Cols(0, 3).Rows(0, 4)
	if b.Equal(d) {
		t.Fatal("b == d")
	}
}

func TestDenseMulBy(t *testing.T) {
	m := MakeDense(6, 5)
	m.SetAll(2)
	start, stop := 1, 4
	m.Rows(start, stop).SetAll(3)
	s := m.Cols(start, stop)
	s.MulBy(s)
	for i := 0; i < m.NumRow(); i++ {
		for k := 0; k < m.NumCol(); k++ {
			v := m.Get(i, k)
			if k >= start && k < stop {
				if i >= start && i < stop && v == 9 {
					continue
				}
				if (i < start || i >= stop) && v == 4 {
					continue
				}
			} else {
				if i >= start && i < stop && v == 3 {
					continue
				}
				if (i < start || i >= stop) && v == 2 {
					continue
				}
			}
			m.wrongValue(t, i, k)
		}
	}
}

func TestDenseMul(t *testing.T) {
	a := MakeDense(7, 6)
	a.setRandom(-1, 1)
	b := MakeDense(7, 6)
	b.setRandom(-1, 1)
	c := MakeDense(7, 6)
	d := MakeDense(7, 6)

	c.Copy(a)
	c.MulBy(b)
	d.Mul(a, b)
	d.checkEqual(t, c)
}

func TestDenseAddTo(t *testing.T) {
	a := MakeDense(4, 5,
		1, 2, 3, 4, 5,
		6, 7, 8, 9, 0,
		0, 9, 8, 7, 6,
		5, 4, 3, 2, 1,
	)
	b := MakeDense(4, 5,
		1, 0, 1, 0, 1,
		0, 1, 0, 1, 0,
		1, 1, 0, 0, 1,
		1, 0, 0, 1, 1,
	)
	c := MakeDense(4, 5)

	c.Copy(a)
	c.AddTo(b, 1)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if c.Get(i, k) != a.Get(i, k)+b.Get(i, k) {
				c.wrongValue(t, i, k)
			}
		}
	}
	c.Copy(a)
	c.AddTo(b, -1)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if c.Get(i, k) != a.Get(i, k)-b.Get(i, k) {
				c.wrongValue(t, i, k)
			}
		}
	}
	c.Copy(a)
	c.AddTo(b, 2)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if c.Get(i, k) != a.Get(i, k)+b.Get(i, k)*2 {
				c.wrongValue(t, i, k)
			}
		}
	}

	start, stop := 1, 4
	x := c.Cols(start, stop)
	y := b.Cols(start, stop)

	c.Copy(a)
	x.AddTo(y, 1)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if k < start || k >= stop {
				if c.Get(i, k) != a.Get(i, k) {
					c.wrongValue(t, i, k)
				}
			} else {
				if c.Get(i, k) != a.Get(i, k)+b.Get(i, k) {
					c.wrongValue(t, i, k)
				}
			}
		}
	}
	c.Copy(a)
	x.AddTo(y, -1)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if k < start || k >= stop {
				if c.Get(i, k) != a.Get(i, k) {
					c.wrongValue(t, i, k)
				}
			} else {
				if c.Get(i, k) != a.Get(i, k)-b.Get(i, k) {
					c.wrongValue(t, i, k)
				}
			}
		}
	}
	c.Copy(a)
	x.AddTo(y, 2)
	for i := 0; i < c.NumRow(); i++ {
		for k := 0; k < c.NumCol(); k++ {
			if k < start || k >= stop {
				if c.Get(i, k) != a.Get(i, k) {
					c.wrongValue(t, i, k)
				}
			} else {
				if c.Get(i, k) != a.Get(i, k)+b.Get(i, k)*2 {
					c.wrongValue(t, i, k)
				}
			}
		}
	}
}

func TestDenseAdd(t *testing.T) {
	a := MakeDense(7, 6)
	a.setRandom(-1, 1)
	b := MakeDense(7, 6)
	b.setRandom(-1, 1)
	c := MakeDense(7, 6)
	d := MakeDense(7, 6)

	c.Copy(a)
	c.AddTo(b, 1)
	d.Add(a, b, 1)
	d.checkEqual(t, c)

	c.Scale(a, 2)
	c.AddTo(b, 2)
	d.Add(a, b, 2)
	d.checkEqual(t, c)

	start, stop := 1, 5
	x := c.Cols(start, stop)
	y := b.Cols(start, stop)
	c = MakeDense(7, 4)
	d = MakeDense(7, 4)

	c.Copy(x)
	c.AddTo(y, 1)
	d.Add(x, y, 1)
	d.checkEqual(t, c)

	c.Scale(x, 2)
	c.AddTo(y, 2)
	d.Add(x, y, 2)
	d.checkEqual(t, c)
}

func TestDenseSub(t *testing.T) {
	a := MakeDense(7, 6)
	a.setRandom(-1, 1)
	b := MakeDense(7, 6)
	b.setRandom(-1, 1)
	c := MakeDense(7, 6)
	d := MakeDense(7, 6)

	c.Copy(a)
	c.AddTo(b, -1)
	d.Sub(a, b, 1)
	d.checkEqual(t, c)

	c.Scale(a, 2)
	c.AddTo(b, -2)
	d.Sub(a, b, 2)
	d.checkEqual(t, c)

	start, stop := 1, 5
	x := c.Cols(start, stop)
	y := b.Cols(start, stop)
	c = MakeDense(7, 4)
	d = MakeDense(7, 4)

	c.Copy(x)
	c.AddTo(y, -1)
	d.Sub(x, y, 1)
	d.checkEqual(t, c)

	c.Scale(x, 2)
	c.AddTo(y, -2)
	d.Sub(x, y, 2)
	d.checkEqual(t, c)

	fmt.Println(d)
}

// Utils

func (m Dense) setRandom(min, max float32) {
	delta := max - min
	for i := 0; i < m.NumRow(); i++ {
		for k := 0; k < m.NumCol(); k++ {
			m.Set(i, k, min+delta*rand.Float32())
		}
	}
}

func (m Dense) wrongValue(t *testing.T, i, k int) {
	t.Fatalf("element (%d,%d) has wrong value %g\n%s", i, k, m.Get(i, k), m)
}

func (m Dense) checkEqual(t *testing.T, a Dense) {
	if !m.Equal(a) {
		t.Fatalf("matrices not equal\n%s\n\n%s", m, a)
	}
}
