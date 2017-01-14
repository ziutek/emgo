package matrix32

// ArrMulBy performs: d *= a
func (d Dense) ArrMulBy(a Dense) {
	d.checkDim(a)
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		k, n := 0, d.numcol-1
		for k < n {
			dr[k+0] *= ar[k+0]
			dr[k+1] *= ar[k+1]
			k += 2
		}
		if k == n {
			dr[k] *= ar[k]
		}
	}
}
