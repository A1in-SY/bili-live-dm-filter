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
		return fmt.Errorf("%s\n\tat %s:%d in %s\nCause by: %s", warpMsg, file, line, fn.Name(), oriErr.Error())
	case ok && oriErr == nil:
		return fmt.Errorf("%s\n\tat %s:%d in %s", warpMsg, file, line, fn.Name())
	case !ok && oriErr != nil:
		return fmt.Errorf("%s\n\tat unknown\nCause by: %s", warpMsg, oriErr.Error())
	case !ok && oriErr == nil:
		return fmt.Errorf("%s\n\tat unknown", warpMsg)
	}
	return
}
