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

	appctrl := ClassForName("NSObject").Subclass("ApplicationController")
	appctrl.AddMethod("applicationWillFinishLaunching:", BApplicationWillFinishLaunching)
	appctrl.AddIvar("scrollTable1", ClassForName("NSScrollView"))
	appctrl.Register()

	app := ClassForName("NSApplication").Instance("sharedApplication")
	bundle := ClassForName("NSBundle").Instance("alloc")
	path := NSString(".")
	dict := NSDictionary("NSOwner", app)

	bundle = bundle.Call("initWithPath:", path)
	bundle.Call("loadNibFile:externalNameTable:withZone:", NSString("TableView"), dict, app.Call("zone"))

	app.Call("run")
}

//export BApplicationWillFinishLaunching
func BApplicationWillFinishLaunching(self C.id, op C.SEL, notification C.id) {
	fmt.Println("applicationWillFinishLaunching:")

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
	scrollTable1 := me.InstanceVariable("scrollTable1")
	fmt.Println("scrollTable1 class:", scrollTable1.Class().Name())
	tableView := scrollTable1.Call("documentView")
	
	
	// add some stuff to the TableView
	// XXX, tableview discovering only the first column
	
	arrayController := ClassForName("NSArrayController").Instance("alloc").Call("init")
	dict := ClassForName("NSMutableDictionary").Instance("alloc").Call("init")
	
	objects := ClassForName("NSMutableArray").Instance("alloc").Call("init")
	objects.Call("addObject:", NSString("Joe"))
	objects.Call("addObject:", NSString("(444) 444-4444"))
	
	keys := ClassForName("NSMutableArray").Instance("alloc").Call("init")
	keys.Call("addObject:", NSString("column1"))
	keys.Call("addObject:", NSString("column2"))
	
	dict.Call("addObjects:forKeys:", objects, keys)
	
	arrayController.Call("addObject:", dict)
	
		
	column1 := tableView.Call("tableColumnWithIdentifier:", NSString("column1"))
	column1.Call("bind:toObject:withKeyPath:options:", NSString("value"), arrayController, NSString("arrangedObjects.column1"), (Object)(0))
	
	column2 := tableView.Call("tableColumnWithIdentifier:", NSString("column2"))
	column2.Call("bind:toObject:withKeyPath:options:", NSString("value"), arrayController, NSString("arrangedObjects.column2"), (Object)(0))
		
//	headerView := tableView.Call("headerView")
//	headerView.Class().ListMethods()
}
