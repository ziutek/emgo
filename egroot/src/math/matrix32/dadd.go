package matrix32

// Add performs: d = (a + b) * s
func (d Dense) Add(a, b Dense, s float32) {
	d.checkDim(a)
	d.checkDim(b)
	switch s {
	case 1:
		for i := 0; i < d.numrow; i++ {
			dr := d.v[i*d.stride:]
			ar := a.v[i*a.stride:]
			br := b.v[i*b.stride:]
			k, n := 0, d.numcol-1
			for k < n {
				dr[k+0] = ar[k+0] + br[k+0]
				dr[k+1] = ar[k+1] + br[k+1]
				k += 2
			}
			if k == n {
				dr[k] = ar[k] + br[k]
			}
		}
	default:
		for i := 0; i < d.numrow; i++ {
			dr := d.v[i*d.stride:]
			ar := a.v[i*a.stride:]
			br := b.v[i*b.stride:]
			k, n := 0, d.numcol-1
			for k < n {
				dr[k+0] = (ar[k+0] + br[k+0]) * s
				dr[k+1] = (ar[k+1] + br[k+1]) * s
				k += 2
			}
			if k == n {
				dr[k] = (ar[k] + br[k]) * s
			}
		}
	}
}
