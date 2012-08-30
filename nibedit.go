package main

/*
* Adding outlets and actions to .xib files is something that you'd want to be able to do without first 
* generating objc headers for classes that were conceived of in some other language. Beginning work on
* a tool that roundtrips xib files through Go's parser.
*/

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Reference struct {
	XMLName	xml.Name	`xml:"reference"`
	Key		string		`xml:"key,attr,omitempty"`
	Ref		string		`xml:"ref,attr,omitempty"`
}

type Bytes struct {
	XMLName	xml.Name	`xml:"bytes"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}

type Nil struct {
	XMLName	xml.Name	`xml:"nil"`
	Key		string		`xml:"key,attr"`
}

type Bool struct {
	XMLName	xml.Name	`xml:"bool"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}

type Boolean struct {
	XMLName	xml.Name	`xml:"boolean"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:"value,attr"`
}

type Float struct {
	XMLName	xml.Name	`xml:"float"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}

type Double struct {
	XMLName	xml.Name	`xml:"double"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}

type Int struct {
	XMLName	xml.Name	`xml:"int"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}

type Integer struct {
	XMLName	xml.Name	`xml:"integer"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:"value,attr"`
}
	
type String struct {
	XMLName	xml.Name	`xml:"string"`
	Key		string		`xml:"key,attr"`
	Value	string		`xml:",innerxml"`
}


type Dictionary struct {
	XMLName			xml.Name		`xml:"dictionary"`
	Class			string			`xml:"class,attr"`
	Id				string			`xml:"id,attr,omitempty"`
	Key				string			`xml:"key,attr,omitempty"`
	
	Nils			[]Nil			`xml:"nil,omitempty"`
	Byteses			[]Bytes			`xml:"bytes,omitempty"`
	Bools			[]Bool			`xml:"bool,omitempty"`
	Booleans		[]Boolean		`xml:"boolean,omitempty"`
	Strings			[]String		`xml:"string,omitempty"`
	Objects			[]Object		`xml:"object,omitempty"`
	Arrays			[]Array			`xml:"array,omitempty"`
	Ints			[]Int			`xml:"int,omitempty"`
	Integers		[]Integer		`xml:"integer,omitempty"`
	Floats			[]Float			`xml:"float,omitempty"`
	Doubles			[]Double		`xml:"double,omitempty"`
	References		[]Reference		`xml:"reference,omitempty"`
	Dictionaries	[]Dictionary	`xml:"dictionary,omitempty"`
}

type Array struct {
	XMLName			xml.Name		`xml:"array"`
	Class			string			`xml:"class,attr,omitempty"`
	Id				string			`xml:"id,attr,omitempty"`
	Key				string			`xml:"key,attr,omitempty"`
	
	Nils			[]Nil			`xml:"nil,omitempty"`
	Byteses			[]Bytes			`xml:"bytes,omitempty"`
	Bools			[]Bool			`xml:"bool,omitempty"`
	Booleans		[]Boolean		`xml:"boolean,omitempty"`
	Strings			[]String		`xml:"string,omitempty"`
	Objects			[]Object		`xml:"object,omitempty"`
	Arrays			[]Array			`xml:"array,omitempty"`
	Ints			[]Int			`xml:"int,omitempty"`
	Integers		[]Integer		`xml:"integer,omitempty"`
	Floats			[]Float			`xml:"float,omitempty"`
	Doubles			[]Double		`xml:"double,omitempty"`
	References		[]Reference		`xml:"reference,omitempty"`
	Dictionaries	[]Dictionary	`xml:"dictionary,omitempty"`
}
	
type Object struct {
	XMLName			xml.Name		`xml:"object"`
	Class			string			`xml:"class,attr"`
	Id				string			`xml:"id,attr,omitempty"`
	Key				string			`xml:"key,attr,omitempty"`
	
	Nils			[]Nil			`xml:"nil,omitempty"`
	Byteses			[]Bytes			`xml:"bytes,omitempty"`
	Bools			[]Bool			`xml:"bool,omitempty"`
	Booleans		[]Boolean		`xml:"boolean,omitempty"`
	Strings			[]String		`xml:"string,omitempty"`
	Objects			[]Object		`xml:"object,omitempty"`
	Arrays			[]Array			`xml:"array,omitempty"`
	Ints			[]Int			`xml:"int,omitempty"`
	Integers		[]Integer		`xml:"integer,omitempty"`
	Floats			[]Float			`xml:"float,omitempty"`
	Doubles			[]Double		`xml:"double,omitempty"`
	References		[]Reference		`xml:"reference,omitempty"`
	Dictionaries	[]Dictionary	`xml:"dictionary,omitempty"`
}
	
type Data struct {
	XMLName			xml.Name		`xml:"data"`
	
	Nils			[]Nil			`xml:"nil,omitempty"`
	Byteses			[]Bytes			`xml:"bytes,omitempty"`
	Bools			[]Bool			`xml:"bool,omitempty"`
	Booleans		[]Boolean		`xml:"boolean,omitempty"`
	Strings			[]String		`xml:"string,omitempty"`
	Objects			[]Object		`xml:"object,omitempty"`
	Arrays			[]Array			`xml:"array,omitempty"`
	Ints			[]Int			`xml:"int,omitempty"`
	Integers		[]Integer		`xml:"integer,omitempty"`
	Floats			[]Float			`xml:"float,omitempty"`
	Doubles			[]Double		`xml:"double,omitempty"`
	References		[]Reference		`xml:"reference,omitempty"`
	Dictionaries	[]Dictionary	`xml:"dictionary,omitempty"`
}
	
type Archive struct {
	XMLName		xml.Name		`xml:"archive"`
	Type		string			`xml:"type,attr"`
	Version		string			`xml:"version,attr"`
	Data		Data			`xml:"data"`
}


func main() {
	
//	fmt.Println(os.Args)
	
	if len(os.Args)<3 {
		fmt.Println("usage\n\tnibedit <inputfile> <outputfile>")
		return
	}
	
	v := Archive{Type: "none", Version: "0"}
	fileName := os.Args[1]
	file, err := os.Open(fileName)
	decoder := xml.NewDecoder(file)
	
	err = decoder.Decode(&v)
	if err != nil {
	    fmt.Printf("error: %v", err)
	    return
	}
	
	file.Close()
	
//	fmt.Printf("XMLName: %#v\n", v.XMLName)
//	fmt.Printf("Type: %#v\n", v.Type)
//	fmt.Printf("Version: %#v\n", v.Version)	
	
	// okay diff the output..
	output, err := xml.MarshalIndent(v, "", "\t")
	if err != nil {
	    fmt.Printf("error: %v\n", err)
		return
	}
	
	fmt.Print(xml.Header)
	os.Stdout.Write(output)
	
	
/*	for i:=0; i<len(v.Strings); i++ {
		fmt.Printf("%v: %v\n", v.Strings[i].Key, v.Strings[i].Value)
	}*/
}

