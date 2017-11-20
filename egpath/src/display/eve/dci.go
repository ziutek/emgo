package eve

type DCI interface {
	Begin() // Begin sets CSN pin to low.
	End() // End sets CSN pin to high
	
	Write([]byte) 
}