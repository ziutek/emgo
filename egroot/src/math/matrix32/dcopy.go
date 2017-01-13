package matrix32

// Copy performs: d = a
func (d Dense) Copy(a Dense) {
	d.checkDim(a)
	for i := 0; i < d.numrow; i++ {
		ao := i * a.stride
		copy(d.v[i*d.stride:], a.v[ao:ao+d.numcol])
	}
}