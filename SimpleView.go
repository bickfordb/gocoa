package main

/*
#include <stdlib.h>
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>
*/
//#cgo CFLAGS: -I/usr/include -I/System/Library/Frameworks/Foundation.framework/Versions/C/Headers/ -I/System/Library/Frameworks/AppKit.framework/Versions/C/Headers/ -I/System/Library/Frameworks/ApplicationServices.framework/Versions/A/Frameworks/HIServices.framework/Versions/A/Headers/ -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
//#cgo LDFLAGS: -lobjc -framework Foundation -framework AppKit -framework ApplicationServices -framework CoreGraphics
import "C"

import (
	"gocoa"
	"fmt"
	"unsafe"
	"encoding/binary"
	"bytes"
)

/*
* The following is an example of something that I don't understand well and represents a problem:
* the exported symbol names seem to determine whether the funtions can be found by dlsym.
* If all are named similarly, only one is found. If they are dissimilar enough in their first
* character, they apparently can be loaded.
*
* Possibly, there's some asynchronous operation that requires the time to iterate over the
* symbol table between calls, and maybe changing the name exposes that behavior.
*		solution: calling dlopen in a blocking goroutine with a wait() call?
*
* One possible solution is trying not to repeatedly call dlopen/dlclose
*/


//export IInitWithFrame
func IInitWithFrame(self C.id, op C.SEL, aRect C.CGRect) C.id {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, aRect)
	
	simpleView := gocoa.NewObject((uintptr)(unsafe.Pointer(self)))
	simpleView = simpleView.Class().Instance("alloc")
	simpleView = simpleView.CallSuperR("initWithFrame:", buf.Bytes())	
	return (C.id)(unsafe.Pointer(simpleView.Id()))
}


//export VDrawRect
func VDrawRect(self C.id, op C.SEL, aRect C.CGRect) {
	fmt.Println("drawRect:")
	
	view := gocoa.NewObject((uintptr)(unsafe.Pointer(self)))
	view.ListInstanceVariables()

	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, aRect)
	fmt.Println("len(aRectBytes)", len(buf.Bytes()))
	fmt.Println("bytes:", buf.Bytes())
	
	gocoa.NSColor(gocoa.RedColor).Call("set")
	
	bezier := gocoa.ClassForName("NSBezierPath").Instance("init")
//	bezier.ListMethods()
	bezier.CallR("fillRect:", buf.Bytes())		// terrible here
	
}


//export ZWindowResize
func ZWindowResize(self C.id, op C.SEL, notification C.id) {
	fmt.Println("windowDidResize:")
	controller := gocoa.NewObject((uintptr)(unsafe.Pointer(self)))
	simpleView := controller.InstanceVariable("itsView")
	simpleView.Call("invalidateIntrinsicContentSize")
}



func init() {
	gocoa.InitMac()
}


/*
* main()
* Main function for testing
*/
func main() {
		
	view := gocoa.ClassForName("NSView").Subclass("SimpleView")
	view.AddMethod("initWithFrame:", IInitWithFrame)
	view.AddMethod("drawRect:", VDrawRect)
	view.Register()
	
	controller := gocoa.ClassForName("NSObject").Subclass("Controller")
	controller.AddMethod("windowDidResize:", ZWindowResize)
	controller.AddIvar("itsView", view)
	controller.Register()
	
	
	app := gocoa.ClassForName("NSApplication").Instance("sharedApplication")
	bundle := gocoa.ClassForName("NSBundle").Instance("alloc")
	path := gocoa.NSString(".")
	dict := gocoa.NSDictionary("NSOwner", app)
	
	bundle = bundle.Call("initWithPath:", path.Id())
	bundle.Call("loadNibFile:externalNameTable:withZone:", gocoa.NSString("SimpleView").Id(), dict.Id(), app.Call("zone").Id())
		
	app.Call("run")
	
}


