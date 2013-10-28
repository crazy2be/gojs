package gojs

// #include <stdlib.h>
// #include <JavaScriptCore/JSStringRef.h>
import "C"
import (
	"syscall"
	"unsafe"
)

//=========================================================
// StringRef
//

type String struct {
}

func NewString(value string) *String {
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	ref := C.JSStringCreateWithUTF8CString(cvalue)
	return (*String)(unsafe.Pointer(ref))
}

func newStringFromRef(ref C.JSStringRef) *String {
	return (*String)(unsafe.Pointer(ref))
}

func (ref *String) Retain() {
	C.JSStringRetain((C.JSStringRef)(unsafe.Pointer(ref)))
}

func (ref *String) Release() {
	C.JSStringRelease((C.JSStringRef)(unsafe.Pointer(ref)))
}

func (ref *String) String() string {
	// Conversion 1, null-terminate UTF-8 string
	len := C.JSStringGetMaximumUTF8CStringSize((C.JSStringRef)(unsafe.Pointer(ref)))
	buffer := C.malloc(len)
	if buffer == nil {
		panic(syscall.ENOMEM)
	}
	defer C.free(buffer)
	C.JSStringGetUTF8CString((C.JSStringRef)(unsafe.Pointer(ref)), (*C.char)(buffer), len)

	// Conversion 2, Go string
	ret := C.GoString((*C.char)(buffer))
	return ret
}

// Bytes returns a byte slice with the bytes of this string.
func (ref *String) Bytes() []byte {
	// Conversion 1, null-terminate UTF-8 string
	len := C.JSStringGetMaximumUTF8CStringSize((C.JSStringRef)(unsafe.Pointer(ref)))
	buffer := C.malloc(len)
	if buffer == nil {
		panic(syscall.ENOMEM)
	}
	defer C.free(buffer)
	strlenWithNul := C.JSStringGetUTF8CString((C.JSStringRef)(unsafe.Pointer(ref)), (*C.char)(buffer), len)

	// Conversion 2, Go []byte
	return C.GoBytes(buffer, C.int(strlenWithNul-1))
}

func (ref *String) Length() uint32 {
	ret := C.JSStringGetLength((C.JSStringRef)(unsafe.Pointer(ref)))
	return uint32(ret)
}

func (ref *String) Equal(rhs *String) bool {
	ret := C.JSStringIsEqual((C.JSStringRef)(unsafe.Pointer(ref)), (C.JSStringRef)(unsafe.Pointer(rhs)))
	return bool(ret)
}

func (ref *String) EqualToString(rhs string) bool {
	crhs := C.CString(rhs)
	defer C.free(unsafe.Pointer(crhs))
	ret := C.JSStringIsEqualToUTF8CString((C.JSStringRef)(unsafe.Pointer(ref)), crhs)
	return bool(ret)
}
