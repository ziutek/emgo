package fmt

func Fprint(w io.Writer, a ...Formatter) (n int, err error) {
	var m int
	for _, v := range a {
		m, err = v.Format(w)
		n += m
		if err != nil {
			return
		}
	}
	return
}
