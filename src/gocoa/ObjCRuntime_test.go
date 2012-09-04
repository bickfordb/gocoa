package gocoa

import (
//	"runtime"
	"testing"
//	"unsafe"
)
	
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

func Test_Class_Method(t *testing.T) {
	class := ClassForName("NSObject")
	method := class.Method("alloc")
	if method == 0 {
		t.Log("method", method)
		t.Fail()
	}
}

func Test_Class_Property(t *testing.T) {
	t.Log("where's a good system class with predefined properties?")
	t.Fail()
/*	class := ClassForName("NSDate")
	class.ListProperties()
	property := class.Property("propname")	
	if property == 0 {
		t.Fail()
	} else {
		t.Log("property", property.Name())
	}*/
}

func Test_Class_InstanceVariables(t *testing.T) {
	t.Log("unimplemented")
	t.Fail()
}

func Test_Class_Methods(t *testing.T) {
	t.Log("unimplemented")
	t.Fail()
}

func Test_Class_Properties(t *testing.T) {
	t.Log("unimplemented")
	t.Fail()
}


func Test_Class_Subclass(t *testing.T) {
	subclass := ClassForName("NSObject").Subclass("NSSubclass")
	if subclass == 0 {
		t.Fail()
	}
	if subclass.Name() != "NSSubclass" {
		t.Fail()
	}
	
	instance := subclass.Instance("alloc")
	if instance == 0 {
		t.Fail()
	}
	if instance.Class().Name() != "NSSubclass" {
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
	subclass.AddIvar("someIvar", dude.Class())
//	causes runtime error:
//	subclass.Register()
	instance := subclass.Instance("alloc").Call("init")
	instance.SetInstanceVariable("someIvar", dude)
	ivar := instance.InstanceVariable("someIvar")
	if ivar == 0 {
		t.Log("ivar", ivar)
		t.Fail()
	}
}



// XXX this isn't getting at the meat of the issue
func Object_I_Instance(cls Class, method string, args ...Passable) Object {
	return ((Object)(cls)).I(method, args...)
}
func Test_Object_I(t *testing.T) {

	bundle := Object_I_Instance(ClassForName("NSBundle"), "alloc")
	bundle = bundle.I("initWithPath:", NSString("."))
	if bundle == 0 {
		t.Log("bundle", bundle)
		t.Fail() 
	}

	someArray := Object_I_Instance(ClassForName("NSMutableArray"), "alloc").I("init")
	if someArray == 0 {
		t.Log("someArray", someArray)
		t.Fail()
	}
	
	someCount := (NSUInteger)(someArray.I("count"))
	if someCount != 0 {
		t.Log("someCount", someCount)
		t.Fail()
	}
	
	nilresult := bundle.I("initWithPath:", NSString("A String"), MakeNSRect(0,0,0,0), someArray, someCount, MakeNSBoolean(true))
	if nilresult != 0 {
		t.Log("nilresult", nilresult)
		t.Fail()
	}

}

func Test_Object_Call(t *testing.T) {
	
	bundle := ClassForName("NSBundle").Instance("alloc").Call("initWithPath:", NSString("."))
	if bundle == 0 {
		t.Log("bundle", bundle)
		t.Fail() 
	}

	someArray := ClassForName("NSMutableArray").Instance("alloc").Call("init")
	if someArray == 0 {
		t.Log("someArray", someArray)
		t.Fail()
	}
	
	someCount := (NSUInteger)(someArray.Call("count"))
	if someCount != 0 {
		t.Log("someCount", someCount)
		t.Fail()
	}
	/*
	// hangs
	nilresult := bundle.Call("initWithPath:", NSString("A String"), Object(someArray), Object(someCount), Object(MakeNSBoolean(true)))
	if nilresult != 0 {
		t.Fail()
	} else {
		t.Log("nilresult", nilresult)
	}
	*/
} 

func Test_Object_CallSuper(t *testing.T) {
	t.Log("unimplemented")
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
	t.Log("unimplemented, wants class.Methods")
	t.Fail()
}

func Test_Method_ArgumentType(t *testing.T) {
	t.Log("unimplemented, wants class.Methods")
	t.Fail()
}

func Test_Method_Name(t *testing.T) {
	t.Log("unimplemented, wants class.Methods")
	t.Fail()
}



func Test_Ivar_Name(t *testing.T) {
	t.Log("unimplemented, wants class.Ivars")
	t.Fail()
}



func Test_Property_Name(t *testing.T) {
	t.Log("unimplemented, wants class.Properties")
	t.Fail()
}

func Test_Property_Attributes(t *testing.T) {
	t.Log("unimplemented, wants class.Properties")
	t.Fail()
}



