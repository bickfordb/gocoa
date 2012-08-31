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
	"fmt"
	"gocoa"
	"unsafe"
)

func init() {
	gocoa.InitMac()
}

/*
* main()
* Main function for testing
 */
func main() {

	hellow := gocoa.ClassForName("NSObject").Subclass("ApplicationController")
	hellow.AddMethod("applicationWillFinishLaunching:", BApplicationWillFinishLaunching)
	hellow.AddMethod("buttonClick:", IButtonClick)
	hellow.AddIvar("textBox1", gocoa.ClassForName("NSTextField"))
	hellow.Register()

	app := gocoa.ClassForName("NSApplication").Instance("sharedApplication")
	bundle := gocoa.ClassForName("NSBundle").Instance("alloc")
	path := gocoa.NSString(".")
	dict := gocoa.NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", gocoa.NSString("HelloWorld"), dict, app.Call("zone"))
	
	
	icon := gocoa.ClassForName("NSImage").Instance("alloc")
	icon = icon.Call("initByReferencingFile:", gocoa.NSString("go.icns"))
	app.Call("setApplicationIconImage:", icon)
	
	
	app.Call("run")
}

//export BApplicationWillFinishLaunching
func BApplicationWillFinishLaunching(self C.id, op C.SEL, notification C.id) {
	fmt.Println("applicationWillFinishLaunching:")

	notify := gocoa.ObjectForId((uintptr)(unsafe.Pointer(notification)))
	application := notify.Call("object")

	windowsArray := application.Call("windows")
	windowsCount := (gocoa.NSUInteger)(windowsArray.Call("count").Pointer)
	var ix gocoa.NSUInteger
	for ix = 0; ix < windowsCount; ix++ {
		window := windowsArray.CallI("objectAtIndex:", ix)
		window.Call("setTitle:", gocoa.NSString("Form Loaded"))
	}

	me := gocoa.ObjectForId((uintptr)(unsafe.Pointer(self)))
	textBox1 := me.InstanceVariable("textBox1")
	textBox1.Call("setStringValue:", gocoa.NSString("Form Loaded"))
}

//export IButtonClick
func IButtonClick(self C.id, op C.SEL, sender C.id) {
	fmt.Println("buttonClick:")
	me := gocoa.ObjectForId((uintptr)(unsafe.Pointer(self)))
	textBox1 := me.InstanceVariable("textBox1")
	textBox1.Call("setStringValue:", gocoa.NSString("Button Pushed"))
}
