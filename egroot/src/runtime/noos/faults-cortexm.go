// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

//emgo:noinline
func nmiHandler() {
	for {
	}
}

//emgo:export
func faultHandler()
