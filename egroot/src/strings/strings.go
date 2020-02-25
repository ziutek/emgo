package strings

// IndexByte returns the index of first c in s or -1 if there is no c in s.
func IndexByte(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

// HasPrefix tests whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		if s[i] != prefix[i] {
			return false
		}
	}
	return true
}

// TrimSpace returns a slice of the string s, with all leading and trailing
// white space removed. BUG: unicode whitespaces not supported.
func TrimSpace(s string) string {
	for len(s) != 0 {
		b := s[0]
		if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
			break
		}
		s = s[1:]
	}
	for len(s) != 0 {
		b := s[len(s)-1]
		if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
			break
		}
		s = s[:len(s)-1]
	}
	return s
}
