// This example writes "Hello world!" string to default debug port, in this
// case to ITM stimulus port 0.
//
// It uses default (reset) clock source (HSI, 16MHz) so use ./load-oocd.sh
// instead ../load-oocd.sh.
package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("Hello world!")
	}
}
