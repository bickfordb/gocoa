package main

/*
#cgo CFLAGS: -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -framework AppKit 
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>
*/
import "C"

import (
	. "gocoa"
	"unsafe"
)

/*
* The following is an example of something that I don't understand well and represents a problem:
* the exported symbol names seem to determine whether the funtions can be found by dlsym.
* If all are named similarly, only one is found. If they are dissimilar enough in their first
* character, they apparently can be loaded.
*
* I have a bug open that may be resolved with some coming linker changes, remains to be seen.
 */


//export IInitWithFrame
func IInitWithFrame(self C.id, op C.SEL, aRect C.CGRect) C.id {
	rect := TypeNSRect(aRect)
	simpleView := (Object)(unsafe.Pointer(self))
	simpleView = simpleView.Class().Instance("alloc")
	simpleView = simpleView.CallSuper("initWithFrame:", rect)
	return (C.id)(unsafe.Pointer(simpleView))
}

//export VDrawRect
func VDrawRect(self C.id, op C.SEL, aRect C.CGRect) {
rect := TypeNSRect(aRect)
	
	NSColor(WhiteColor).Call("set")
	ClassForName("NSBezierPath").Instance("fillRect:", rect)
	
	NSColor(GreenColor).Call("set")
	rect2 := MakeNSRect(5,5,50,50)
	
	ClassForName("NSBezierPath").Instance("fillRect:", rect2)
}


/*
* main()
* Main function for testing
 */
func main() {

	view := ClassForName("NSView").Subclass("SimpleView")
	view.AddMethod("initWithFrame:", IInitWithFrame)
	view.AddMethod("drawRect:", VDrawRect)
	view.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	dict := NSDictionary("NSOwner", app)
	
	bundle = bundle.Call("initWithPath:", NSString("."))
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("SimpleView"), dict, app.Call("zone"))
	
	app.Call("run")
}
