package main

/*
#cgo CFLAGS: -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -framework AppKit 
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>
*/
import "C"

import (
	"fmt"
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

//export AcceptsFirstResponder
func AcceptsFirstResponder(self C.id, op C.SEL) C.BOOL {
	fmt.Println("acceptsFirstResponder")
	return 1
}

//export IInitWithFrame
func IInitWithFrame(self C.id, op C.SEL, aRect C.CGRect) C.id {
	rect := TypeNSRect(aRect)
	simpleView := (Object)(unsafe.Pointer(self))
	simpleView = simpleView.Class().Instance("alloc")
	simpleView = simpleView.CallSuperR("initWithFrame:", rect)
	return (C.id)(unsafe.Pointer(simpleView))
}

//export VDrawRect
func VDrawRect(self C.id, op C.SEL, aRect C.CGRect) {
	view := (Object)(unsafe.Pointer(self))
	rect := TypeNSRect(aRect)
	
//	C.NSRectFill(rect)
	
	fmt.Println("drawRect:", rect)
//	view.Call("lockFocus")
	view.CallSuperR("drawRect:", rect)
//	view.CallSuper("lockFocus")
	
/*
fixes:

get NSGraphicsContext directly
	lockFocus and unlockFocus - get apple docs for use

test nib by creating an objc app that uses the same nib	

super is broken, is callsuper doing anything?
	are methods passing through?									- test case
	is SimpleView actually instantiated properly as a subclass?		- test case

can we call a "delegate" method?

*/
	
	NSColor(WhiteColor).Call("set")
	ClassForName("NSBezierPath").InstanceR("fillRect:", rect)
	
	NSColor(GreenColor).Call("set")
	rect2 := MakeNSRect(5,5,50,50)
	
	ClassForName("NSBezierPath").InstanceR("fillRect:", rect2)
	
//	view.CallSuper("unlockFocus")
}

//export ZWindowResize
func ZWindowResize(self C.id, op C.SEL, notification C.id) {
	controller := (Object)(unsafe.Pointer(self))
	simpleView := controller.InstanceVariable("itsView")
	simpleView.Call("invalidateIntrinsicContentSize")
//	simpleView.CallSuper("setNeedsDisplay")
}

/*
* main()
* Main function for testing
 */
func main() {

	view := ClassForName("NSView").Subclass("SimpleView")
	view.AddMethod("acceptsFirstResponder", AcceptsFirstResponder)
	view.AddMethod("initWithFrame:", IInitWithFrame)
	view.AddMethod("drawRect:", VDrawRect)
	view.Register()

	controller := ClassForName("NSObject").Subclass("Controller")
	controller.AddMethod("windowDidResize:", ZWindowResize)
	controller.AddIvar("itsView", view)
	controller.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", NSString("."))
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("SimpleView"), dict, app.Call("zone"))
	
	
	
	// testing message passing
	/*
	windowsArray := app.Call("windows")
	windowsCount := (NSUInteger)(windowsArray.Call("count"))	
	bundle.I("initWithPath:", NSString("A String"), MakeNSRect(0,0,0,0), windowsArray, windowsCount, MakeNSBoolean(true))
	fmt.Println("windowsCount", windowsCount)
	
	windowsCount = (NSUInteger)(windowsArray.I("count"))
	fmt.Println("windowsCount", windowsCount)*/
	
	app.Call("run")

}
