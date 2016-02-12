// This example writes "Hello world!" string to default debug port, in this
// case to ITM stimulus port 0.
package main

import (
	"fmt"
)

func main() {
	for {
		fmt.Println("Hello world!")
	}
}
