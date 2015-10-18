// +build noos

package runtime

// This file imports noos package into runtime.
// Additionally it can contain code that can not be included into noos package
// because of dependency loops.

import _ "runtime/noos"
