package gocoa

import (
	"C"
	"unsafe"
)

func NSDictionary(key string, value *Object) *Object {
	return ClassForName("NSDictionary").Instance("alloc").Call("initWithObject:forKey:", value, NSString(key))
}

func NSString(toNSString string) *Object {
	cStringPtr := (uintptr)(unsafe.Pointer(C.CString(toNSString)))
	return ClassForName("NSString").Instance("stringWithUTF8String:", cStringPtr)
}

func GoString(nsString *Object) string {
	return C.GoString((*C.char)(unsafe.Pointer(nsString.Call("UTF8String").Pointer)))
}

const (
	BlackColor     = "blackColor"
	BlueColor      = "blueColor"
	BrownColor     = "brownColor"
	ClearColor     = "clearColor"
	CyanColor      = "cyanColor"
	DarkGrayColor  = "darkGrayColor"
	GrayColor      = "grayColor"
	GreenColor     = "greenColor"
	LightGrayColor = "lightGrayColor"
	MagentaColor   = "magentaColor"
	OrangeColor    = "orangeColor"
	PurpleColor    = "purpleColor"
	RedColor       = "redColor"
	WhiteColor     = "whiteColor"
	YellowColor    = "yellowColor"
)

func NSColor(color string) *Object {
	return ClassForName("NSColor").Instance(color)
}
