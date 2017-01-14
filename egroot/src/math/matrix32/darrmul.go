package matrix32

// ArrMul performs: d = a * b
func (d Dense) ArrMul(a, b Dense) {
	d.checkDim(a)
	d.checkDim(b)
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		br := b.v[i*b.stride:]
		k, n := 0, d.numcol-1
		for k < n {
			dr[k+0] = ar[k+0] * br[k+0]
			dr[k+1] = ar[k+1] * br[k+1]
			k += 2
		}
		if k == n {
			dr[k] = ar[k] * br[k]
		}
	}
}
