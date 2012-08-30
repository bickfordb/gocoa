package gocoa

/*
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <objc/objc-runtime.h>
*/
//#cgo CFLAGS: -I/usr/include 
//#cgo LDFLAGS: -lobjc 
import "C"

import (
	"unsafe"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"math"
)


// ------------- functions for registering Go-defined classes

/*
// requires -I/usr/include/mach-o, <mach-o/dyld.h>
func loadThySelfMac(symbol string) *[0]byte {
	

	fmt.Println("symbol:", symbol)
	
	nssymbol := NSLookupSymbolInImage(nil,C.CString(symbol))
	fmt.Println("nssymbol:", nssymbol)
	if nssymbol == C.NULL {
		fmt.Println("error:", nssymbol)
	}
	
	//void* NSAddressOfSymbol(NSSymbol symbol);
	symbol_address := C.NSAddressOfSymbol(nssymbol)
	fmt.Println("symbol_address:", symbol_address)
	if symbol_address == C.NULL {
		fmt.Println("error:", symbol_address)
	}
	
	return (*[0]byte)(unsafe.Pointer(symbol_address))
}*/


func loadThySelf(symbol string) *[0]byte {
	
	fmt.Println("symbol:", symbol)
	
	this_process := C.dlopen(nil, C.RTLD_NOW )
//	fmt.Println("this_process:", this_process)
//	if this_process == C.NULL {
//		fmt.Println("error:", C.GoString(C.dlerror()))
//	}
	
	symbol_address := C.dlsym(this_process, C.CString(symbol))
	fmt.Println("\tsymbol_address:", symbol_address)
	if symbol_address == nil {
		fmt.Println("\terror:", C.GoString(C.dlerror()))
	}
	
	defer C.dlclose(this_process)
	return (*[0]byte)(unsafe.Pointer(symbol_address))
}


func (cls *Class) Subclass(subclassName string) *Class {
	class_id := C.objc_allocateClassPair((C.Class)(unsafe.Pointer(cls.id)), C.CString(subclassName), (C.size_t)(0) )
	return &Class{(uintptr)(unsafe.Pointer(class_id))}
}


func trimPackage(typePath string) string {

	if i := strings.LastIndex(typePath, "."); i != -1 {
		return typePath[i+1:]
	}
	return typePath
}


func (cls *Class) AddMethod(methodName string, implementor interface{}) bool {
	
	v := reflect.ValueOf(implementor)
//	fmt.Println("type:", v.Type())
	
	if (v.Kind() == reflect.Func) && (v.Type().NumIn() > 1) {
		
		types := "v"
		
		if v.Type().NumOut() == 1 {
			types = "@"
		}
		
		implementorName := trimPackage(runtime.FuncForPC(v.Pointer()).Name())
		
//		fmt.Println("name:", implementorName)
		numArgs := v.Type().NumIn()
		for i:=0; i < numArgs; i++ {
			argType := trimPackage(v.Type().In(i).String())
//			fmt.Println("\t", i, argType)
			types = types + objcArgTypeString(argType)
			
		}
//		fmt.Println("types:", types)
		
		sel := C.sel_registerName(C.CString(methodName))
		imp := loadThySelf(implementorName)
		
//		fmt.Println("imp:", imp)
		result := C.class_addMethod((C.Class)(unsafe.Pointer(cls.id)), sel, imp, C.CString(types))
		return (result == 1)
	}
	
//	fmt.Println("****implementor did not satisfy requirements****")
	return false;
}


func (cls *Class) AddIvar(ivarName string, ivarClass *Class) bool {
	
	types := objcArgTypeString(ivarClass.Name())
	size := (C.size_t)(unsafe.Sizeof(cls.id))
	alignment := (C.uint8_t)(math.Log2((float64)(unsafe.Sizeof(cls.id))))
	result := C.class_addIvar((C.Class)(unsafe.Pointer(cls.id)), C.CString(ivarName), size, alignment, C.CString(types))
	
	return (result == 1)
//BOOL class_addIvar(Class cls, const char *name, size_t size, uint8_t alignment, const char *types)
}

func (cls *Class) Register() {
	C.objc_registerClassPair((C.Class)(unsafe.Pointer(cls.id)))
}


// XXX somewhat incomplete 
func objcArgTypeString(argType string) string {
	
	switch(argType) {
		case "_Ctype_id":		return "@"
		case "_Ctype_SEL":		return ":"
		case "_Ctype_CGRect":	return "{_NSRect={_NSPoint=ff}{_NSSize=ff}}"
		default: return "@"
	}
	return ""
}

