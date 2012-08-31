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
	"strings"
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
example from HelloWorld.nib
421919869 - ApplicationController

<object class="IBObjectContainer" key="IBDocument.Objects">
	<array class="NSMutableArray" key="connectionRecords">...

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


func usage(err string) {
	if len(err) > 0 {
		fmt.Println("nibble:", err)
	}
	usage := `
usage: nibble [-a (options)] [-d (options)] [-l class] [input file]
	
        -a outlet|action name source dest
            add outlet or action 'name', connect 'source' to 'dest' 
        -d class
            delete all outlets and actions for 'class'
        -l class
            list all objects of type 'class'
			
example: nibble -a outlet mainWindow ApplicationController NSWindowTemplate file.nib
            adds an outlet, 'mainWindow' from ApplicationController to 
            NSWindowTemplate in file.nib and prints the output
`
	fmt.Println(usage)
}

func loadNib(fileName string) Archive {

	result := Archive{Type: "none", Version: "0"}
	file, err := os.Open(fileName)
	decoder := xml.NewDecoder(file)

	err = decoder.Decode(&result)
	if err != nil {
		fmt.Println("error:", err, " fileName:", fileName)
		return result
	}

	file.Close()
	
	return result
}

func printNib(archive Archive) {

	output, err := xml.MarshalIndent(&archive, "", "\t")
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
	fmt.Print(xml.Header)
	os.Stdout.Write(output)
}


func main() {

	if len(os.Args) < 3 || os.Args[1] == "-help" {
		usage("")
		return
	}
	
	switch(os.Args[1]) {
	
		case "-a":
			if len(os.Args) == 7 {
				if os.Args[2] == "outlet" {
					v := loadNib(os.Args[6])
					AddOutlet(v.Data, os.Args[3], os.Args[4], os.Args[5])
					printNib(v)
				} else if os.Args[2] == "action" {
					v := loadNib(os.Args[6])
					AddAction(v.Data, os.Args[3], os.Args[4], os.Args[5])
					printNib(v)
				
/*				} else if os.Args[2] == "appdelegate" {
					v := loadNib(os.Args[6])
					AddAppDelegate(v.Data, os.Args[3], os.Args[4], os.Args[5])
					printNib(v)
*/					
				} else {
					usage(os.Args[2] + ", invalid IBConnection type")
				}
			} else {
				usage("invalid args for -a:" + strings.Join(os.Args, " "))
			}
			break
		case "-d":
			if len(os.Args) == 4 {
				v := loadNib(os.Args[3])
				DeleteOutlets(v.Data, os.Args[2])
				printNib(v)
			} else {
				usage("invalid args for -d:" + strings.Join(os.Args, " "))
			}
			break
		case "-l":
			if len(os.Args) == 4 {
				v := loadNib(os.Args[3])
				ListClasses(v.Data, os.Args[2])
				
			} else {
				usage("invalid args for -l:" + strings.Join(os.Args, " "))
			}
			break
		
		default: 
			usage("invalid option:" + os.Args[1])
	}
	
	
	// the source of an outlet is the custom object, slightly counterintuitive
/*	AddOutlet(v.Data, "scrollTable1", "ApplicationController", "NSScrollView")
	AddOutlet(v.Data, "mainWindow", "ApplicationController", "NSWindowTemplate")
	
	AddAction(v.Data, "applicationWillFinishLaunching:", "ApplicationController", "NSApplication")*/	
}



func ListClasses(v Searchable, name string) {


}


func DeleteOutlets(v Searchable, class string) {

}


func AddOutlet(v Searchable, name string, classSource string, customObjectDestination string) {
	
	source := FindCustomObject(v, classSource)
	sourceId, _ := strconv.ParseInt(source.Id, 10, 32)
	
	destination := FindClass(v, customObjectDestination)
	destinationId, _ := strconv.ParseInt(destination.Id, 10, 32)
	
	maxId := MaxValueForConnectionId(v) + 1
	
	connection := CreateIBConnection(IBOutletConnection, name, int(sourceId), int(destinationId), maxId)
	AddIBConnection(v, connection, maxId)
	maxId++
	
	connection = CreateIBConnection(IBOutletConnection, "delegate", int(destinationId), int(sourceId), maxId)
	AddIBConnection(v, connection, maxId)
	
}


func AddAction(v Searchable, action string, customObjectSource string, customObjectDestination string) {
	
	source := FindCustomObject(v, customObjectSource)
	sourceId, _ := strconv.ParseInt(source.Id, 10, 32)
	
	destination := FindCustomObject(v, customObjectDestination)
	destinationId, _ := strconv.ParseInt(destination.Id, 10, 32)
	
	maxId := MaxValueForConnectionId(v) + 1
		
	connection := CreateIBConnection(IBActionConnection, action, int(destinationId), int(sourceId), maxId)
	AddIBConnection(v, connection, maxId)
}

/*
* something else, adding an application delegate
*/
func AddAppDelegate(v Searchable, action string, customObjectSource string, customObjectDestination string) {
	
	source := FindCustomObject(v, customObjectSource)
	sourceId, _ := strconv.ParseInt(source.Id, 10, 32)
	
	destination := FindCustomObject(v, "NSApplication")
	destinationId, _ := strconv.ParseInt(destination.Id, 10, 32)
	
	maxId := MaxValueForConnectionId(v) + 1
		
	connection := CreateIBConnection(IBActionConnection, "delegate", int(destinationId), int(sourceId), maxId)
	AddIBConnection(v, connection, maxId)
}
