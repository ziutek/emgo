package matrix32

// ScaleBy performs: a *= s
func (d Dense) ScaleBy(s float32) {
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		k, n := 0, d.numcol-1
		for k < n {
			dr[k+0] *= s
			dr[k+1] *= s
			k += 2
		}
		if k == n {
			dr[k] *= s
		}
	}
}
