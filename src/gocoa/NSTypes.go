package gocoa

//#include <CoreGraphics.h>
import "C"


func NSMakeRect(X float64, Y float64, Width float64, Height float64) NSRect {
	return NSRect{ Origin: NSPoint{ X, Y }, Size: NSSize{ Width, Height } }
}

type NSRect struct {
	Origin	NSPoint
	Size	NSSize
}

type NSSize struct {
	Width	float64
	Height	float64
}

type NSPoint struct{
	X	float64
	Y	float64
}

func ToNSRect(cgrect C.CGRect) *NSRect {
	var result NSRect

	return &result
}

func (nsr *NSRect) Bytes() []byte {
	var result []byte
	
	return result
}

func (nsr *NSRect) CGRect() C.CGRect {
	var result C.CGRect
	
	return result
}


type NSUInteger uint64



// XXX somewhat incomplete 
func objcArgTypeString(argType string) string {
	
	switch argType {
		case "_Ctype_id":		return "@"
		case "_Ctype_SEL":		return ":"
		case "_Ctype_CGRect":	return "{_NSRect={_NSPoint=ff}{_NSSize=ff}}"
		case "_Ctype_BOOL":		return "B" 
		default: 				return "@"
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

