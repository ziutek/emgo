package uart

type circnt struct {
	idx, cnt uint32
}

func (c circnt) isrAdd(n, length uint32) {
	c.idx += n
	if c.idx >= length {
		c.idx -= length
		c.cnt++
	}
}
