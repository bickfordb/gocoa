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

 /*
//export BIsOpaque
func BIsOpaque(self C.id, op C.SEL) C.BOOL {
	fmt.Println("isOpaque")
	return (C.BOOL)(1)
}*/

//export IInitWithFrame
func IInitWithFrame(self C.id, op C.SEL, aRect C.CGRect) C.id {
	rect := TypeNSRect(aRect)
	simpleView := ObjectForId((uintptr)(unsafe.Pointer(self)))
	simpleView = simpleView.Class().Instance("alloc")
	simpleView = simpleView.CallSuperR("initWithFrame:", rect)
	
	return (C.id)(unsafe.Pointer(simpleView))
}


/*
Fixing this is going to involve msgSend being able to return arbitrary data
*/


//export VDrawRect
func VDrawRect(self C.id, op C.SEL, aRect C.CGRect) {
	rect := TypeNSRect(aRect)
	fmt.Println("drawRect:", rect.String())
		
//	view := ObjectForId((uintptr)(unsafe.Pointer(self)))
//	view.CallSuperR("drawRect:", rect)
	
//	view.ListInstanceVariables()
//	view.ListMethods()
//	view.Class().Super().ListMethods()

//	"To kill you must know your enemy. And in this case, my enemy, is a varmint."
	
//	rect2 := NSMakeRect(rect)
	
//	transform := ClassForName("NSAffineTransform").Instance("transform")
//	bezier := ClassForName("NSBezierPath").Instance("bezierPath")
	
	
// need window location I guess
	
	NSColor(GreenColor).Call("set")
	
	rect2 := NSMakeRect(5,5,50,50)
	fmt.Println("rect2:", rect2.String())
	
	ClassForName("NSBezierPath").InstanceR("fillRect:", rect2)
	
	
//	bezier.Call("fill")
	
	
	//.Call("init")
//	bezier.ListMethods()
//	bezier.Class().Super().ListMethods()
	
	
//	bounds := bezier.Call("bounds") 
//	fmt.Println("bounds:", bounds.String())
	
//	bezier.Call("fill") 
//	bezier.CallR("fillRect:",rect) 

}

//export ZWindowResize
func ZWindowResize(self C.id, op C.SEL, notification C.id) {
	fmt.Println("windowDidResize:")
	controller := ObjectForId((uintptr)(unsafe.Pointer(self)))
	simpleView := controller.InstanceVariable("itsView")
	simpleView.Call("invalidateIntrinsicContentSize")
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

	controller := ClassForName("NSObject").Subclass("Controller")
	controller.AddMethod("windowDidResize:", ZWindowResize)
	controller.AddIvar("itsView", view)
	controller.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	path := NSString(".")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("SimpleView"), dict, app.Call("zone"))
	
	
	// testing message passing
	
//	windowsArray := app.Call("windows")
//	windowsCount := (NSUInteger)(windowsArray.Call("count"))
	
//	bundle.I("initWithPath:", path, NSMakeRect(0,0,0,0), windowsCount, NSMakeBoolean(true))
	
	
	
	app.Call("run")

}
