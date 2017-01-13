package matrix32

// Scale performs: d = a * s
func (d Dense) Scale(a Dense, s float32) {
	d.checkDim(a)
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		k := d.numcol
		for k >= 2 {
			k--
			dr[k] = ar[k] * s
			k--
			dr[k] = ar[k] * s
		}
		if k != 0 {
			dr[0] = ar[0] * s
		}
	}
}
