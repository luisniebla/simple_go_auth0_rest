package user

import (
	"fmt"
	"math/cmplx"
	"unsafe"
)

var c, python, java bool // declare list of variables
// When you need an integer value you should use int unless you have a specific reason to use a sized or unsigned integer type.
var a, b int = 0, 0

// outside function, every statement must begin with keyword, can't use :=
// c := 0

func ExportFunc() {
	a := 1
	fmt.Print("variable shadow", a)
	fmt.Print("Hello user")
}

func notExported() {
	fmt.Println("yo")
}

var (
	BoolVar bool = true
	IntVar  int  = 12412521541254
	Int8Var int8 = 122
)

var (
	ToBe       bool       = false
	MaxInt     uint64     = 1<<64 - 1
	z          complex128 = cmplx.Sqrt(-5 + 12i)
	Int16Var   int16      = 12412
	Intuintptr uintptr    = 12311515125125
	RuneVar    rune       = 'a'
)

func testUintptr() {
	var x int = 42

	p := &x

	var addr2 uintptr

	addr2 = uintptr(unsafe.Pointer(p))

	fmt.Printf("memmeory addr4ess %d", addr2)
	fmt.Printf("Type: %T Value: %v\n", ToBe, ToBe)
	fmt.Printf("Type: %T Value: %v\n", MaxInt, MaxInt)
	fmt.Printf("Type: %T Value: %v\n", z, z)
	fmt.Printf("Type: %T, Value: %v\n", Int8Var, Int8Var)
	fmt.Printf("Type: %T, Value: %v\n", Int16Var, Int16Var)
	fmt.Printf("Type: %T, Value: %v\n", Intuintptr, Intuintptr)

	var x2 int = 42

	p2 := &x2

	var addr uintptr

	addr = uintptr(unsafe.Pointer(p2))

	fmt.Printf("memmeory addr4ess %d\n", addr)

	// rune is alias for int32
	fmt.Printf("Type: %T, Value: %v\n", RuneVar, RuneVar)
}

func testZeroValues() {
	var i int
	var f float64
	var b bool
	var s string
	var i8 int8
	var i16 int16
	var i32 int32
	var runi rune
	fmt.Printf("%v %v %v %q\n", i, f, b, s)
	fmt.Println(i8, i16, i32, runi)
}

func testFor() {
	sum := 1
	for i := 1; i < 10; i += 1 {
		fmt.Println(i)
	}
	for sum < 10 {
		fmt.Println(sum)
	}
	for sum < 10 { // while loop
		fmt.Println(sum)
	}
	for { // ifninite loop

	}
}

func testIf() bool {
	a := 1
	if b := 2; a == b {
		return true
	} else if b > a { // scope continues to else
		return false
	}
	// b no longer exists
	return false
}

// pointers
func testPointers() {
	i := 42
	p := &i         // p is a pointer to the memory address of the value of i
	fmt.Println(*p) // read i through the pointer p
	*p = 21         // sets the value of i using the point
}

// structs are collections of fields
type V struct {
	X int
	Y int
}

func testVertex() {
	fmt.Println(V{1, 2})
	fmt.Println(V{X: 1, Y: 2})
}

// maps map keys to values
func testMap() {
	// var m map[string]interface{}
	// m["1"] = 2 returns error because you need ot make, zero value of map is nil
	m := make(map[string]interface{})
	m["1"] = 1
	delete(m, "1")
	elem, ok := m["1"]
	if ok == true {
		fmt.Println(elem)
	} else {
		fmt.Println("1 does not exist")
	}
}

// basically just a constant
func testMapLiteratal() {
	var m = map[string]string{
		"1": "2",
	}
	fmt.Println(m)
}
