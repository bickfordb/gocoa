package gocoa

//#cgo LDFLAGS: -lobjc
//#include <CoreGraphics.h>
//#include <objc/objc-runtime.h>
import "C"
import (
	"bytes"
	"encoding/binary"
	"strconv"
	"unsafe"
)

/*
might be able to get away with omitting the Id() method by returning ids as byte arrays too
*/
type Passable interface {
	Id() C.id
	Bytes() []byte
	TypeString() string
}


/* Object **************************************************************************/

func (obj Object) Id() C.id           { return obj.idPointer() }
func (obj Object) Bytes() []byte      { return make([]byte, 0) }
func (obj Object) TypeString() string { return "@" }

/* Selector ************************************************************************/

func (sel Selector) Id() C.id           { return (C.id)(unsafe.Pointer(sel.selPointer())) }
func (sel Selector) Bytes() []byte      { return make([]byte, 0) }
func (sel Selector) TypeString() string { return ":" }

/* NSSize **************************************************************************/

// define
type NSSize struct {
	Width  float64
	Height float64
}

// create
func MakeNSSize(Width float64, Height float64) NSSize {
	return NSSize{Width, Height}
}

// implement Passable
func (nss NSSize) Id() C.id { return (C.id)(unsafe.Pointer(&(nss.Bytes()[0]))) }
func (nss NSSize) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &nss)
	return buf.Bytes()
}
func (nss NSSize) TypeString() string { return "{_NSSize=ff}" }

/* NSPoint *************************************************************************/

// define
type NSPoint struct {
	X float64
	Y float64
}

// create
func MakeNSPoint(X float64, Y float64) NSPoint {
	return NSPoint{X, Y}
}

// implement Passable
func (nsp NSPoint) Id() C.id { return (C.id)(unsafe.Pointer(&(nsp.Bytes()[0]))) }
func (nsp NSPoint) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &nsp)
	return buf.Bytes()
}
func (nsp NSPoint) TypeString() string { return "{_NSPoint=ff}" }

/* NSRect **************************************************************************/

// define
type NSRect struct {
	Origin NSPoint
	Size   NSSize
}

// create
func MakeNSRect(X float64, Y float64, Width float64, Height float64) NSRect {
	return NSRect{Origin: NSPoint{X, Y}, Size: NSSize{Width, Height}}
}

// implement Passable
func (nsr NSRect) Id() C.id { return (C.id)(unsafe.Pointer(&(nsr.Bytes()[0]))) }
func (nsr NSRect) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, &nsr)
	return buf.Bytes()
}
func (nsr NSRect) TypeString() string { return "{_NSRect={_NSPoint=ff}{_NSSize=ff}}" }

// convert to/from CGRect
func TypeNSRect(cgrect interface{}) NSRect {
	var result NSRect
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, cgrect)
	binary.Read(buf, binary.LittleEndian, &result)
	return result
}

func (nsr NSRect) CGRect() C.CGRect {
	var result C.CGRect
	buf := bytes.NewBuffer(nsr.Bytes())
	binary.Read(buf, binary.LittleEndian, result)
	return result
}

// pretty print
func (nsr NSRect) String() string {
	result := "["
	result = result + strconv.FormatFloat(nsr.Origin.X, 'f', -1, 64) + " "
	result = result + strconv.FormatFloat(nsr.Origin.Y, 'f', -1, 64) + " "
	result = result + strconv.FormatFloat(nsr.Size.Width, 'f', -1, 64) + " "
	result = result + strconv.FormatFloat(nsr.Size.Height, 'f', -1, 64) + "]"
	return result
}

/* NSUInteger **************************************************************************/

type NSUInteger uint64

// implement Passable
func (ui NSUInteger) Id() C.id { return (C.id)(unsafe.Pointer(&ui)) }
func (ui NSUInteger) Bytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ui)
	return buf.Bytes()
}
func (ui NSUInteger) TypeString() string { return "I" }

/* NSBoolean **************************************************************************/

type NSBoolean byte

// create
func MakeNSBoolean(value bool) NSBoolean {
	var result NSBoolean
	if value {
		result = 1
	}
	return result
}

// implement Passable
func (nsb NSBoolean) Id() C.id { return (C.id)(unsafe.Pointer(&nsb)) }
func (nsb NSBoolean) Bytes() []byte {
	result := make([]byte, 1)
	result[0] = (byte)(nsb)
	return result
}
func (nsb NSBoolean) TypeString() string { return "B" }

/* type strings ************************************************************************/

// XXX this really isn't quite the same as type safety...
func objcArgTypeString(argType string) string {

	switch argType {
	case "_Ctype_id":
		return "@"
	case "_Ctype_SEL":
		return ":"
	case "_Ctype_CGRect":
		return "{_NSRect={_NSPoint=ff}{_NSSize=ff}}"
	case "_Ctype_BOOL":
		return "B"
	default:
		return "@"
	}
	return ""
}

/*
#define _C_ID       '@'
#define _C_CLASS    '#'
#define _C_SEL      ':'
#define _C_CHR      'c'
#define _C_UCHR     'C'
#define _C_SHT      's'
#define _C_USHT     'S'
#define _C_INT      'i'
#define _C_UINT     'I'
#define _C_LNG      'l'
#define _C_ULNG     'L'
#define _C_LNG_LNG  'q'
#define _C_ULNG_LNG 'Q'
#define _C_FLT      'f'
#define _C_DBL      'd'
#define _C_BFLD     'b'
#define _C_BOOL     'B'
#define _C_VOID     'v'
#define _C_UNDEF    '?'
#define _C_PTR      '^'
#define _C_CHARPTR  '*'
#define _C_ATOM     '%'
#define _C_ARY_B    '['
#define _C_ARY_E    ']'
#define _C_UNION_B  '('
#define _C_UNION_E  ')'
#define _C_STRUCT_B '{'
#define _C_STRUCT_E '}'
#define _C_VECTOR   '!'
#define _C_CONST    'r'

*/
