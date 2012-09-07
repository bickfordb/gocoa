package gocoa

import (
	"fmt"
//	"runtime"
	"testing"
//	"unsafe"
)

func printlist(name string, methods []Method) {
	fmt.Println(name)
	for i:=0; i<len(methods); i++ {
		fmt.Println("\t", methods[i].Name())
	}
}

func printproperties(name string, properties []Property) {
	fmt.Println(name)
	for i:=0; i<len(properties); i++ {
		fmt.Println("\t", properties[i].Name())
	}
}


	
/*
* If it has a test, it's an interface. Most of these are naive tests, when I start discovering
* more interesting failure scenarios, they'll be incorporated.
*/


func Test_ClassForName(t *testing.T) {
	nso := ClassForName("NSObject").Instance("alloc").Call("init")
	if nso == 0 {
		t.Log("nso", nso)
		t.Fail()
	}
}

func Test_SelectorForName(t *testing.T) {
	sel := SelectorForName("alloc")
	if sel == 0 {
		t.Log("sel", sel)
		t.Fail()
	}
}


func Test_Class_RespondsTo(t *testing.T) {
	nso := ClassForName("NSObject").Instance("alloc").Call("init")
	if !nso.Class().RespondsTo("classForPortCoder") {
		t.Log("classForPortCoder not found")
		t.Fail()
	}
	if nso.Class().RespondsTo("bogusMethod:withBogusArgs:") {
		t.Log("bogusMethod:withBogusArgs:, found")
		t.Fail()
	}
}

func Test_Class_Name(t *testing.T) {
	name := ClassForName("NSObject").Name()
	if name != "NSObject" {
		t.Log("ClassForName", name)
		t.Fail()
	}
}


func Test_Class_Super(t *testing.T) {
	class := ClassForName("NSObject")
	superclass := class.Super()
	name := superclass.Name()
	
	if superclass == 0 {
		t.Log("superclass", superclass)
		t.Fail()
	}
	
	if len(name) == 0 || name != "Object" {
		t.Log("name", name)
		t.Fail()
	}
}

func Test_Class_Instance(t *testing.T) {
	nso := ClassForName("NSObject").Instance("alloc")
	if nso == 0 {
		t.Fail()
	}
	if !nso.Class().RespondsTo("init") {
		t.Fail()
	}
	if nso.Class().Name() != "NSObject" {
		t.Fail()
	}
}

func Test_Class_Ivar(t *testing.T) {
	object := ClassForName("NSObject").Instance("alloc")
	ivar := object.Method("isa")
	if ivar == 0 {
		t.Log("ivar", ivar)
		t.Fail()
	}
}

func Test_Class_Method(t *testing.T) {
	class := ClassForName("NSObject")
	method := class.Method("alloc")
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
}

func Test_Class_Property(t *testing.T) {
	sel := "terminationHandler"
	nst := ClassForName("NSTask")
	property := nst.Property(sel)
	
	if nst == 0 {
		t.Log("nst", nst)
		t.Fail()
	}
	if nst.Name() != "NSTask" {
		t.Log("nst", nst.Name())
		t.Fail()
	}
	
	if property == 0 {
		t.Log("property", property)
		t.Fail()
	}

}

func Test_Class_Ivars(t *testing.T) {
	class := ClassForName("NSObject")
	ivars := class.Ivars()
	if len(ivars) < 1 {
		t.Log("ivars", ivars)
		t.Fail()
	}
}

func Test_Class_Methods(t *testing.T) {
	class := ClassForName("NSObject")
	methods := class.Methods()
	if len(methods) < 1 {
		t.Log("methods", methods)
		t.Fail()
	}
}

func Test_Class_Properties(t *testing.T) {	
	class := ClassForName("NSTask")
	properties := class.Properties()
	if len(properties) < 1 {
		t.Log("properties", properties)
		t.Fail()
	}
}


func Test_Class_Subclass(t *testing.T) {
	subclass := ClassForName("NSObject").Subclass("NSSubclass")
	if subclass == 0 {
		t.Log("subclass", subclass)
		t.Fail()
	}
	if subclass.Name() != "NSSubclass" {
		t.Log("subclass", subclass.Name())
		t.Fail()
	}
	
	instance := subclass.Instance("alloc")
	if instance == 0 {
		t.Log("instance", instance)
		t.Fail()
	}
	if instance.Class().Name() != "NSSubclass" {
		t.Log("instance", instance.Class().Name())
		t.Fail()
	}
}


func Test_Class_AddMethod(t *testing.T) {
	t.Log("gotta write a harness, no cgo in test")
	t.Fail()
}

func Test_Class_AddIvar(t *testing.T) {
	
	dude := NSString("dude")
	
	subclass := ClassForName("NSObject").Subclass("NSSubclass")
	subclass.AddIvar("ivarOne", dude.Class())
	subclass.AddIvar("ivarTwo", subclass)
	
	instance := subclass.Instance("alloc").Call("init")
	
	instance.SetInstanceVariable("ivarOne", dude)
	instance.SetInstanceVariable("ivarTwo", instance)
	
	ivarOne := instance.InstanceVariable("ivarOne")
	ivarTwo := instance.InstanceVariable("ivarTwo")
	
	if ivarOne == 0 {
		t.Log("ivarOne", ivarOne)
		t.Fail()
	}
	if ivarTwo == 0 {
		t.Log("ivarTwo", ivarTwo)
		t.Fail()
	}
	
	if NSStringToString(ivarOne) != "dude" {
		t.Log("ivarOne", NSStringToString(ivarOne))
		t.Fail()
	}
	if ivarTwo.Class().Name() != "NSSubclass" {
		t.Log("ivarTwo", ivarTwo.Class().Name())
		t.Fail()
	}
	
}

func Test_Class_Register(t *testing.T) {
	
	classreg := ClassForName("NSObject").Subclass("NSSubclass")
	classreg.Register()
	
	subclass := ClassForName("NSSubclass")
	superclass := subclass.Super()
	instance := subclass.Instance("alloc").Call("init")
	
	if subclass == 0 {
		t.Log("subclass", subclass)
		t.Fail()
	}
	if superclass == 0 {
		t.Log("superclass", superclass)
		t.Fail()
	}
	if instance == 0 {
		t.Log("instance", instance)
		t.Fail()
	}
	
}

func Test_Object_Method(t *testing.T) {
	object := ClassForName("NSObject").Instance("alloc")
	method := object.Method("init")
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
}

func Test_Object_Class(t *testing.T) {
	object := ClassForName("NSObject").Instance("alloc")
	class := object.Class()
	if class == 0 {
		t.Log("class", class)
		t.Fail()
	}
}

func Test_Object_InstanceVariable(t *testing.T) {
	object := ClassForName("NSDictionary").Instance("alloc").Call("initWithObject:forKey:", NSString("foo"), NSString("bar"))
	ivar := object.InstanceVariable("_used")
	if ivar == 0 {
		t.Log("ivar", ivar)
		t.Fail()
	}
}

func Test_Object_SetInstanceVariable(t *testing.T) {
	
	dude := NSString("foo")
	subclass := ClassForName("NSObject").Subclass("NSSubclass")
	
	if subclass == 0 {
		t.Log("subclass", subclass)
		t.Fail()
	}
	subclass.AddIvar("someIvar", dude.Class())
	// XXX causes runtime hang:
/*	subclass.Register()
	
	instance := subclass.Instance("alloc").Call("init")
	instance.SetInstanceVariable("someIvar", dude)
	ivar := instance.InstanceVariable("someIvar")
	if ivar == 0 {
		t.Log("ivar", ivar)
		t.Fail()
	}*/
	
	t.Log("runtime hang case not tested")
	t.Fail()
}



func Test_Object_Call(t *testing.T) {

	t.Log("initbundle")
	bundle := ClassForName("NSBundle").Instance("alloc").Call("initWithPath:", NSString("."))
	if bundle == 0 {
		t.Log("bundle", bundle)
		t.Fail() 
	}
	
	t.Log("initarray")
	someArray := ClassForName("NSMutableArray").Instance("alloc").Call("init")
	if someArray == 0 {
		t.Log("someArray", someArray)
		t.Fail()
	}
	
	t.Log("getcount")
	someCount := (NSUInteger)(someArray.Call("count"))
	if someCount != 0 {
		t.Log("someCount", someCount)
		t.Fail()
	}
	
	// hangs
	nilresult := bundle.Call("initWithPath:", NSString("A String"), someArray, someCount, MakeNSBoolean(true))
	if nilresult != 0 {
		t.Fail()
	} else {
		t.Log("nilresult", nilresult)
	}
	
} 


func Test_Object_CallSuper(t *testing.T) {
	
	bundle := ClassForName("NSBundle").Instance("alloc").Call("initWithPath:", NSString("."))
//	nso := ClassForName("NSObject").Instance("alloc").Call("init")
	
	if bundle == 0 {
		t.Log("bundle", bundle)
		t.Fail() 
	}

// XXX hang
/* 	nso2 := bundle.CallSuper("init")
	
	if nso2 == 0 {
		t.Log("nso2", nso2)
		t.Fail() 
	}
	
	// XXX either I'm misunderstanding how this is supposed to work or it's broken
	if nso.Class().Name() != nso2.Class().Name() {
		t.Log("nso", nso.Class().Name(), "nso2", nso2.Class().Name())
		t.Fail() 
	}
	
	if len(nso.Class().Methods()) != len(nso2.Class().Methods()) {
		t.Log("nso", len(nso.Class().Methods()), "nso2", len(nso2.Class().Methods()))
		t.Fail() 
	}*/
	
	t.Log("runtime hang case not tested")
	t.Fail()
}




func Test_Selector_Name(t *testing.T) {
	sel := SelectorForName("alloc")
	name := sel.Name()
	if len(name) == 0 || name != "alloc" {
		t.Log("sel", sel)
		t.Fail()
	}
}



func Test_Method_ArgumentCount(t *testing.T) {
	class := ClassForName("NSObject")
	method := class.Method("replacementObjectForPortCoder:")
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
	if method.ArgumentCount() != 3 {
		t.Log("method", method.Name(), method.ArgumentCount(), "arguments")
		t.Fail()
	}
	
}

func Test_Method_ArgumentType(t *testing.T) {
	class := ClassForName("NSObject")
	method := class.Method("replacementObjectForPortCoder:")
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
	if method.ArgumentType(3) != "@" {
		t.Log("method", method.Name(), "arg", 1, method.ArgumentType(3), "type")
		t.Fail()
	}
	
}

func Test_Method_Name(t *testing.T) {
	sel := "replacementObjectForPortCoder:"
	class := ClassForName("NSObject")
	method := class.Method(sel)
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
	if method.Name() != sel {
		t.Log("method", method.Name())
		t.Fail()
	}
}


func Test_Ivar_Name(t *testing.T) {
	sel := "isa"
	class := ClassForName("NSObject")
	ivar := class.Ivar(sel)
	if ivar == 0 || ivar.Name() != sel {
		t.Log("ivar", ivar)
		t.Fail()
	}
}



func Test_Property_Name(t *testing.T) {
	sel := "terminationHandler"
	nst := ClassForName("NSTask")
	property := nst.Property(sel)
	
	if nst == 0 {
		t.Log("nst", nst)
		t.Fail()
	}
	if nst.Name() != "NSTask" {
		t.Log("nst", nst.Name())
		t.Fail()
	}
	
	if property == 0 {
		t.Log("property", property)
		t.Fail()
	}
	if property.Name() != sel {
		t.Log("property", property.Name())
		t.Fail()
	}
	
}

func Test_Property_Attributes(t *testing.T) {
	sel := "terminationHandler"
	nst := ClassForName("NSTask")
	
	
	if nst == 0 {
		t.Log("nst", nst)
		t.Fail()
	}
	if nst.Name() != "NSTask" {
		t.Log("nst", nst.Name())
		t.Fail()
	}
	
	property := nst.Property(sel)
	
	if property == 0 {
		t.Log("property", property)
		t.Fail()
	}
	if property.Name() != sel {
		t.Log("property", property.Name())
		t.Fail()
	}
	
	attributes := property.Attributes()
	if len(attributes) == 0 {
		t.Log("attributes", attributes)
		t.Fail()
	}
}



