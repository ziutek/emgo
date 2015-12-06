// build +noos

package setup

import (
	"rtos"
)

func sysclkChanged() {
	if rtos.MaxTasks() == 0 {
		return
	}
	// TODO: Setup there RTC0 (RTC1 ???).
}
