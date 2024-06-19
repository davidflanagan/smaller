package smaller

import (
	"encoding/binary"
	"reflect"
	"sync/atomic"
	"unsafe"
)

type String struct {
	p unsafe.Pointer
}

const MaxCapacity = 1<<17 - 4

func NewString(s string) String {
	var (
		length    = len(s)
		prefixLen = 0
		totalLen  int
		numWords  int
		up        unsafe.Pointer
		bytes     []byte
	)

	switch {
	case length < 1<<7:
		prefixLen = 1
	case length < 1<<14:
		prefixLen = 2
	case length < 1<<21:
		prefixLen = 3
	case length < 1<<28:
		prefixLen = 4
	case length < 1<<35:
		prefixLen = 5
	case length < 1<<42:
		prefixLen = 6
	case length < 1<<49:
		prefixLen = 7
	case length < 1<<56:
		prefixLen = 8
	default:
		panic("string is too long")
	}

	totalLen = length + prefixLen        // This is how many bytes we need
	numWords = ((totalLen - 1) >> 3) + 1 // Which rounds up to this many uint64s

	switch numWords {
	case 1:
		up = unsafe.Pointer(new([1]uint64))
	case 2:
		up = unsafe.Pointer(new([2]uint64))
	case 3:
		up = unsafe.Pointer(new([3]uint64))
	case 4:
		up = unsafe.Pointer(new([4]uint64))
	case 5:
		up = unsafe.Pointer(new([5]uint64))
	case 6:
		up = unsafe.Pointer(new([6]uint64))
	case 7:
		up = unsafe.Pointer(new([7]uint64))
	case 8:
		up = unsafe.Pointer(new([8]uint64))
	case 9:
		up = unsafe.Pointer(new([9]uint64))
	case 10:
		up = unsafe.Pointer(new([10]uint64))
	case 11:
		up = unsafe.Pointer(new([11]uint64))
	case 12:
		up = unsafe.Pointer(new([12]uint64))
	case 16:
		up = unsafe.Pointer(new([16]uint64))
	case 32:
		up = unsafe.Pointer(new([32]uint64))
	case 33:
		up = unsafe.Pointer(new([33]uint64))
	default:
		slice := make([]uint64, 0, numWords)
		header := (*reflect.SliceHeader)(unsafe.Pointer(&slice))
		up = unsafe.Pointer(header.Data)
	}

	// we now have an unsafe.Pointer to allocated memory that we can use to
	// store the string length and content. We need to view it as a byte slice
	// so we can initialize it
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Cap = numWords * 8
	header.Len = totalLen
	header.Data = uintptr(up)

	binary.PutUvarint(bytes, uint64(length))
	copy(bytes[prefixLen:], []byte(s))

	return String{p: up}
}

func (s String) Len() int {
	var bytes []byte
	up := atomic.LoadPointer(&s.p)
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Cap = 8
	header.Len = 8
	header.Data = uintptr(up)
	len, _ := binary.Uvarint(bytes)
	return int(len)
}

func (s String) String() string {
	var bytes []byte
	up := atomic.LoadPointer(&s.p)
	header := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	header.Cap = 8
	header.Len = 8
	header.Data = uintptr(up)
	len, offset := binary.Uvarint(bytes)

	var rv string
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&rv))
	stringHeader.Len = int(len)
	stringHeader.Data = uintptr(up) + uintptr(offset)
	return rv
}
