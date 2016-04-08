// +build linux

package fmt

import "os"

func init() {
	DefaultWriter = os.Stdout
}
