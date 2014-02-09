package matrix

type Dense struct {
	v          []float32 // [row, row, ..., row]
	rows, cols int
	stride     int // distance between vertically adjacent elements
}

// NewDense creates new matrix that refers to v
func NewDense(rows, cols, stride int, v []float32) (d Dense) {
	n := rows * stride
	if n > len(v) {
		panic("rows * stride > len(v)")
	}
	d.rows = rows
	d.cols = cols
	d.stride = stride
	d.v = v
	return
}

// Zero sets d to zero matrix
func (d *Dense) Zero() {
	for i := 0; i < len(d.v); i++ {
		d.v[i] = 0
	}
}

// Identity sets d to identity matrix (panics id d isnt a square matrix)
func (d *Dense) Identity() {
	if d.rows != d.cols {
		panic("attemt to create not square identity matrix")
	}
	for i := 0; i < len(d.v); i += d.stride + 1 {
		d.v[i] = 1
	}
}

// Size returns dimensions of the matrix (rows, cols)
func (d *Dense) Size() (int, int) {
	return d.rows, d.cols
}

// Rows returns number of rows
func (d *Dense) Rows() int {
	return d.rows
}

// Cols returns number of columns
func (d *Dense) Cols() int {
	return d.cols
}

// Stride returns distance between vertically adjacent elements
func (d *Dense) Stride() int {
	return d.stride
}

// Elems returns internal buffer of elements
func (d *Dense) Elems() []float32 {
	return d.v
}

// Get returns element from row i and column k
func (d *Dense) Get(i, k int) float32 {
	return d.v[i*d.stride+k]
}

// Set sets element in row i and column k
func (d *Dense) Set(i, k int, a float32) {
	d.v[i*d.stride+k] = a
}

// SetAll sets all elements to a
func (d *Dense) SetAll(a float32) {
	for i := 0; i < d.rows; i++ {
		row := d.v[i*d.stride:]
		k := d.cols
		for k >= 2 {
			k--
			row[k] = a
			k--
			row[k] = a
		}
		if k != 0 {
			row[0] = a
		}
	}
}

// Utils

func (d *Dense) checkEqualDims(a *Dense) {
	if d.Rows() != a.Rows() || d.Cols() != d.Cols() {
		panic("matrix dimensions not equal")
	}
}
