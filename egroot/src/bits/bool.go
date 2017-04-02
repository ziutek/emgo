package bits

//c:inline
func one(b bool) int

// One returns 1 if b == true and 0 if b == false.
func One(b bool) int {
	return one(b)
}