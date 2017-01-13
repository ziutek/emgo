package matrix32

// MulBy performs: d *= a
func (d Dense) MulBy(a Dense) {
	d.checkDim(a)
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		k := d.numcol
		for k >= 2 {
			k--
			dr[k] *= ar[k]
			k--
			dr[k] *= ar[k]
		}
		if k != 0 {
			dr[0] *= ar[0]
		}
	}
}
