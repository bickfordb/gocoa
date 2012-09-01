package main

/*
#cgo LDFLAGS: -framework AppKit 
#include <objc/objc-runtime.h>
*/
import "C"

import (
	"fmt"
	. "gocoa"
	"unsafe"
)

/*
* main()
* a simple application with a datagrid in a scrollview
 */
func main() {

	hellow := ClassForName("NSObject").Subclass("ApplicationController")
	hellow.AddMethod("applicationWillFinishLaunching:", BApplicationWillFinishLaunching)
	hellow.AddIvar("scrollTable1", ClassForName("NSScrollView"))
	hellow.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	path := NSString(".")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("Application"), dict, app.Call("zone"))

	app.Call("run")
}

//export BApplicationWillFinishLaunching
func BApplicationWillFinishLaunching(self C.id, op C.SEL, notification C.id) {
	fmt.Println("applicationWillFinishLaunching:")

	notify := ObjectForId((uintptr)(unsafe.Pointer(notification)))
	application := notify.Call("object")

	windowsArray := application.Call("windows")
	windowsCount := (NSUInteger)(windowsArray.Call("count").Pointer)
	var ix NSUInteger
	for ix = 0; ix < windowsCount; ix++ {
		window := windowsArray.CallI("objectAtIndex:", ix)
		window.Call("setTitle:", NSString("Form Loaded"))
	}

	me := ObjectForId((uintptr)(unsafe.Pointer(self)))
	scrollTable1 := me.InstanceVariable("scrollTable1")
	fmt.Println("scrollTable1 class:", scrollTable1.Class().Name())
	//	textBox1.Call("setStringValue:", NSString("Form Loaded"))
}
