package matrix32

type Dense struct {
	v          []float32 // [row, row, ..., row]
	rows, cols int
	stride     int // distance between vertically adjacent elements
}

// NewDense creates new matrix that refers to v
func NewDense(rows, cols, stride int, v []float32) Dense {
	n := rows * stride
	if n > len(v) {
		panic("matrix32: rows * stride > len(v)")
	}
	return Dense{v: v, rows: rows, cols: cols, stride: stride}
}

// Zero sets d to zero matrix.
func (d *Dense) Zero() {
	v := d.v
	k := len(v) - 1
	for k > 0 {
		v[k] = 0
		v[k-1] = 0
		k -= 2
	}
	if k == 0 {
		v[0] = 0
	}
}

// Identity sets d to identity matrix (panics if d isn't a square matrix).
func (d *Dense) Identity() {
	if d.rows != d.cols {
		panic("matrix32: attempt to create not square identity matrix")
	}
	d.Zero()
	for i := 0; i < len(d.v); i += d.stride + 1 {
		d.v[i] = 1
	}
}

// Size returns dimensions of the matrix (rows, cols).
func (d *Dense) Size() (int, int) {
	return d.rows, d.cols
}

// Rows returns number of rows.
func (d *Dense) Rows() int {
	return d.rows
}

// Cols returns number of columns.
func (d *Dense) Cols() int {
	return d.cols
}

// Stride returns distance between vertically adjacent elements.
func (d *Dense) Stride() int {
	return d.stride
}

// Elems returns internal buffer of elements.
func (d *Dense) Elems() []float32 {
	return d.v
}

// Get returns element from row i and column k.
func (d *Dense) Get(i, k int) float32 {
	return d.v[i*d.stride+k]
}

// Set sets element in row i and column k.
func (d *Dense) Set(i, k int, a float32) {
	d.v[i*d.stride+k] = a
}

// SetAll sets all elements to a.
func (d *Dense) SetAll(a float32) {
	for i := 0; i < d.rows; i++ {
		row := d.v[i*d.stride:]
		k := d.cols - 1
		for k > 0 {
			row[k] = a
			row[k-1] = a
			k -= 2
		}
		if k == 0 {
			row[0] = a
		}
	}
}

// Hslice returns a slice of a matrix that contains rows from start to stop - 1.
func (d *Dense) Hslice(start, stop int) Dense {
	if start > stop || start < 0 || stop > d.rows {
		panic("matrix32: bad indexes for horizontal slice")
	}
	return Dense{
		v:      d.v[start*d.stride : stop*d.stride],
		rows:   stop - start,
		cols:   d.cols,
		stride: d.stride,
	}
}

// Vslice returns a slice of a matrix that contains cols from start to stop - 1
func (d *Dense) Vslice(start, stop int) Dense {
	if start > stop || start < 0 || stop > d.cols {
		panic("matrix32: bad indexes for vertical slice")
	}
	return Dense{
		v:      d.v[start : (d.rows-1)*d.stride+stop],
		rows:   d.rows,
		cols:   stop - start,
		stride: d.stride,
	}

}

// Hvec returns horizontal vector that refers to d. Panics if cols != stride.
func (d *Dense) Hvec() Dense {
	if d.cols != d.stride {
		panic("matrix32: can't convert matrix to horizontal vector: cols != stride")
	}
	return Dense{v: d.v, rows: 1, cols: len(d.v), stride: len(d.v)}
}

// Vvec returns vertical vector that refers to d. Panics if cols != stride.
func (d *Dense) Vvec() Dense {
	if d.cols != d.stride {
		panic("matrix32: can't convert matrix to vertical vector: cols != stride")
	}
	return Dense{v: d.v, rows: len(d.v), cols: 1, stride: 1}
}

// Equal returns true if matrices are equal
func (d *Dense) Equal(a *Dense) bool {
	if d.rows != a.rows || d.cols != a.cols {
		return false
	}
	for i := 0; i < d.rows; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		for k := 0; k < d.cols; k++ {
			if dr[k] != ar[k] {
				return false
			}
		}
	}
	return true
}

// Utils

func (d *Dense) checkDims(a *Dense) {
	if d.rows != a.rows || d.cols != a.cols {
		panic("matrix32: dimensions not equal")
	}
}
