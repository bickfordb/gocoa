package gocoa

/*
#cgo CFLAGS: -I/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/HIServices.framework/Versions/A/Headers/
#cgo LDFLAGS: -framework ApplicationServices 
#include <Processes.h>
*/
import "C"

import (
	"runtime"
)

/*
* init() 
* LockOSThread() is necessary to ensure that cocoa is being called from the main thread,
* Go may and does spawn main() as a secondary thread in the Mac implementation. TransformProcessType()
* on the current ProcessSerialNumber allows a unix process to be promoted to a full-fledged Mac app
* with a UI.
 */
func init() {
	runtime.LockOSThread()
	var psn C.ProcessSerialNumber
	C.GetCurrentProcess(&psn)
	C.TransformProcessType(&psn, 1)
	C.SetFrontProcess(&psn)
}
