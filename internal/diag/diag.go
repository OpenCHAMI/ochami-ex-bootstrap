package diag

import (
	"fmt"
	"os"
)

// Debug enables extra logging when true.
var Debug bool

// Logf writes formatted debug logs to stderr when Debug is true.
func Logf(format string, args ...any) {
	if !Debug {
		return
	}
	fmt.Fprintf(os.Stderr, "[DEBUG] "+format+"\n", args...)
}
