// +build linux

package runtime

// This file imports noos package into runtime.
// Additionally it can contain code that can not be included into linux package
// because of dependency loops.

import _ "runtime/linux"
