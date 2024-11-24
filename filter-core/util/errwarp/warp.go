package errwarp

import (
	"fmt"
	"runtime"
)

func Warp(warpMsg string, oriErr error) (warpErr error) {
	pc, file, line, ok := runtime.Caller(1)
	fn := runtime.FuncForPC(pc)
	switch {
	case ok && oriErr != nil:
		return fmt.Errorf("%s\n\t\tat %s:%d in %s\n\tCause by: %s", warpMsg, file, line, fn.Name(), oriErr.Error())
	case ok && oriErr == nil:
		return fmt.Errorf("%s\n\t\tat %s:%d in %s", warpMsg, file, line, fn.Name())
	case !ok && oriErr != nil:
		return fmt.Errorf("%s\n\t\tat unknown\n\tCause by: %s", warpMsg, oriErr.Error())
	case !ok && oriErr == nil:
		return fmt.Errorf("%s\n\t\tat unknown", warpMsg)
	}
	return
}
