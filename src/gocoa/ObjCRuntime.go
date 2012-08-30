package gocoa

/*
#cgo CFLAGS: -I/usr/include -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -lobjc -framework CoreGraphics
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>

static inline id gocoa_objc_msgSendI(id self, SEL op, long message) {
	return objc_msgSend(self, op, message);
}

static inline id gocoa_objc_msgSendR(id self, SEL op, CGRect message) {
	return objc_msgSend(self, op, message);
}

static inline id gocoa_objc_msgSendSuperR(struct objc_super *super, SEL op, CGRect message) {
	return objc_msgSendSuper(super, op, message);
}

static inline id gocoa_objc_msgSend(id self, SEL op, id messages[], int messagesLen) {
	switch (messagesLen) {
		case 1: return objc_msgSend(self, op, messages[0]);
		case 2: return objc_msgSend(self, op, messages[0], messages[1]);
		case 3: return objc_msgSend(self, op, messages[0], messages[1], messages[2]);
		case 4: return objc_msgSend(self, op, messages[0], messages[1], messages[2], messages[3]);
		case 5: return objc_msgSend(self, op, messages[0], messages[1], messages[2], messages[3], messages[4]);
		default: return objc_msgSend(self, op);
	}
}

static inline id gocoa_objc_msgSendSuper(struct objc_super *super, SEL op, id messages[], int messagesLen) {
	switch (messagesLen) {
		case 1: return objc_msgSendSuper(super, op, messages[0]);
		case 2: return objc_msgSendSuper(super, op, messages[0], messages[1]);
		case 3: return objc_msgSendSuper(super, op, messages[0], messages[1], messages[2]);
		case 4: return objc_msgSendSuper(super, op, messages[0], messages[1], messages[2], messages[3]);
		case 5: return objc_msgSendSuper(super, op, messages[0], messages[1], messages[2], messages[3], messages[4]);
		default: return objc_msgSendSuper(super, op);
	}
}
*/
import "C"

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"unsafe"
	
	"math"
	"reflect"
	"runtime"
	"strings"
)

/* class implementation *************************************************************** */

type Class struct{
	Pointer uintptr
}


func (cls *Class) classPointer() C.Class {
	return (C.Class)(unsafe.Pointer(cls.Pointer))
}


func (cls *Class) idPointer() C.id {
	return (C.id)(unsafe.Pointer(cls.Pointer))
}


func (cls *Class) respondsTo(selector string) bool {
	sel := C.sel_registerName(C.CString(selector))
	return (C.class_respondsToSelector(cls.classPointer(), sel) == 1)
}


func (cls *Class) Name() string {
	return C.GoString(C.class_getName(cls.classPointer()))
}


func (cls *Class) Superclass() *Class {
	return &Class{(uintptr)(unsafe.Pointer(C.class_getSuperclass(cls.classPointer())))}
}


func (cls *Class) Instance(calling string, args ...uintptr) *Object {
	return &Object{msgSend(&Object{cls.Pointer}, calling, args...)}
}


func (cls *Class) Method(name string) *Method {
	sel := C.sel_registerName(C.CString(name))
	return &Method{(uintptr)(unsafe.Pointer(C.class_getClassMethod(cls.classPointer(), sel)))}
}


func (cls *Class) ListInstanceVariables() {

	fmt.Println(cls.Name(), ": ivars")

	var outCount C.uint
	var ivarPointers []C.Ivar

	p := (C.class_copyIvarList(cls.classPointer(), &outCount))

	if p != nil {
		ivarPointers = (*[1 << 30]C.Ivar)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := Ivar{(uintptr)(unsafe.Pointer(ivarPointers[i]))}
			fmt.Println("\t", tmp.Name())
		}
		C.free(unsafe.Pointer(p))
	}
}

func (cls *Class) ListMethods() {

	fmt.Println(cls.Name(), ": methods")

	var outCount C.uint
	var methodPointers []C.Method

	p := (C.class_copyMethodList(cls.classPointer(), &outCount))

	if p != nil {
		methodPointers = (*[1 << 30]C.Method)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := Method{(uintptr)(unsafe.Pointer(methodPointers[i]))}
			fmt.Println("\t", tmp.Name())
		}
		C.free(unsafe.Pointer(p))
	}
}

/* object implementation *************************************************************** */

type Object struct {
	Pointer uintptr
}

// XXX name collision in C.types, can't use C.Object, need to create a namespace
// however, C.types seem to coerce
func (obj *Object) idPointer() C.id {
	return (C.id)(unsafe.Pointer(obj.Pointer))
}
/*
func (obj *Object) Id() uintptr {
	return obj.Pointer
}

func NewObject(object_id uintptr) *Object {
	return &Object{object_id}
}

func (obj *Object) AsClass() *Class {
	return &Class{obj.Pointer}
}*/

func (obj *Object) getMethod(name string) *Method {
	sel := C.sel_registerName(C.CString(name))
	method_id := (uintptr)(unsafe.Pointer(C.class_getInstanceMethod(obj.Class().classPointer(), sel)))
	return &Method{method_id}
}

func (obj *Object) Class() *Class {
	object_id := (uintptr)(unsafe.Pointer(C.object_getClass(obj.idPointer())))
	return &Class{object_id}
}

// XXX should take arguments 
func (obj *Object) Init() *Object {
	sel := C.sel_registerName(C.CString("init"))
	return &Object{(uintptr)(unsafe.Pointer(C.objc_msgSend(obj.idPointer(), sel)))}
}

func (obj *Object) CallR(method string, arg []byte) *Object {
	var copiedArg C.CGRect
	buf := bytes.NewBuffer(arg)
	/*	err := */ binary.Read(buf, binary.LittleEndian, copiedArg)
	return &Object{msgSendR(obj, method, copiedArg)}
}

func (obj *Object) CallI(method string, arg NSUInteger) *Object {
	return &Object{msgSendI(obj, method, arg)}
}

func (obj *Object) Call(method string, args ...*Object) *Object {
//	var outArgs [len(args)]uintptr
	outArgs := make([]uintptr, len(args))
	for i := 0; i < len(args); i++ {
		outArgs[i] = args[i].Pointer
	}
	return &Object{msgSend(obj, method, outArgs...)}
}

func (obj *Object) CallSuper(method string, args ...*Object) *Object {
	outArgs := make([]uintptr, len(args))
	for i := 0; i < len(args); i++ {
		outArgs[i] = args[i].Pointer
	}
	return &Object{msgSendSuper(obj, method, outArgs...)}
}

func (obj *Object) CallSuperR(method string, arg []byte) *Object {
	var copiedArg C.CGRect
	buf := bytes.NewBuffer(arg)
	/*	err := */ binary.Read(buf, binary.LittleEndian, copiedArg)
	return &Object{msgSendSuperR(obj, method, copiedArg)}
}

// XXX this does not necessarily return an object pointer
func (obj *Object) InstanceVariable(name string) *Object {
	var val uintptr
	ivar := C.object_getInstanceVariable(obj.idPointer(), C.CString(name), (*unsafe.Pointer)(unsafe.Pointer(&val)))
	typeenc := C.GoString(C.ivar_getTypeEncoding(ivar))
	
	if typeenc == "@" {
		return &Object{val}
	}
	return nil
}

func (obj *Object) SetInstanceVariable(name string, val *Object) {
	C.object_setInstanceVariable(obj.idPointer(), C.CString(name), unsafe.Pointer(val.Pointer))
}

func (obj *Object) ListInstanceVariables() {
	obj.Class().ListInstanceVariables()
}

func (obj *Object) ListMethods() {
	obj.Class().ListMethods()
}

/* class creation methods ------------------------------------------------------------- */

func (cls *Class) Subclass(subclassName string) *Class {
	class_id := C.objc_allocateClassPair((C.Class)(unsafe.Pointer(cls.Pointer)), C.CString(subclassName), (C.size_t)(0))
	return &Class{(uintptr)(unsafe.Pointer(class_id))}
}

func (cls *Class) AddMethod(methodName string, implementor interface{}) bool {

	v := reflect.ValueOf(implementor)

	if (v.Kind() == reflect.Func) && (v.Type().NumIn() > 1) {

		types := "v"
		impName := trimPackage(runtime.FuncForPC(v.Pointer()).Name())
		numArgs := v.Type().NumIn()
		
		if v.Type().NumOut() == 1 {
			types = "@"
		}
		
		for i := 0; i < numArgs; i++ {
			argType := trimPackage(v.Type().In(i).String())
			types = types + objcArgTypeString(argType)

		}

		sel := C.sel_registerName(C.CString(methodName))
		imp := loadThySelf(impName)
		result := C.class_addMethod(cls.classPointer(), sel, imp, C.CString(types))
		
		return (result == 1)
	}
	
	return false
}

func (cls *Class) AddIvar(ivarName string, ivarClass *Class) bool {

	types := objcArgTypeString(ivarClass.Name())
	size := (C.size_t)(unsafe.Sizeof(cls.Pointer))
	alignment := (C.uint8_t)(math.Log2((float64)(unsafe.Sizeof(cls.Pointer))))
	result := C.class_addIvar(cls.classPointer(), C.CString(ivarName), size, alignment, C.CString(types))

	return (result == 1)
}

func (cls *Class) Register() {
	C.objc_registerClassPair(cls.classPointer())
}



/* method implementation ************************************************************** */

type Method struct{
	Pointer uintptr
}

func (mthd *Method) methodPointer() C.Method {
	return (C.Method)(unsafe.Pointer(mthd.Pointer))
}

func (mthd *Method) ArgumentCount() int {
	return (int)(C.method_getNumberOfArguments(mthd.methodPointer()))
}

func (mthd *Method) ArgumentType(index int) string {
	var dst_len C.size_t
	var dst *C.char
	C.method_getArgumentType(mthd.methodPointer(), (C.uint)(index), dst, dst_len)
	return C.GoString(dst)
}

func (mthd *Method) Name() string {
	return C.GoString(C.sel_getName(C.method_getName(mthd.methodPointer())))
}

/* ivar implementation ************************************************************** */

type Ivar struct{ id uintptr }

func (ivr *Ivar) Name() string {
	return C.GoString(C.ivar_getName((C.Ivar)(unsafe.Pointer(ivr.id))))
}


/* utility methods ******************************************************************* */


func ClassForName(name string) *Class {
	class_id := (uintptr)(unsafe.Pointer(C.objc_getClass(C.CString(name))))
	return &Class{class_id}
}


func NewObject(object_id uintptr) *Object {
	return &Object{object_id}
}


// clumsy hacks abound
func msgSendI(receiver *Object, selector string, number NSUInteger) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendI(receiver.idPointer(), sel, (C.long)(number))))
}

func msgSendR(receiver *Object, selector string, rect C.CGRect) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendR(receiver.idPointer(), sel, rect)))
}

type superStruct struct {
	receiver uintptr
	class    uintptr
}

func msgSendSuperR(receiver *Object, selector string, rect C.CGRect) uintptr {
	var super superStruct
	super.receiver = receiver.Pointer
	super.class = receiver.Class().Superclass().Pointer
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendSuperR((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, rect)))
}

/*
* msgSend()
* Notice that you have to pass a pointer to the first array element to match the c array calling convention.
* The gocoa_ bridge function essentially, clumsily replicates Apple's deprecated objc_msgSendv().
 */
func msgSend(receiver *Object, selector string, args ...uintptr) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	if len(args) > 0 { // due to cgo calling convention, can't pass an empty array
		return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSend(receiver.idPointer(), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)))))
	}
	return (uintptr)(unsafe.Pointer(C.objc_msgSend(receiver.idPointer(), sel)))
}

/*
* msgSendSuper()
* The distinction here involves messaging to superclasses, with message receipt to the subclass. This
* requires initing a structure that refers to both the receiver and its superclass. 
 */
func msgSendSuper(receiver *Object, selector string, args ...uintptr) uintptr {
	var super superStruct
	super.receiver = receiver.Pointer
	super.class = receiver.Class().Superclass().Pointer

	sel := C.sel_registerName(C.CString(selector))
	if len(args) > 0 { // due to cgo calling convention, can't pass an empty array
		return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)))))
	}
	return (uintptr)(unsafe.Pointer(C.objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel)))
}



func loadThySelf(symbol string) *[0]byte {

	fmt.Println("symbol:", symbol)

	this_process := C.dlopen(nil, C.RTLD_NOW)
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


/*
* trimPackage()
* Trim the leading package name (up to including the last index of "."), leave the type
*/
func trimPackage(typePath string) string {

	if i := strings.LastIndex(typePath, "."); i != -1 {
		return typePath[i+1:]
	}
	return typePath
}

// XXX somewhat incomplete 
func objcArgTypeString(argType string) string {

	switch argType {
	case "_Ctype_id":
		return "@"
	case "_Ctype_SEL":
		return ":"
	case "_Ctype_CGRect":
		return "{_NSRect={_NSPoint=ff}{_NSSize=ff}}"
	default:
		return "@"
	}
	return ""
}

