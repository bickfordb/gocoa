package gocoa


import (
	"reflect"
	"fmt"
)


func DoSomething() {
	
	type testStruct struct {
		methodCount int
	}
	t := testStruct{0}
	
	AddMethod(&t)
//	t.Test()
	
	fmt.Println("end DoSomething")
}


func AddMethod(addToStruct interface{}) {
	
	v := reflect.ValueOf(addToStruct)
	
	if v.Kind() == reflect.Struct {
		
	/*	newMethod := reflect.Type.New(reflect.Method)
		newMethod.Name = "Test"
		newMethod.Func = func() { fmt.Println("t.Test() called successfully") }
		
//	    PkgPath string
//	    Type  Type  // method type
		
		lastMethod := v.Method(v.NumMethod())
		v.Append(lastMethod, newMethod)
		*/
	}
	
}