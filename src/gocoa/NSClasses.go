package gocoa

import (
	"C"
	"unsafe"
)

// NSDictionary requires alloc/init
func NSDictionary(key string, value *Object) *Object {
	dict := ClassForName("NSDictionary").Instance("alloc").Call("initWithObject:forKey:", value.id, NSString(key).id)
	return &Object{dict.id}
}

// NSString class method stringWithUTFString is helpful
func NSString(toNSString string) *Object {
	cStringPtr := (uintptr)(unsafe.Pointer(C.CString(toNSString)))
	class := ClassForName("NSString").Instance("stringWithUTF8String:", cStringPtr )
	return &Object{class.id}
}

func GoString(nsString *Object) string {
	return C.GoString((*C.char)(unsafe.Pointer(nsString.Call("UTF8String").id)))
}

const (
	blackColor = "blackColor"
	blueColor = "blueColor"
	/*brownColor
	clearColor
	cyanColor
	darkGrayColor
	grayColor
	greenColor
	lightGrayColor
	magentaColor
	orangeColor
	purpleColor*/
	RedColor = "redColor"
	whiteColor = "whiteColor"
	yellowColor = "yellowColor"
)


func NSColor(color string) *Object {
	class := ClassForName("NSColor").Instance(color)
	return &Object{class.id}
}


