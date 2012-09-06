package main

/*
#cgo LDFLAGS: -framework AppKit
#include <objc/objc-runtime.h>
*/
import "C"

import (
	. "gocoa"
	"unsafe"
)


func main() {
	hellow := ClassForName("NSObject").Subclass("ApplicationController")
	hellow.AddMethod("applicationWillFinishLaunching:", BApplicationWillFinishLaunching)
	hellow.AddMethod("buttonClick:", IButtonClick)
	hellow.AddIvar("textBox1", ClassForName("NSTextField"))
	hellow.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	path := NSString(".")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("HelloWorld"), dict, app.Call("zone"))

	icon := ClassForName("NSImage").Instance("alloc")
	icon = icon.Call("initByReferencingFile:", NSString("go.icns"))
	app.Call("setApplicationIconImage:", icon)

	app.Call("run")
}

//export BApplicationWillFinishLaunching
func BApplicationWillFinishLaunching(self C.id, op C.SEL, notification C.id) {
	notify := (Object)(unsafe.Pointer(notification))
	application := notify.Call("object")

	windowsArray := application.Call("windows")
	windowsCount := (NSUInteger)(windowsArray.Call("count"))
	var ix NSUInteger
	for ix = 0; ix < windowsCount; ix++ {
		window := windowsArray.Call("objectAtIndex:", ix)
		window.Call("setTitle:", NSString("Form Loaded"))
	}

	me := (Object)(unsafe.Pointer(self))
	textBox1 := me.InstanceVariable("textBox1")
	textBox1.Call("setStringValue:", NSString("Form Loaded"))
}

//export IButtonClick
func IButtonClick(self C.id, op C.SEL, sender C.id) {
	me := (Object)(unsafe.Pointer(self))
	textBox1 := me.InstanceVariable("textBox1")
	textBox1.Call("setStringValue:", NSString("Button Pushed"))
}
