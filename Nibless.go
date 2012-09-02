package main

//#cgo LDFLAGS: -framework AppKit 
import "C"

import (
	. "gocoa"
)

/*
* main()
* A minimal nibless Cocoa application, bypasses Interface Builder
 */
func main() {
	
	app := ClassForName("NSApplication").Instance("sharedApplication")
	menubar := ClassForName("NSMenu").Instance("new") 
	appMenuItem := ClassForName("NSMenuItem").Instance("new")
	menubar.Call("addItem:", appMenuItem)
	app.Call("setMainMenu:", menubar)
	appMenu := ClassForName("NSMenu").Instance("new")
	appName := ClassForName("NSProcessInfo").Instance("processInfo").Call("processName")
	quitTitle := NSString("Quit ").Call("stringByAppendingString:", appName)
	quitMenuItem := ClassForName("NSMenuItem").Instance("alloc")
//	third argument wants: app.Selector("terminate:")
	quitMenuItem.Call("initWithTitle:action:keyEquivalent:", quitTitle, (Object)(0), NSString("q"))
	appMenu.Call("addItem:", quitMenuItem)
	appMenuItem.Call("setSubmenu:", appMenu)
	window := ClassForName("NSWindow").Instance("alloc")
	
	window.Call("init")
//	wants: window.Call("initWithContentRect:styleMask:backing:defer:", NSMakeRect(0, 0, 200, 200), 
//		NSTitledWindowMask, NSBackingStoreBuffered, NSBoolean(false)).Call("autorelease")
//	wants: window.Call("cascadeTopLeftFromPoint:", NSPoint{20,20})
	window.Call("setTitle:", appName)
	window.Call("makeKeyAndOrderFront:", (Object)(0))
//	wants: app.Call("activateIgnoringOtherApps:", NSBoolean(true))
	app.Call("run")

}
