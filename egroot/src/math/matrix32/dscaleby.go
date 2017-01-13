package matrix32

// ScaleBy performs: a *= s
func (d Dense) ScaleBy(s float32) {
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		k := d.numcol
		for k >= 2 {
			k--
			dr[k] *= s
			k--
			dr[k] *= s
		}
		if k != 0 {
			dr[0] *= s
		}
	}
}
