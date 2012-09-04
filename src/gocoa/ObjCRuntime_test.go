package gocoa

import (
	"testing"
//	"unsafe"
)
	
/*
* If it has a test, it's an interface. Most of these are naive tests, when I start discovering
* more interesting failure scenarios, they'll be incorporated.
*/


func Test_ClassForName(t *testing.T) {
	nso := ClassForName("NSObject").Instance("alloc").Call("init")
	t.Log("ClassForName(\"NSObject\").Instance(\"alloc\").Call(\"init\")", nso)
	if nso == 0 {
		t.Fail()
	}
}

func Test_SelectorForName(t *testing.T) {
	sel := SelectorForName("alloc")
	t.Log("SelectorForName(\"alloc\")", sel)
	if sel == 0 {
		t.Fail()
	}
}


func Test_Class_RespondsTo(t *testing.T) {
	nso := ClassForName("NSObject").Instance("alloc").Call("init")
	if !nso.Class().RespondsTo("classForPortCoder") {
		t.Fail()
	}
	if nso.Class().RespondsTo("bogusMethod:withBogusArgs:") {
		t.Fail()
	}
}

func Test_Class_Name(t *testing.T) {
	name := ClassForName("NSObject").Name()
	if name != "NSObject" {
		t.Fail()
	}
}

func Test_Class_Super(t *testing.T) {
	name := ClassForName("NSObject").Super().Name()
	if len(name) == 0 {
		t.Fail()
	}
	if name != "Object" {
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
	t.Fail()
}

func Test_Class_Property(t *testing.T) {
	t.Fail()
}

func Test_Class_ListInstanceVariables(t *testing.T) {
	t.Fail()
}

func Test_Class_ListMethods(t *testing.T) {
	t.Fail()
}

func Test_Class_ListProperties(t *testing.T) {
	t.Fail()
}

func Test_Class_Subclass(t *testing.T) {
	t.Fail()
}

func Test_Class_AddMethod(t *testing.T) {
	t.Fail()
}

func Test_Class_AddIvar(t *testing.T) {
	t.Fail()
}

func Test_Class_Register(t *testing.T) {
	t.Fail()
}



func Test_Object(t *testing.T) {
	t.Fail()
}

func Test_Object_Method(t *testing.T) {
	t.Fail()
}

func Test_Object_Class(t *testing.T) {
	t.Fail()
}

func Test_Object_InstanceVariable(t *testing.T) {
	t.Fail()
}

func Test_Object_SetInstanceVariable(t *testing.T) {
	t.Fail()
}


func Object_I_Instance(cls Class, method string, args ...Passable) Object {
	return ((Object)(cls)).I(method, args...)
}

// XXX this should match Object.Call, isolate and figure out the problem
func Test_Object_I(t *testing.T) {
	
	
	view := ClassForName("NSView").Subclass("SimpleView")
	view.Register()
	
	controller := ClassForName("NSObject").Subclass("Controller")
	controller.Register()
	
	controllerInst := Object_I_Instance(ClassForName("Controller"), "alloc").I("init")
	t.Log("alloc", controllerInst)
/*	
	bundle := Object_I_Instance(ClassForName("NSBundle"), "alloc")
	dict := NSDictionary("NSOwner", app)
	
	bundle = bundle.I("initWithPath:", NSString(".."))
	bundle.I("loadNibFile:externalNameTable:withZone:", NSString("SimpleView"), dict, app.I("zone"))

//

	bundle = Object_I_Instance(ClassForName("NSBundle"), "alloc")
	bundle = bundle.I("initWithPath:", NSString("."))
	if bundle == 0 { 
		t.Fail() 
	} else {
		t.Log("bundle", bundle)
	}

	someArray := Object_I_Instance(ClassForName("NSMutableArray"), "alloc").I("init")
	someCount := (NSUInteger)(someArray.I("count"))
	if someCount != 0 {
		t.Fail()
	} else {
		t.Log("someCount", someCount)
	}
	
	nilresult := bundle.I("initWithPath:", NSString("A String"), MakeNSRect(0,0,0,0), someArray, someCount, MakeNSBoolean(true))
	if nilresult != 0 {
		t.Fail()
	} else {
		t.Log("nilresult", bundle)
	}*/

}

func Test_Object_Call(t *testing.T) {
	t.Fail()
} 

func Test_Object_CallSuper(t *testing.T) {
	t.Fail()
}




func Test_Selector_Name(t *testing.T) {
	t.Fail()
}



func Test_Method(t *testing.T) {
	t.Fail()
}

func Test_Method_ArgumentCount(t *testing.T) {
	t.Fail()
}

func Test_Method_ArgumentType(t *testing.T) {
	t.Fail()
}

func Test_Method_Name(t *testing.T) {
	t.Fail()
}



func Test_Ivar(t *testing.T) {
	t.Fail()
}

func Test_Ivar_Name(t *testing.T) {
	t.Fail()
}



func Test_Property(t *testing.T) {
	t.Fail()
}

func Test_Property_Name(t *testing.T) {
	t.Fail()
}

func Test_Property_Attributes(t *testing.T) {
	t.Fail()
}





