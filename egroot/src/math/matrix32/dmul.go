package matrix32

// Mul performs: d = a × b
func (d Dense) Mul(a, b Dense) {
	if d.numrow != a.numrow || d.numcol != b.numcol || a.numcol != b.numrow {
		panic("matrix32: MP: bad dimensions")
	}
	for i := 0; i < d.numrow; i++ {
		dr := d.v[i*d.stride:]
		ar := a.v[i*a.stride:]
		for k := 0; k < d.numcol; k++ {
			var p float32
			j, l := 0, k
			n := a.numcol - 1
			for j < n {
				p += ar[j+0] * b.v[l+0]
				p += ar[j+1] * b.v[l+b.stride]
				j += 2
				l += b.stride * 2
			}
			if j == n {
				p += ar[j] * b.v[l]
			}
			dr[k] = p
		}
	}
}
