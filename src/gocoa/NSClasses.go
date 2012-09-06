package gocoa

import (
	"C"
	"unsafe"
)

func NSDictionary(key string, value Object) Object {
	return ClassForName("NSDictionary").Instance("alloc").Call("initWithObject:forKey:", value, NSString(key))
}

func NSString(inString string) Object {
	return ClassForName("NSString").Instance("stringWithUTF8String:", (charptr)(unsafe.Pointer(C.CString(inString))))
}

func NSStringToString(inString Object) string {
	return C.GoString((*C.char)(unsafe.Pointer(inString.Call("UTF8String"))))
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

func NSColor(color string) Object {
	return ClassForName("NSColor").Instance(color)
}


const (
	NSBorderlessWindowMask		NSUInteger = 0
    NSTitledWindowMask			NSUInteger = 1 << 0
    NSClosableWindowMask		NSUInteger = 1 << 1
    NSMiniaturizableWindowMask	NSUInteger = 1 << 2
    NSResizableWindowMask		NSUInteger = 1 << 3
)

const (
    NSBackingStoreRetained	 	NSUInteger = 0
    NSBackingStoreNonretained	NSUInteger = 1
    NSBackingStoreBuffered		NSUInteger = 2
)