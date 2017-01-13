package matrix32

// Mul performs: d = a * b
func (d Dense) Mul(a, b Dense) {
	d.checkDim(a)
	d.checkDim(b)
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		br := b.v[i*b.stride:]
		k := d.numcol
		for k >= 2 {
			k--
			dr[k] = ar[k] * br[k]
			k--
			dr[k] = ar[k] * br[k]
		}
		if k != 0 {
			dr[0] = ar[0] * br[0]
		}
	}
}
