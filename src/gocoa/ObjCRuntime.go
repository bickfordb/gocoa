package gocoa

/*
#include <stdlib.h>
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
//#cgo CFLAGS: -I/usr/include -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
//#cgo LDFLAGS: -lobjc -framework CoreGraphics
import "C"

import (
	"fmt"
	"unsafe"
	"encoding/binary"
	"bytes"
)

/* class implementation *************************************************************** */

type Class struct { id uintptr }

// XXX maybe just a macro for converting between C.Class and C.id ?
func (cls *Class) ptr() C.Class {
	return (C.Class)(unsafe.Pointer(cls.id))
}
	
func (cls *Class) Name() string {
	return C.GoString(C.class_getName(cls.ptr()))
}

func (cls *Class) Superclass() *Class {
	class_id := (uintptr)(unsafe.Pointer(C.class_getSuperclass(cls.ptr())))
	return &Class{class_id}
}

func (cls *Class) Instance(calling string, args...uintptr) *Object {
	return &Object{ msgSend(&Object{cls.id}, calling, args...) }
}
	
func (cls *Class) Method(name string) *Method {
	sel := C.sel_registerName(C.CString(name))
	method_id := (uintptr)(unsafe.Pointer(C.class_getClassMethod(cls.ptr(), sel)))
	return &Method{method_id}
}

/*
* respondsTo()
* 
*/
func (cls *Class) respondsTo(selector string) bool {
	sel := C.sel_registerName(C.CString(selector))
	return (C.class_respondsToSelector(cls.ptr(), sel) == 1)
}



type Ivar struct { id uintptr }

func (ivr *Ivar) ptr() C.Ivar {
	return (C.Ivar)(unsafe.Pointer(ivr.id))
}

func (ivr *Ivar) Name() string {
	return C.GoString(C.ivar_getName(ivr.ptr()))
}

func (cls *Class) ListInstanceVariables() {
	
	fmt.Println(cls.Name(), ": ivars")

	var outCount C.uint
	var ivarPointers []C.Ivar
	
	p := (C.class_copyIvarList(cls.ptr(), &outCount))
	
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
	
	p := (C.class_copyMethodList(cls.ptr(), &outCount))
	
	if p != nil {
		methodPointers = (*[1 << 30]C.Method)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := Method{(uintptr)(unsafe.Pointer(methodPointers[i]))}
			fmt.Println("\t", tmp.Name())
		}
		C.free(unsafe.Pointer(p))
	}
}

/* object implementation ************************************************************** */

type Object struct {
	id		uintptr
}

func NewObject(object_id uintptr) *Object {
	return &Object{object_id}
}

func (obj *Object) IsClass() *Class {
	return &Class{obj.id}
}

func (obj *Object) getMethod(name string) *Method {
	sel := C.sel_registerName(C.CString(name))
	method_id := (uintptr)(unsafe.Pointer(C.class_getInstanceMethod(obj.Class().ptr(), sel)))
	return &Method{method_id}
}
	
func (obj *Object) Class() *Class {
	object_id := (uintptr)(unsafe.Pointer(C.object_getClass((C.id)(obj.Ptr()))))
	return &Class{object_id}
}
	
// XXX should take arguments 
func (obj *Object) Init() *Object {
//	msgSend(obj.id, "init")
	sel := C.sel_registerName(C.CString("init"))
	obj.id = (uintptr)(unsafe.Pointer(C.objc_msgSend((C.id)(unsafe.Pointer(obj.id)), sel)))
	return obj
}


func (obj *Object) CallR(method string, arg []byte) *Object {
	var copiedArg C.CGRect	
	buf := bytes.NewBuffer(arg)
/*	err := */binary.Read(buf, binary.LittleEndian, copiedArg)
	return &Object{msgSendR(obj, method, copiedArg)}
}

func (obj *Object) CallI(method string, arg NSUInteger) *Object {
	return &Object{msgSendI(obj, method, arg)}
}
	
func (obj *Object) Call(method string, args...uintptr) *Object {
	return &Object{msgSend(obj, method, args...)}
}

func (obj *Object) CallSuper(method string, args...uintptr) *Object {
	return &Object{msgSendSuper(obj, method, args...)}
}

func (obj *Object) CallSuperR(method string, arg []byte) *Object {
	var copiedArg C.CGRect	
	buf := bytes.NewBuffer(arg)
/*	err := */binary.Read(buf, binary.LittleEndian, copiedArg)
	return &Object{msgSendSuperR(obj, method, copiedArg)}
}

// XXX this does not necessarily return an object pointer
func (obj *Object) InstanceVariable(name string) *Object {
	var val uintptr
	ivar := C.object_getInstanceVariable(obj.Ptr(), C.CString(name), (*unsafe.Pointer)(unsafe.Pointer(&val)))
	typeenc := C.GoString(C.ivar_getTypeEncoding(ivar))
	
//	fmt.Println("typeenc:", typeenc)
	
	if(typeenc == "@") {
		return &Object{val}
	}
	return nil
}

func (obj *Object) SetInstanceVariable(name string, val *Object) {
	C.object_setInstanceVariable(obj.Ptr(), C.CString(name), unsafe.Pointer(val.Ptr()))
}

func (obj *Object) ListInstanceVariables() {
	obj.Class().ListInstanceVariables()
}

func (obj *Object) ListMethods() {
	obj.Class().ListMethods()
}

// XXX name collision in C.types, can't use C.Object, need to create a namespace
// however, C.types seem to coerce
func (obj *Object) Ptr() C.id { 
	return (C.id)(unsafe.Pointer(obj.id))
}

func (obj *Object) Id() uintptr { 
	return obj.id
}

	
/* method implementation ************************************************************** */

type Method struct { id uintptr }

func (mthd *Method) ptr() C.Method {
	return (C.Method)(unsafe.Pointer(mthd.id))
}

func (mthd *Method) ArgumentCount() int {
	return (int)(C.method_getNumberOfArguments(mthd.ptr()))
}
	
func (mthd *Method) ArgumentType(index int) string {
	var dst_len C.size_t
	var dst *C.char
	C.method_getArgumentType(mthd.ptr(), (C.uint)(index), dst, dst_len)
	return C.GoString(dst)
}
	
func (mthd *Method) Name() string {
	return C.GoString(C.sel_getName(C.method_getName(mthd.ptr())))
}
	

/* utility methods ******************************************************************* */


// ------------- functions for interacting with classes

func ClassForName(name string) *Class {
	class_id := (uintptr)(unsafe.Pointer(C.objc_getClass(C.CString(name))))
	return &Class{class_id}
}

// clumsy hacks abound
func msgSendI(receiver *Object, selector string, number NSUInteger) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendI((C.id)(unsafe.Pointer(receiver.id)), sel, (C.long)(number))))
}

func msgSendR(receiver *Object, selector string, rect C.CGRect) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendR((C.id)(unsafe.Pointer(receiver.id)), sel, rect)))
}

type superStruct struct {
	receiver	uintptr
	class		uintptr
}

func msgSendSuperR(receiver *Object, selector string, rect C.CGRect) uintptr {
	var super superStruct
	super.receiver = receiver.id
	super.class = receiver.Class().Superclass().id
	sel := C.sel_registerName(C.CString(selector))
	return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendSuperR((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, rect )))
}



/*
* msgSend()
* Notice that you have to pass a pointer to the first array element to match the c array calling convention.
* The gocoa_ bridge function essentially, clumsily replicates Apple's deprecated objc_msgSendv().
*/ 
func msgSend(receiver *Object, selector string, args...uintptr) uintptr {
	sel := C.sel_registerName(C.CString(selector))
	if len(args) > 0 {		// due to cgo calling convention, can't pass an empty array
		return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSend((C.id)(unsafe.Pointer(receiver.id)), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)) )))
	}
	return (uintptr)(unsafe.Pointer(C.objc_msgSend((C.id)(unsafe.Pointer(receiver.id)), sel)))
}

/*
* msgSendSuper()
* The distinction here involves messaging to superclasses, with message receipt to the subclass. This
* requires initing a structure that refers to both the receiver and its superclass. 
*/
func msgSendSuper(receiver *Object, selector string, args...uintptr) uintptr {
	var super superStruct
	super.receiver = receiver.id
	super.class = receiver.Class().Superclass().id
	
	sel := C.sel_registerName(C.CString(selector))
	if len(args) > 0 {		// due to cgo calling convention, can't pass an empty array
		return (uintptr)(unsafe.Pointer(C.gocoa_objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)) )))
	}
	return (uintptr)(unsafe.Pointer(C.objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel)))
}

