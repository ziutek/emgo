package matrix32

// AddTo performs: d += a * s
func (d *Dense) AddTo(a *Dense, s float32) {
	d.checkDims(a)
	switch s {
	case 1:
		for i := 0; i < d.rows; i++ {
			dr := d.v[i*d.stride:]
			ar := a.v[i*a.stride:]
			k := d.cols - 1
			for k > 0 {
				dr[k] += ar[k]
				k--
				dr[k] += ar[k]
				k--
			}
			if k == 0 {
				dr[0] += ar[0]
			}
		}
	case -1:
		for i := 0; i < d.rows; i++ {
			dr := d.v[i*d.stride:]
			ar := a.v[i*a.stride:]
			k := d.cols - 1
			for k > 0 {
				dr[k] -= ar[k]
				k--
				dr[k] -= ar[k]
				k--
			}
			if k == 0 {
				dr[0] -= ar[0]
			}
		}
	default:
		for i := 0; i < d.rows; i++ {
			dr := d.v[i*d.stride:]
			ar := a.v[i*a.stride:]
			k := d.cols - 1
			for k > 0 {
				dr[k] += ar[k] * s
				k--
				dr[k] += ar[k] * s
				k--
			}
			if k == 0 {
				dr[0] += ar[0] * s
			}
		}
	}
}
