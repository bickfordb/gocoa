package gocoa

/*
#cgo CFLAGS: -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -lobjc
#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>

CGRect*	fpCGRect(void* in) { return (CGRect*)in; }
id*		fpId	(void* in) { return (id*)in; }


static inline id gocoa_I(id self, SEL op, void* items[], char** types, int argsCount) {
	printf("gocoa_I(%d, %d)", -1, argsCount);
	int i=0;
	for(; i<argsCount; i++) {
		printf(", type:'%s'", types[i]);
	}
	printf(")\n");

	switch (argsCount) {
		case 1: return objc_msgSend(self, op, *fpId(items[0]));
//		case 2: return objc_msgSend(self, op, items[0], items[1]);
//		case 3: return objc_msgSend(self, op, items[0], items[1], items[2]);
//		case 4: return objc_msgSend(self, op, items[0], items[1], items[2], items[3]);
//		case 5: return objc_msgSend(self, op, items[0], items[1], items[2], items[3], items[4]);
		default: return objc_msgSend(self, op);
	}

}



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
	"fmt"
	"math"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

/* class implementation ********************************************************************/

type Class uintptr

func (cls Class) classPointer() C.Class {
	return (C.Class)(unsafe.Pointer(cls))
}

func (cls Class) idPointer() C.id {
	return (C.id)(unsafe.Pointer(cls))
}

func (cls Class) respondsTo(selector string) bool {
	sel := C.sel_registerName(C.CString(selector))
	return (C.class_respondsToSelector(cls.classPointer(), sel) == 1)
}

func (cls Class) Name() string {
	return C.GoString(C.class_getName(cls.classPointer()))
}

func (cls Class) Super() Class {
	return (Class)(unsafe.Pointer(C.class_getSuperclass(cls.classPointer())))
}

func (cls Class) Instance(method string, args ...Object) Object {
	obj := (Object)(cls)
	return obj.Call(method, args...)
}

func (cls Class) InstanceR(method string, arg NSRect) Object {
	obj := (Object)(cls)
	return obj.CallR(method, arg)
}

func (cls Class) Method(name string) Method {
	sel := C.sel_registerName(C.CString(name))
	return (Method)(unsafe.Pointer(C.class_getClassMethod(cls.classPointer(), sel)))
}

func (cls Class) Property(name string) Property {
	return (Property)(unsafe.Pointer(C.class_getProperty(cls.classPointer(), C.CString(name))))
}

func (cls Class) ListInstanceVariables() {

	fmt.Println(cls.Name(), ": ivars")

	var outCount C.uint
	var ivarPointers []C.Ivar

	p := (C.class_copyIvarList(cls.classPointer(), &outCount))

	if p != nil {
		ivarPointers = (*[1 << 30]C.Ivar)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := (Ivar)(unsafe.Pointer(ivarPointers[i]))
			fmt.Println("\t", tmp.Name())
		}
		C.free(unsafe.Pointer(p))
	}
}

func (cls Class) ListMethods() {

	fmt.Println(cls.Name(), ": methods")

	var outCount C.uint
	var methodPointers []C.Method

	p := (C.class_copyMethodList(cls.classPointer(), &outCount))

	if p != nil {
		methodPointers = (*[1 << 30]C.Method)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := (Method)(unsafe.Pointer(methodPointers[i]))
			fmt.Println("\t", tmp.Name())
		}
		C.free(unsafe.Pointer(p))
	}
}

func (cls Class) ListProperties() {

	fmt.Println(cls.Name(), ": properties")

	var outCount C.uint
	var properties []C.objc_property_t

	p := (C.class_copyPropertyList(cls.classPointer(), &outCount))

	if p != nil {
		properties = (*[1 << 30]C.objc_property_t)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			tmp := (Property)(unsafe.Pointer(properties[i]))
			fmt.Println("\t", tmp.Name(), "(", tmp.Attributes(), ")")
		}
		C.free(unsafe.Pointer(p))
	}
}

/* object implementation ********************************************************************/

type Object uintptr

// XXX name collision in C.types, can't use C.Object
func (obj Object) idPointer() C.id {
	return (C.id)(unsafe.Pointer(obj))
}

func (obj Object) getMethod(name string) Method {
	sel := C.sel_registerName(C.CString(name))
	return (Method)(unsafe.Pointer(C.class_getInstanceMethod(obj.Class().classPointer(), sel)))
}

func (obj Object) Class() Class {
	return (Class)(unsafe.Pointer(C.object_getClass(obj.idPointer())))
}

// XXX this does not necessarily return an object pointer
func (obj Object) InstanceVariable(name string) Object {
	var val uintptr
	ivar := C.object_getInstanceVariable(obj.idPointer(), C.CString(name), (*unsafe.Pointer)(unsafe.Pointer(&val)))
	typeenc := C.GoString(C.ivar_getTypeEncoding(ivar))

	if typeenc == "@" {
		return (Object)(val)
	}
	return 0
}

func (obj Object) SetInstanceVariable(name string, val Object) {
	C.object_setInstanceVariable(obj.idPointer(), C.CString(name), unsafe.Pointer(val))
}

/* class creation methods ------------------------------------------------------------- */

func (cls Class) Subclass(subclassName string) Class {
	class_id := C.objc_allocateClassPair(cls.classPointer(), C.CString(subclassName), (C.size_t)(0))
	return (Class)(unsafe.Pointer(class_id))
}

func (cls Class) AddMethod(methodName string, implementor interface{}) bool {

	v := reflect.ValueOf(implementor)

	if (v.Kind() == reflect.Func) && (v.Type().NumIn() > 1) {

		types := "v"
		impName := trimPackage(runtime.FuncForPC(v.Pointer()).Name())
		numArgs := v.Type().NumIn()

		if v.Type().NumOut() == 1 {
			types = objcArgTypeString(trimPackage(v.Type().Out(0).String()))
			fmt.Println("typeout:", types)
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

func (cls Class) AddIvar(ivarName string, ivarClass Class) bool {

	types := objcArgTypeString(ivarClass.Name())
	size := (C.size_t)(unsafe.Sizeof(cls))
	alignment := (C.uint8_t)(math.Log2((float64)(unsafe.Sizeof(cls))))
	result := C.class_addIvar(cls.classPointer(), C.CString(ivarName), size, alignment, C.CString(types))

	return (result == 1)
}

func (cls Class) Register() {
	C.objc_registerClassPair(cls.classPointer())
}

/* method implementation ************************************************************** */

type Selector uintptr

func SelectorForName(name string) Selector {
	return (Selector)(unsafe.Pointer(C.sel_registerName(C.CString(name))))
}

func (sel Selector) selPointer() C.SEL {
	return (C.SEL)(unsafe.Pointer(sel))
}

func (sel Selector) Name() string {
	return C.GoString(C.sel_getName(sel.selPointer()))
}

/* method implementation ************************************************************** */

type Method uintptr

func (mthd Method) methodPointer() C.Method {
	return (C.Method)(unsafe.Pointer(mthd))
}

func (mthd Method) ArgumentCount() int {
	return (int)(C.method_getNumberOfArguments(mthd.methodPointer()))
}

func (mthd Method) ArgumentType(index int) string {
	var dst_len C.size_t
	var dst *C.char
	C.method_getArgumentType(mthd.methodPointer(), (C.uint)(index), dst, dst_len)
	return C.GoString(dst)
}

func (mthd Method) Name() string {
	return C.GoString(C.sel_getName(C.method_getName(mthd.methodPointer())))
}

/***************************************************************************************
* Ivar
* An Ivar is an Objective-C structure describing an instance variable, accessible via
* object.InstanceVariable(name)
 */

type Ivar uintptr

func (ivr Ivar) Name() string {
	return C.GoString(C.ivar_getName((C.Ivar)(unsafe.Pointer(ivr))))
}

/* property methods ********************************************************************/

type Property uintptr

func (prop Property) Name() string {
	return C.GoString(C.property_getName((C.objc_property_t)(unsafe.Pointer(prop))))
}

func (prop Property) Attributes() string {
	return C.GoString(C.property_getAttributes((C.objc_property_t)(unsafe.Pointer(prop))))
}

/* utility methods ********************************************************************/

func ClassForName(name string) Class {
	return (Class)(unsafe.Pointer(C.objc_getClass(C.CString(name))))
}

func ObjectForId(object_id uintptr) Object {
	return (Object)(object_id)
}

/* messaging functions *************************************************************** 
*
* There are a few hairy issues with the messaging funtions. First, the C funtion 
* interface doesn't translate Go variadic functions, which is just as well, because
* C variadic functions can accept variables of any type.
*
* The ultimate solution is to define a type of linked list in C that accepts arbitrary 
* types, and one Go object method: Call(selector string, args...interface{}) that 
* iterates over the argument list and uses reflect to unpack the arguments into the C 
* list, then calls objc_msgSend.
*
* As I understand, other platforms possibly including the iPhone may require
* objc_msgSend_stret and objc_msgSend_fret, and they certainly have slightly different 
* data types, which will be something for consideration when fixing the following.
*
* As of yet, it's a mess that seems to work.
 */

func (obj Object) I(selector string, args ...Passable) Passable {

	items := make([]unsafe.Pointer, len(args))
	types := make([]*C.char, len(args))

	for i := 0; i < len(args); i++ {

		value := reflect.ValueOf(args[i])

		types[i] = C.CString(args[i].TypeString())

		fmt.Println("value.String()", value.String()) // 

		if value.String() == "<gocoa.Object Value>" {
			items[i] = unsafe.Pointer(args[i].Id())
		} else {
			items[i] = unsafe.Pointer(&(args[i].Bytes()[0]))
		}
	}

	sel := C.sel_registerName(C.CString(selector))
	var result C.id

	result = C.gocoa_I(obj.idPointer(), sel, &items[0], (**C.char)(&types[0]), (C.int)(len(items)))

	// XXX output conversion needed
	return (Object)(unsafe.Pointer(result))
}

// clumsy hacks abound

type superStruct struct {
	receiver Object
	class    Class
}

// XXX all of these are pointless, fix
func (obj Object) CallR(method string, arg NSRect) Object {
	sel := C.sel_registerName(C.CString(method))
	return (Object)(unsafe.Pointer(C.gocoa_objc_msgSendR(obj.idPointer(), sel, arg.CGRect())))
}

func (obj Object) CallI(method string, arg NSUInteger) Object {
	sel := C.sel_registerName(C.CString(method))
	return (Object)(unsafe.Pointer(C.gocoa_objc_msgSendI(obj.idPointer(), sel, (C.long)(arg))))
}

/*
* Call()
* Notice that you have to pass a pointer to the first array element to match the c array calling convention.
 */
func (obj Object) Call(method string, args ...Object) Object {
	sel := C.sel_registerName(C.CString(method))
	if len(args) > 0 { // due to cgo calling convention, can't pass an empty array
		return (Object)(unsafe.Pointer(C.gocoa_objc_msgSend(obj.idPointer(), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)))))
	}
	return (Object)(unsafe.Pointer(C.objc_msgSend(obj.idPointer(), sel)))
}

/*
* CallSuper()
* The distinction here involves messaging to superclasses, with message receipt to the subclass. This
* requires initing a structure that refers to both the receiver and its superclass. 
 */
func (obj Object) CallSuper(method string, args ...Object) Object {
	var super superStruct
	super.receiver = obj
	super.class = obj.Class().Super()

	sel := C.sel_registerName(C.CString(method))
	if len(args) > 0 { // due to cgo calling convention, can't pass an empty array
		return (Object)(unsafe.Pointer(C.gocoa_objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, (*C.id)(unsafe.Pointer(&args[0])), (C.int)(len(args)))))
	}
	return (Object)(unsafe.Pointer(C.objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel)))
}

func (obj Object) CallSuperR(method string, arg NSRect) Object {
	var super superStruct
	super.receiver = obj
	super.class = obj.Class().Super()
	sel := C.sel_registerName(C.CString(method))
	return (Object)(unsafe.Pointer(C.gocoa_objc_msgSendSuperR((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, arg.CGRect())))
}

/*
* loadThySelf()
* Go doesn't support dynamic linking. However, it supports a C interface that supports
* dynamic linking. And it supports symbol export allowing callbacks into go functions
* using a C calling convention. So, Go supports dynamic linking. 
*
* XXX this function will probably fail with a panic rather than a message when I figure 
* out why it's unreliable. 
 */
func loadThySelf(symbol string) *[0]byte {

	fmt.Println("symbol:", symbol)

	this_process := C.dlopen(nil, C.RTLD_NOW)
	if this_process == nil {
		fmt.Println("********** error:", C.GoString(C.dlerror()))
	}

	symbol_address := C.dlsym(this_process, C.CString(symbol))
	if symbol_address == nil {
		fmt.Println("********** error:", C.GoString(C.dlerror()))
	}

	C.dlclose(this_process)
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
