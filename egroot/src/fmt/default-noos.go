// +build noos

package fmt

import "rtos"

func init() {
	DefaultWriter = rtos.Debug(0)
}
