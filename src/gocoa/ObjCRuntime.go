package gocoa

/*
#cgo CFLAGS: -I/System/Library/Frameworks/CoreGraphics.framework/Versions/A/Headers/
#cgo LDFLAGS: -lobjc -lffi

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <ffi/ffi.h>
#include <objc/objc-runtime.h>
#include <CoreGraphics.h>


// for the ambitious, the complete solution would be to build a ffi_type compiler for objc type strings
// XXX irresponsible with memory
static inline ffi_type* gocoa_ffiForNSType(char* nstype) {
	if(strcmp(nstype, "{_NSRect={_NSPoint=ff}{_NSSize=ff}}") == 0) {
		ffi_type*	result = (ffi_type*) malloc(sizeof(ffi_type));
		ffi_type**	type_elements = (ffi_type**) malloc(sizeof(ffi_type*) * 5);
		result->size = result->alignment = 0;
		result->type = FFI_TYPE_STRUCT;
		type_elements[0] = &ffi_type_double;
		type_elements[1] = &ffi_type_double;
		type_elements[2] = &ffi_type_double;
		type_elements[3] = &ffi_type_double;
		type_elements[4] = NULL;
		result->elements = type_elements;
		return result;
	} else if(strcmp(nstype, "B") == 0) {
		return &ffi_type_uint8;
	} else if(strcmp(nstype, "I") == 0) {
		return &ffi_type_uint64;
	} else {
		return &ffi_type_pointer;
	}
}

// beginnings of a proper solution, debugging
static inline void gocoa_I(id self, SEL op, id* result, void* args[], char** types, int argsCount) {
//	printf("gocoa_I(%s, %s, %p", object_getClassName(self), sel_getName(op), *result);
	
	int			i;
	ffi_cif		cif;
	ffi_type	**ffi_types;
	void		**ffi_values;
		
	ffi_types  = (ffi_type **) malloc((argsCount+2)*sizeof(ffi_type *));
	ffi_values = (void **) malloc((argsCount+2)*sizeof(void *));
	
	ffi_types[0] = &ffi_type_pointer;
	ffi_values[0] = &self;
  	ffi_types[1] = &ffi_type_pointer;
	ffi_values[1] = &op;
	
	for (i = 0; i < argsCount; i++) {
//		printf(", %p ('%s')", args[i], types[i]);
		ffi_types[2+i] = gocoa_ffiForNSType(types[i]);
		if(ffi_types[2+i] == &ffi_type_pointer) {
			ffi_values[2+i] = &args[i];
		} else {
			ffi_values[2+i] = args[i];
		}
	}
	
	if (ffi_prep_cif(&cif, FFI_DEFAULT_ABI, argsCount+2, &ffi_type_pointer, ffi_types) == FFI_OK) {
		ffi_call(&cif, (void (*)()) objc_msgSend, result, ffi_values);
	}
		
	free(ffi_types);
	free(ffi_values);
	
//	printf(", %d) out: %p (%s)\n", argsCount, *result, object_getClassName(*result));
	
}


static inline void gocoa_ISuper(struct objc_super *super, SEL op, id* result, void* args[], char** types, int argsCount) {
//	printf("gocoa_ISuper(%s, %s, %p", class_getName(super->class), sel_getName(op), *result);
	
	int			i;
	ffi_cif		cif;
	ffi_type	**ffi_types;
	void		**ffi_values;

	ffi_types  = (ffi_type **) malloc((argsCount+2)*sizeof(ffi_type *));
	ffi_values = (void **) malloc((argsCount+2)*sizeof(void *));
	
	ffi_types[0] = &ffi_type_pointer;
	ffi_values[0] = &super;
  	ffi_types[1] = &ffi_type_pointer;
	ffi_values[1] = &op;
	
	for (i = 0; i < argsCount; i++) {
//		printf(", %p ('%s')", args[i], types[i]);
		ffi_types[2+i] = gocoa_ffiForNSType(types[i]);
		if(ffi_types[2+i] == &ffi_type_pointer) {
			ffi_values[2+i] = &args[i];
		} else {
			ffi_values[2+i] = args[i];
		}
	}
	
	if (ffi_prep_cif(&cif, FFI_DEFAULT_ABI, argsCount+2, &ffi_type_pointer, ffi_types) == FFI_OK) {
		ffi_call(&cif, (void (*)()) objc_msgSendSuper, result, ffi_values);
	} 
	
	free(ffi_types);
	free(ffi_values);
	
//	printf(", %d) out: %p (%s)\n", argsCount, *result, object_getClassName(*result));
}

*/
import "C"

import (
	"math"
	"reflect"
	"runtime"
	"strings"
	"unsafe"
)

/* class implementation ********************************************************************/

type Class uintptr

func ClassForName(name string) Class {
	return (Class)(unsafe.Pointer(C.objc_getClass(C.CString(name))))
}

func (cls Class) cclass() C.Class {
	return (C.Class)(unsafe.Pointer(cls))
}

func (cls Class) cid() C.id {
	return (C.id)(unsafe.Pointer(cls))
}

func (cls Class) RespondsTo(selector string) bool {
	sel := C.sel_registerName(C.CString(selector))
	return (C.class_respondsToSelector(cls.cclass(), sel) == 1)
}

func (cls Class) Name() string {
	return C.GoString(C.class_getName(cls.cclass()))
}

func (cls Class) Super() Class {
	return (Class)(unsafe.Pointer(C.class_getSuperclass(cls.cclass())))
}

func (cls Class) Ivar(name string) Ivar {
	return (Ivar)(unsafe.Pointer(C.class_getInstanceVariable(cls.cclass(), C.CString(name))))
}

func (cls Class) Method(name string) Method {
	sel := C.sel_registerName(C.CString(name))
	return (Method)(unsafe.Pointer(C.class_getClassMethod(cls.cclass(), sel)))
}

func (cls Class) Property(name string) Property {
	return (Property)(unsafe.Pointer(C.class_getProperty(cls.cclass(), C.CString(name))))
}

func (cls Class) Ivars() []Ivar {
	var outCount C.uint
	var ivarPointers []C.Ivar
	p := (C.class_copyIvarList(cls.cclass(), &outCount))
	result := make([]Ivar, outCount)
	if p != nil {
		ivarPointers = (*[1 << 30]C.Ivar)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			result[i] = (Ivar)(unsafe.Pointer(ivarPointers[i]))
		}
		C.free(unsafe.Pointer(p))
	}
	return result
}

func (cls Class) Methods() []Method {
	var outCount C.uint
	var methodPointers []C.Method
	p := (C.class_copyMethodList(cls.cclass(), &outCount))
	result := make([]Method, outCount)
	if p != nil {
		methodPointers = (*[1 << 30]C.Method)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			result[i] = (Method)(unsafe.Pointer(methodPointers[i]))
		}
		C.free(unsafe.Pointer(p))
	}
	return result
}

func (cls Class) Properties() []Property {
	var outCount C.uint
	var properties []C.objc_property_t
	p := (C.class_copyPropertyList(cls.cclass(), &outCount))
	result := make([]Property, outCount)
	if p != nil {
		properties = (*[1 << 30]C.objc_property_t)(unsafe.Pointer(p))[0:outCount]
		for i := 0; i < int(outCount); i++ {
			result[i] = (Property)(unsafe.Pointer(properties[i]))
		}
		C.free(unsafe.Pointer(p))
	}
	return result
}

/* object implementation ********************************************************************/

type Object uintptr

// XXX name collision in C.types, can't use C.Object
func (obj Object) cid() C.id {
	return (C.id)(unsafe.Pointer(obj))
}

func (obj Object) Method(name string) Method {
	sel := C.sel_registerName(C.CString(name))
	return (Method)(unsafe.Pointer(C.class_getInstanceMethod(obj.Class().cclass(), sel)))
}

func (obj Object) Class() Class {
	return (Class)(unsafe.Pointer(C.object_getClass(obj.cid())))
}

// XXX this does not necessarily return an object pointer
func (obj Object) InstanceVariable(name string) Object {
	var val uintptr
	ivar := C.object_getInstanceVariable(obj.cid(), C.CString(name), (*unsafe.Pointer)(unsafe.Pointer(&val)))
	typeenc := C.GoString(C.ivar_getTypeEncoding(ivar))

	if typeenc == "@" {
		// do something
	}
	return (Object)(val)
}

func (obj Object) SetInstanceVariable(name string, val Object) {
	C.object_setInstanceVariable(obj.cid(), C.CString(name), unsafe.Pointer(val))
}

/* class creation methods ------------------------------------------------------------- */

func (cls Class) Subclass(subclassName string) Class {
	class_id := C.objc_allocateClassPair(cls.cclass(), C.CString(subclassName), (C.size_t)(0))
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
		}

		for i := 0; i < numArgs; i++ {
			argType := trimPackage(v.Type().In(i).String())
			types = types + objcArgTypeString(argType)

		}

		sel := C.sel_registerName(C.CString(methodName))
		imp := loadThySelf(impName)
		result := C.class_addMethod(cls.cclass(), sel, imp, C.CString(types))

		return (result == 1)
	}

	return false
}

func (cls Class) AddIvar(ivarName string, ivarClass Class) bool {

	types := objcArgTypeString(ivarClass.Name())
	size := (C.size_t)(unsafe.Sizeof(cls))
	alignment := (C.uint8_t)(math.Log2((float64)(unsafe.Sizeof(cls))))
	result := C.class_addIvar(cls.cclass(), C.CString(ivarName), size, alignment, C.CString(types))

	return (result == 1)
}

func (cls Class) Register() {
	C.objc_registerClassPair(cls.cclass())
}

/* method implementation ************************************************************** */

type Selector uintptr

func SelectorForName(name string) Selector {
	return (Selector)(unsafe.Pointer(C.sel_registerName(C.CString(name))))
}

func (sel Selector) csel() C.SEL {
	return (C.SEL)(unsafe.Pointer(sel))
}

func (sel Selector) Name() string {
	return C.GoString(C.sel_getName(sel.csel()))
}

/* method implementation ************************************************************** */

type Method uintptr

func (mthd Method) cmethod() C.Method {
	return (C.Method)(unsafe.Pointer(mthd))
}

func (mthd Method) ArgumentCount() int {
	return (int)(C.method_getNumberOfArguments(mthd.cmethod()))
}

func (mthd Method) ArgumentType(index int) string {
	var dst_len C.size_t
	var dst *C.char
	C.method_getArgumentType(mthd.cmethod(), (C.uint)(index), dst, dst_len)
	return C.GoString(dst)
}

func (mthd Method) Name() string {
	return C.GoString(C.sel_getName(C.method_getName(mthd.cmethod())))
}

/* ivar methods *********************************************************************** */

type Ivar uintptr

func (ivr Ivar) Name() string {
	return C.GoString(C.ivar_getName((C.Ivar)(unsafe.Pointer(ivr))))
}

/* property methods ******************************************************************* */

type Property uintptr

func (prop Property) Name() string {
	return C.GoString(C.property_getName((C.objc_property_t)(unsafe.Pointer(prop))))
}

func (prop Property) Attributes() string {
	return C.GoString(C.property_getAttributes((C.objc_property_t)(unsafe.Pointer(prop))))
}


/* messaging functions *************************************************************** */

// Better understanding of when to employ objc_msgSend_stret and objc_msgSend_fret
// is going to be a prerequisite to calling this finished 

func constructArgs(args...Passable) ([]unsafe.Pointer, []*C.char) {
	items := make([]unsafe.Pointer, len(args))
	types := make([]*C.char, len(args))

	for i := 0; i < len(args); i++ {
		types[i] = C.CString(args[i].TypeString())
		items[i] = unsafe.Pointer(args[i].Ptr())
	}
	return items, types
}

type superStruct struct {
	receiver Object
	class    Class
}


// one possibility is that this always returns an object, always converting structs to 
// pointers on output, assuming that the runtime can generally use them on input
func (obj Object) Call(selector string, args ...Passable) Object {
//	if obj == 0 {
//		panic("can't call with nil class pointer")
//	}
	
	var result C.id
	sel := C.sel_registerName(C.CString(selector))
	
	if len(args) > 0 {
		items, types := constructArgs(args...)
		C.gocoa_I(obj.cid(), sel, &result, &items[0], (**C.char)(&types[0]), (C.int)(len(items)))
	} else {
		result = C.objc_msgSend(obj.cid(), sel)
	}
		
	// XXX output conversion needed
	return (Object)(unsafe.Pointer(result))
}

/*
* The distinction here involves messaging to superclasses, with message receipt to the subclass. This
* requires initing a structure that refers to both the receiver and its superclass. 
 */
func (obj Object) CallSuper(method string, args ...Passable) Object {
	if obj == 0 {
		panic("can't call super with nil object pointer")
	}
	if obj.Class().Super() == 0 {
		panic("can't call super with nil superclass pointer")
	}
	
	var super superStruct
	super.receiver = obj
	super.class = obj.Class().Super()
		
	var result C.id
	sel := C.sel_registerName(C.CString(method))
	
	if len(args) > 0 {
		items, types := constructArgs(args...)
		C.gocoa_ISuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel, &result, &items[0], (**C.char)(&types[0]), (C.int)(len(items)))
	} else {
		result = C.objc_msgSendSuper((*C.struct_objc_super)(unsafe.Pointer(&super)), sel)
	}
	
	// XXX output conversion needed
	return (Object)(unsafe.Pointer(result))
}


func (cls Class) Instance(method string, args ...Passable) Object {
	if cls == 0 {
		panic("can't instantiate with nil class pointer")
	}
	return (Object)(cls).Call(method, args...)
}


/*
* loadThySelf()
* Go doesn't support dynamic linking. However, it supports a C interface that supports
* dynamic linking. And it supports symbol export allowing callbacks into go functions
* using a C calling convention. So, Go supports dynamic linking. 
 */
func loadThySelf(symbol string) *[0]byte {

	this_process := C.dlopen(nil, C.RTLD_NOW)
	if this_process == nil {
		panic(C.GoString(C.dlerror()))
	}

	symbol_address := C.dlsym(this_process, C.CString(symbol))
	if symbol_address == nil {
		panic(C.GoString(C.dlerror()))
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
