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
	"strconv"
)

type Reference struct {
	XMLName xml.Name `xml:"reference"`
	Key     string   `xml:"key,attr,omitempty"`
	Ref     string   `xml:"ref,attr,omitempty"`
}

type Bytes struct {
	XMLName xml.Name `xml:"bytes"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",innerxml"`
}

type Nil struct {
	XMLName xml.Name `xml:"nil"`
	Key     string   `xml:"key,attr"`
}

type Bool struct {
	XMLName xml.Name `xml:"bool"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",innerxml"`
}

type Boolean struct {
	XMLName xml.Name `xml:"boolean"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

type Float struct {
	XMLName xml.Name `xml:"float"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",innerxml"`
}

type Double struct {
	XMLName xml.Name `xml:"double"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",innerxml"`
}

type Int struct {
	XMLName xml.Name `xml:"int"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:",innerxml"`
}

type Integer struct {
	XMLName xml.Name `xml:"integer"`
	Key     string   `xml:"key,attr"`
	Value   string   `xml:"value,attr"`
}

type String struct {
	XMLName xml.Name `xml:"string"`
	Key     string   `xml:"key,attr"`
	Type 	string   `xml:"type,attr,omitempty"`
	Value   string   `xml:",innerxml"`
}

type Dictionary struct {
	XMLName      xml.Name     `xml:"dictionary"`
	Class        string       `xml:"class,attr"`
	Id           string       `xml:"id,attr,omitempty"`
	Key          string       `xml:"key,attr,omitempty"`
	Nils         []Nil        `xml:"nil,omitempty"`
	Byteses      []Bytes      `xml:"bytes,omitempty"`
	Bools        []Bool       `xml:"bool,omitempty"`
	Booleans     []Boolean    `xml:"boolean,omitempty"`
	Strings      []String     `xml:"string,omitempty"`
	Objects      []Object     `xml:"object,omitempty"`
	Arrays       []Array      `xml:"array,omitempty"`
	Ints         []Int        `xml:"int,omitempty"`
	Integers     []Integer    `xml:"integer,omitempty"`
	Floats       []Float      `xml:"float,omitempty"`
	Doubles      []Double     `xml:"double,omitempty"`
	References   []Reference  `xml:"reference,omitempty"`
	Dictionaries []Dictionary `xml:"dictionary,omitempty"`
}

type Array struct {
	XMLName      xml.Name     `xml:"array"`
	Class        string       `xml:"class,attr,omitempty"`
	Id           string       `xml:"id,attr,omitempty"`
	Key          string       `xml:"key,attr,omitempty"`
	Nils         []Nil        `xml:"nil,omitempty"`
	Byteses      []Bytes      `xml:"bytes,omitempty"`
	Bools        []Bool       `xml:"bool,omitempty"`
	Booleans     []Boolean    `xml:"boolean,omitempty"`
	Strings      []String     `xml:"string,omitempty"`
	Objects      []Object     `xml:"object,omitempty"`
	Arrays       []Array      `xml:"array,omitempty"`
	Ints         []Int        `xml:"int,omitempty"`
	Integers     []Integer    `xml:"integer,omitempty"`
	Floats       []Float      `xml:"float,omitempty"`
	Doubles      []Double     `xml:"double,omitempty"`
	References   []Reference  `xml:"reference,omitempty"`
	Dictionaries []Dictionary `xml:"dictionary,omitempty"`
}

type Object struct {
	XMLName      xml.Name     `xml:"object"`
	Class        string       `xml:"class,attr"`
	Id           string       `xml:"id,attr,omitempty"`
	Key          string       `xml:"key,attr,omitempty"`
	Nils         []Nil        `xml:"nil,omitempty"`
	Byteses      []Bytes      `xml:"bytes,omitempty"`
	Bools        []Bool       `xml:"bool,omitempty"`
	Booleans     []Boolean    `xml:"boolean,omitempty"`
	Strings      []String     `xml:"string,omitempty"`
	Objects      []Object     `xml:"object,omitempty"`
	Arrays       []Array      `xml:"array,omitempty"`
	Ints         []Int        `xml:"int,omitempty"`
	Integers     []Integer    `xml:"integer,omitempty"`
	Floats       []Float      `xml:"float,omitempty"`
	Doubles      []Double     `xml:"double,omitempty"`
	References   []Reference  `xml:"reference,omitempty"`
	Dictionaries []Dictionary `xml:"dictionary,omitempty"`
}

type Data struct {
	XMLName      xml.Name     `xml:"data"`
	Nils         []Nil        `xml:"nil,omitempty"`
	Byteses      []Bytes      `xml:"bytes,omitempty"`
	Bools        []Bool       `xml:"bool,omitempty"`
	Booleans     []Boolean    `xml:"boolean,omitempty"`
	Strings      []String     `xml:"string,omitempty"`
	Objects      []Object     `xml:"object,omitempty"`
	Arrays       []Array      `xml:"array,omitempty"`
	Ints         []Int        `xml:"int,omitempty"`
	Integers     []Integer    `xml:"integer,omitempty"`
	Floats       []Float      `xml:"float,omitempty"`
	Doubles      []Double     `xml:"double,omitempty"`
	References   []Reference  `xml:"reference,omitempty"`
	Dictionaries []Dictionary `xml:"dictionary,omitempty"`
}

type Archive struct {
	XMLName xml.Name `xml:"archive"`
	Type    string   `xml:"type,attr"`
	Version string   `xml:"version,attr"`
	Data    Data     `xml:"data"`
}

//	Data Object Array Dictionary, all searchable
type Searchable interface {
	GetArrays() []Array
	GetObjects() []Object
	GetDictionaries() []Dictionary
	MatchesConnectionId(connectionID string) (bool, *Object)
	ConnectionId() int
	MatchesClassAndKey(class string, key string) (bool, *Object)
	MatchesCustomObject(class string) (bool, *Object)
	MatchesClass(class string) (bool, *Object)
}

func (arr Array) GetArrays() []Array                                      { return arr.Arrays }
func (arr Array) GetObjects() []Object                                    { return arr.Objects }
func (arr Array) GetDictionaries() []Dictionary                           { return arr.Dictionaries }
func (arr Array) MatchesConnectionId(connectionID string) (bool, *Object) { return false, nil }
func (arr Array) ConnectionId() int                                       { return 0 }
func (arr Array) MatchesClassAndKey(class string, key string) (bool, *Object) { return false, nil }
func (arr Array) MatchesCustomObject(class string) (bool, *Object) { return false, nil }
func (arr Array) MatchesClass(class string) (bool, *Object) { return false, nil }

func (data Data) GetArrays() []Array                                      { return data.Arrays }
func (data Data) GetObjects() []Object                                    { return data.Objects }
func (data Data) GetDictionaries() []Dictionary                           { return data.Dictionaries }
func (data Data) MatchesConnectionId(connectionID string) (bool, *Object) { return false, nil }
func (data Data) ConnectionId() int                                       { return 0 }
func (data Data) MatchesClassAndKey(class string, key string) (bool, *Object) { return false, nil }
func (data Data) MatchesCustomObject(class string) (bool, *Object) { return false, nil }
func (data Data) MatchesClass(class string) (bool, *Object) { return false, nil }

func (dic Dictionary) GetArrays() []Array                                      { return dic.Arrays }
func (dic Dictionary) GetObjects() []Object                                    { return dic.Objects }
func (dic Dictionary) GetDictionaries() []Dictionary                           { return dic.Dictionaries }
func (dic Dictionary) MatchesConnectionId(connectionID string) (bool, *Object) { return false, nil }
func (dic Dictionary) ConnectionId() int                                       { return 0 }
func (dic Dictionary) MatchesClassAndKey(class string, key string) (bool, *Object) { return false, nil }
func (dic Dictionary) MatchesCustomObject(class string) (bool, *Object) { return false, nil }
func (dic Dictionary) MatchesClass(class string) (bool, *Object) { return false, nil }

func (obj Object) GetArrays() []Array            { return obj.Arrays }
func (obj Object) GetObjects() []Object          { return obj.Objects }
func (obj Object) GetDictionaries() []Dictionary { return obj.Dictionaries }
func (obj Object) MatchesConnectionId(connectionID string) (bool, *Object) {

	if len(obj.Ints) == 1 && obj.Ints[0].Key == "connectionID" && obj.Ints[0].Value == connectionID {
		//		fmt.Println("found:", obj.Class, obj.Ints[0].Value)
		return true, &obj
	}
	return false, nil
}

func (obj Object) ConnectionId() int {

	if len(obj.Ints) == 1 && obj.Ints[0].Key == "connectionID" {
		i, _ := strconv.ParseInt(obj.Ints[0].Value, 10, 32)
		return int(i)
	}
	return 0
}


func (obj Object) MatchesClassAndKey(class string, key string) (bool, *Object) {
	if obj.Class == class && obj.Key == key {
		return true, &obj
	}
	return false, nil
}


func (obj Object) MatchesClass(class string) (bool, *Object) {
	if obj.Class == class {
		return true, &obj
	}
	return false, nil
}

func (obj Object) MatchesCustomObject(class string) (bool, *Object) {
	if obj.Class == "NSCustomObject" {
		for i:=0; i<len(obj.Strings); i++ {
			if obj.Strings[i].Value == class {
				return true, &obj
			}
		}
	}
	return false, nil
}


func max(val1 int, val2 int) int {
	if val1 > val2 {
		return val1
	}
	return val2
}



func FindClass(dig Searchable, name string) *Object {

	// if type is object, check the connectionID and return self if found
	has, obj := dig.MatchesClass(name)
	if has {
		return obj
	}

	//	else, recurse into Objects, Arrays, and Dictionaries
	for i := 0; i < len(dig.GetObjects()); i++ {
		result := FindClass(dig.GetObjects()[i], name)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetArrays()); i++ {
		result := FindClass(dig.GetArrays()[i], name)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetDictionaries()); i++ {
		result := FindClass(dig.GetDictionaries()[i], name)
		if result != nil {
			return result
		}
	}
	return nil
}


func FindCustomObject(dig Searchable, name string) *Object {

	// if type is object, check the connectionID and return self if found
	has, obj := dig.MatchesCustomObject(name)
	if has {
		return obj
	}

	//	else, recurse into Objects, Arrays, and Dictionaries
	for i := 0; i < len(dig.GetObjects()); i++ {
		result := FindCustomObject(dig.GetObjects()[i], name)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetArrays()); i++ {
		result := FindCustomObject(dig.GetArrays()[i], name)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetDictionaries()); i++ {
		result := FindCustomObject(dig.GetDictionaries()[i], name)
		if result != nil {
			return result
		}
	}
	return nil
}


/*
* would be nice to get this at load with a SAX style interface
 */
func MaxValueForConnectionId(dig Searchable) int {

	// if type is object, check the connectionID and total max
	result := dig.ConnectionId()

	//	else, recurse into Objects, Arrays, and Dictionaries
	for i := 0; i < len(dig.GetObjects()); i++ {
		result = max(result, MaxValueForConnectionId(dig.GetObjects()[i]))
	}

	for i := 0; i < len(dig.GetArrays()); i++ {
		result = max(result, MaxValueForConnectionId(dig.GetArrays()[i]))
	}

	for i := 0; i < len(dig.GetDictionaries()); i++ {
		result = max(result, MaxValueForConnectionId(dig.GetDictionaries()[i]))
	}

	return result
}


func FindClassAndKey(dig Searchable, class string, key string) *Object {
	
	// if type is object, check the connectionID and return self if found
	has, obj := dig.MatchesClassAndKey(class, key)
	if has {
		return obj
	}

	//	else, recurse into Objects, Arrays, and Dictionaries
	for i := 0; i < len(dig.GetObjects()); i++ {
		result := FindClassAndKey(dig.GetObjects()[i], class, key)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetArrays()); i++ {
		result := FindClassAndKey(dig.GetArrays()[i], class, key)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetDictionaries()); i++ {
		result := FindClassAndKey(dig.GetDictionaries()[i], class, key)
		if result != nil {
			return result
		}
	}
	return nil
}


func FindObjectByConnectionId(dig Searchable, connectionID string) *Object {

	// if type is object, check the connectionID and return self if found
	has, obj := dig.MatchesConnectionId(connectionID)
	if has {
		return obj
	}

	//	else, recurse into Objects, Arrays, and Dictionaries
	for i := 0; i < len(dig.GetObjects()); i++ {
		result := FindObjectByConnectionId(dig.GetObjects()[i], connectionID)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetArrays()); i++ {
		result := FindObjectByConnectionId(dig.GetArrays()[i], connectionID)
		if result != nil {
			return result
		}
	}

	for i := 0; i < len(dig.GetDictionaries()); i++ {
		result := FindObjectByConnectionId(dig.GetDictionaries()[i], connectionID)
		if result != nil {
			return result
		}
	}

	return nil
}

/*
<object class="IBObjectContainer" key="IBDocument.Objects">
	<array class="NSMutableArray" key="connectionRecords">

		<object class="IBConnectionRecord">
			<object class="IBActionConnection" key="connection">
				<string key="label">orderFrontStandardAboutPanel:</string>
				<reference key="source" ref="1021"/>
				<reference key="destination" ref="238522557"/>
			</object>
			<int key="connectionID">142</int>
		</object>

		<object class="IBConnectionRecord">
			<object class="IBOutletConnection" key="connection">
				<string key="label">delegate</string>
				<reference key="source" ref="879354517"/>
				<reference key="destination" ref="889664573"/>
			</object>
			<int key="connectionID">884</int>
		</object>

*/

const (
	IBOutletConnection = "IBOutletConnection"
	IBActionConnection = "IBActionConnection"
)

func CreateIBConnection(class string, label string, source int, destination int, connectionID int) Object {

	references := make([]Reference, 2)
	references[0] = Reference{Key: "source", Ref: strconv.FormatInt(int64(source), 10)}
	references[1] = Reference{Key: "destination", Ref: strconv.FormatInt(int64(destination), 10)}

	strings := make([]String, 1)
	strings[0] = String{Key: "label", Value: label}

	ints := make([]Int, 1)
	ints[0] = Int{Key: "connectionID", Value: strconv.FormatInt(int64(connectionID), 10)}

	objects := make([]Object, 1)
	objects[0] = Object{Class: class, Key: "connection", Strings: strings, References: references}

	return Object{Class: "IBConnectionRecord", Objects: objects, Ints: ints}
}


func (obj *Object) ArrayForClassAndKey(class string, key string) *Array {
	for i:=0; i< len(obj.Arrays); i++ {
		if obj.Arrays[i].Key == key && obj.Arrays[i].Class == class {
			return &obj.Arrays[i]
		}
	}
	return nil
}


func (obj *Object) IntForKey(name string) *Int {
	for i:=0; i< len(obj.Ints); i++ {
		if obj.Ints[i].Key == name {
			return &obj.Ints[i]
		}
	}
	return nil
}


func AddIBConnection(v Searchable, connection Object, connectionID int) {
	objectContainer := FindClassAndKey(v, "IBObjectContainer", "IBDocument.Objects")
	intMaxID := objectContainer.IntForKey("maxID")
	connectionRecordsArray := objectContainer.ArrayForClassAndKey("NSMutableArray", "connectionRecords")
	connectionRecordsArray.Objects = append(connectionRecordsArray.Objects, connection)
	i, _ := strconv.ParseInt(intMaxID.Value, 10, 32)
	i32 := max(int(i), connectionID)
	intMaxID.Value = strconv.FormatInt(int64(i32), 10)
}


func main() {

	if len(os.Args) < 3 {
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
	
	// locate the id of the ApplicationController
/*	
	<object class="NSCustomObject" id="889664573">
		<string key="NSClassName">ApplicationController</string>
	</object>*/
	
	appController := FindCustomObject(v.Data, "ApplicationController")
	appControllerId, _ := strconv.ParseInt(appController.Id, 10, 32)
	
	app := FindCustomObject(v.Data, "NSApplication")
	appId, _ := strconv.ParseInt(app.Id, 10, 32)
	
	scrollView := FindClass(v.Data, "NSScrollView")
	scrollViewId, _ := strconv.ParseInt(scrollView.Id, 10, 32)
	
	windowTemplate := FindClass(v.Data, "NSWindowTemplate")
	windowTemplateId, _ := strconv.ParseInt(windowTemplate.Id, 10, 32)
	
	/*
	421919869 - ApplicationController
	
	<object class="IBConnectionRecord">
		<object class="IBOutletConnection" key="connection">
			<string key="label">delegate</string>
			<reference key="source" ref="426487645"/>			'NSApplication - 1'
			<reference key="destination" ref="421919869"/>		'ApplicationController'
		</object>
		<int key="connectionID">255</int>
	</object>
	
	<object class="IBConnectionRecord">
		<object class="IBOutletConnection" key="connection">
			<string key="label">delegate</string>
			<reference key="source" ref="852095695"/>			'NSWindowTemplate'
			<reference key="destination" ref="421919869"/>		'ApplicationController'
		</object>
		<int key="connectionID">256</int>
	</object>
	
	<object class="IBConnectionRecord">
		<object class="IBOutletConnection" key="connection">
			<string key="label">textBox1</string>
			<reference key="source" ref="421919869"/>			'ApplicationController'
			<reference key="destination" ref="388213687"/>		'NSTextField'
		</object>
		<int key="connectionID">252</int>
	</object>
	
	<object class="IBConnectionRecord">
		<object class="IBOutletConnection" key="connection">
			<string key="label">mainWindow</string>
			<reference key="source" ref="421919869"/>			'ApplicationController'
			<reference key="destination" ref="852095695"/>		'NSWindowTemplate'
		</object>
		<int key="connectionID">253</int>
	</object>
	
	<object class="IBConnectionRecord">
		<object class="IBActionConnection" key="connection">
			<string key="label">buttonClick:</string>
			<reference key="source" ref="421919869"/>			'ApplicationController'
			<reference key="destination" ref="697635292"/>		'NSButton'
		</object>
		<int key="connectionID">257</int>
	</object>
	
	*/
	
	
	// add the outlets
	maxId := MaxValueForConnectionId(v.Data) + 1
	connection := CreateIBConnection(IBOutletConnection, "scrollTable1", int(appControllerId), int(scrollViewId), maxId)
	AddIBConnection(v.Data, connection, maxId)
	maxId++
	
	connection = CreateIBConnection(IBOutletConnection, "delegate", int(appId), int(appControllerId), maxId)
	AddIBConnection(v.Data, connection, maxId)
	maxId++
	
	connection = CreateIBConnection(IBOutletConnection, "delegate", int(windowTemplateId), int(appControllerId), maxId)
	AddIBConnection(v.Data, connection, maxId)
	maxId++
	
	connection = CreateIBConnection(IBOutletConnection, "mainWindow", int(appControllerId), int(windowTemplateId), maxId)
	AddIBConnection(v.Data, connection, maxId)
	maxId++
	
//	connection = CreateIBConnection(IBActionConnection, "applicationWillFinishLaunching:", int(appId), int(appControllerId), maxId)
//	AddIBConnection(v.Data, connection, maxId)
	
	
	// ok, output the difference
	output, err := xml.MarshalIndent(&v, "", "\t")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}

	fmt.Print(xml.Header)
	os.Stdout.Write(output)

}
