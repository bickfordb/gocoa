package main

//#cgo LDFLAGS: -framework AppKit 
import "C"

import (
	. "gocoa"
)

/*
* main() A minimal nibless Cocoa application, bypasses Interface Builder
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
	
	quitMenuItem.Call("initWithTitle:action:keyEquivalent:", quitTitle, SelectorForName("terminate:"), NSString("q"))
	appMenu.Call("addItem:", quitMenuItem)
	appMenuItem.Call("setSubmenu:", appMenu)
	window := ClassForName("NSWindow").Instance("alloc")

	window.Call("initWithContentRect:styleMask:backing:defer:", MakeNSRect(0, 0, 200, 200), NSTitledWindowMask, NSBackingStoreBuffered, MakeNSBoolean(false))
	window.Call("cascadeTopLeftFromPoint:", NSPoint{20,20})
	window.Call("setTitle:", appName)
	window.Call("makeKeyAndOrderFront:", (Object)(0))
	
	app.Call("activateIgnoringOtherApps:", MakeNSBoolean(true))
	app.Call("run")

}
