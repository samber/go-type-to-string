package typetostring

import (
	"reflect"
	"runtime"
	"strings"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func check[T any](equalReflectString bool, t *testing.T, expected string) {
	t.Helper()

	assert.Equal(t, expected, GetType[T](), "GetType")

	var v T
	assert.Equal(t, expected, GetValueType(v), "GetValueType")

	reflectType := reflect.TypeOf(&v).Elem()
	assert.Equal(t, expected, GetReflectType(reflectType), "GetReflectType")

	reflectValue := reflect.ValueOf(&v).Elem()
	assert.Equal(t, expected, GetReflectValueType(reflectValue), "GetReflectValueType")

	checkReflectString := assert.NotEqual

	if equalReflectString {
		checkReflectString = assert.Equal
	}

	checkReflectString(t, expected, reflectType.String(), "equalReflectString")
}

func Test(t *testing.T) {
	type testStruct struct{}          //nolint:unused
	type testInterface interface{}    //nolint:unused
	type testGen[T any] struct{ t T } //nolint:unused

	// simple types
	check[int](true,
		t, "int")
	check[string](true,
		t, "string")
	check[complex128](true,
		t, "complex128")
	check[uint32](true,
		t, "uint32")
	check[rune](true,
		t, "int32")

	// stdlib types
	check[error](true,
		t, "error")
	check[*error](true,
		t, "*error")

	// simple types with pointer and slices
	check[[]int](true,
		t, "[]int")
	check[*int](true,
		t, "*int")
	check[*[]int](true,
		t, "*[]int")
	check[[]*int](true,
		t, "[]*int")
	check[*[]*int](true,
		t, "*[]*int")
	check[*[]*[]**int](true,
		t, "*[]*[]**int")

	// structs and interfaces
	check[testStruct](false,
		t, "github.com/samber/go-type-to-string.testStruct")
	check[testInterface](false,
		t, "github.com/samber/go-type-to-string.testInterface")

	// structs and interfaces with pointer and slices
	check[[]testStruct](false,
		t, "[]github.com/samber/go-type-to-string.testStruct")
	check[*testStruct](false,
		t, "*github.com/samber/go-type-to-string.testStruct")
	check[*[]testStruct](false,
		t, "*[]github.com/samber/go-type-to-string.testStruct")
	check[[]*testStruct](false,
		t, "[]*github.com/samber/go-type-to-string.testStruct")
	check[*[]*testStruct](false,
		t, "*[]*github.com/samber/go-type-to-string.testStruct")
	check[*[]*[]**testStruct](false,
		t, "*[]*[]**github.com/samber/go-type-to-string.testStruct")
	check[***testStruct](false,
		t, "***github.com/samber/go-type-to-string.testStruct")
	check[*testInterface](false,
		t, "*github.com/samber/go-type-to-string.testInterface")
	check[***testInterface](false,
		t, "***github.com/samber/go-type-to-string.testInterface")

	// generic types
	check[testGen[int]](false, t, "github.com/samber/go-type-to-string.testGen[int]")
	check[testGen[testing.T]](false, t, "github.com/samber/go-type-to-string.testGen[testing.T]")
	check[testGen[testing.B]](false, t, "github.com/samber/go-type-to-string.testGen[testing.B]")
	check[testGen[assert.Assertions]](false, t, "github.com/samber/go-type-to-string.testGen[github.com/stretchr/testify/assert.Assertions]")
	check[testGen[func(assert.Assertions)]](false, t, "github.com/samber/go-type-to-string.testGen[func(github.com/stretchr/testify/assert.Assertions)]")
	check[testGen[func(testing.T, ...assert.Assertions)]](false, t, "github.com/samber/go-type-to-string.testGen[func(testing.T, ...github.com/stretchr/testify/assert.Assertions)]")
	// @TODO: fix this
	// check[testGen[testStruct]](false, t, "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testStruct]")

	{ // generic with nested local types
		type testInt int
		var expected struct{ testStruct, testInterface, testGenInt, testInt string }

		switch strings.Join(strings.Split(runtime.Version(), ".")[:2], ".") {
		case "go1.18":
			expected.testStruct = "github.com/samber/go-type-to-string.testGen[typetostring.testStruct·1]"
			expected.testInterface = "github.com/samber/go-type-to-string.testGen[typetostring.testInterface·2]"
			expected.testGenInt = "github.com/samber/go-type-to-string.testGen[typetostring.testGen[int]]"
			expected.testInt = "github.com/samber/go-type-to-string.testGen[typetostring.testInt·4]"
		case "go1.19":
			expected.testStruct = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testStruct·1]"       // as 1.20
			expected.testInterface = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testInterface·2]" // as 1.20
			expected.testGenInt = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testGen[int]]"       // no `·3]` for local generic type
			expected.testInt = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testInt·4]"             // as 1.20
		default: // go1.20 and later
			expected.testStruct = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testStruct·1]"
			expected.testInterface = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testInterface·2]"
			expected.testGenInt = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testGen[int]·3]"
			expected.testInt = "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testInt·4]"

			check[testGen[testGen[testInterface]]](false,
				t, "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testInterface·2]·3]")
		}

		check[testGen[testStruct]](false, t, expected.testStruct)
		check[testGen[testInterface]](false, t, expected.testInterface)
		check[testGen[testGen[int]]](false, t, expected.testGenInt)
		check[testGen[testInt]](false, t, expected.testInt)
	}

	// generic types with pointer and slices
	check[[]testGen[int]](false,
		t, "[]github.com/samber/go-type-to-string.testGen[int]")
	check[*testGen[int]](false,
		t, "*github.com/samber/go-type-to-string.testGen[int]")
	check[*[]testGen[int]](false,
		t, "*[]github.com/samber/go-type-to-string.testGen[int]")
	check[[]*testGen[int]](false,
		t, "[]*github.com/samber/go-type-to-string.testGen[int]")
	check[*[]*testGen[int]](false,
		t, "*[]*github.com/samber/go-type-to-string.testGen[int]")
	check[*[]*[]**testGen[int]](false,
		t, "*[]*[]**github.com/samber/go-type-to-string.testGen[int]")

	// maps
	check[map[string]int](true,
		t, "map[string]int")
	check[map[*string]int](true,
		t, "map[*string]int")
	check[*map[string]int](true,
		t, "*map[string]int")
	check[*[]*map[*testStruct]testInterface](false,
		t, "*[]*map[*github.com/samber/go-type-to-string.testStruct]github.com/samber/go-type-to-string.testInterface")
	check[*[]*map[*testStruct][]map[int]*testInterface](false,
		t, "*[]*map[*github.com/samber/go-type-to-string.testStruct][]map[int]*github.com/samber/go-type-to-string.testInterface")
	check[map[testStruct]int](false,
		t, "map[github.com/samber/go-type-to-string.testStruct]int")

	// arrays
	check[[1]int](true,
		t, "[1]int")
	check[[2]*int](true,
		t, "[2]*int")
	check[[3]*[4]testStruct](false,
		t, "[3]*[4]github.com/samber/go-type-to-string.testStruct")

	// channels
	check[chan int](true,
		t, "chan int")
	check[<-chan int](true,
		t, "<-chan int")
	check[chan<- int](true,
		t, "chan<- int")
	check[chan testStruct](false,
		t, "chan github.com/samber/go-type-to-string.testStruct")
	check[chan testInterface](false,
		t, "chan github.com/samber/go-type-to-string.testInterface")
	check[chan *[]*map[*testStruct][]map[chan int]*testInterface](false,
		t, "chan *[]*map[*github.com/samber/go-type-to-string.testStruct][]map[chan int]*github.com/samber/go-type-to-string.testInterface")

	// functions
	check[func()](true,
		t, "func()")
	check[func(string, assert.TestingT) bool](false,
		t, "func(string, github.com/stretchr/testify/assert.TestingT) bool")
	check[func(...string)](true,
		t, "func(...string)")
	check[func(int, ...**testStruct) (string, *int)](false,
		t, "func(int, ...**github.com/samber/go-type-to-string.testStruct) (string, *int)")
	check[func() *testStruct](false,
		t, "func() *github.com/samber/go-type-to-string.testStruct")
	check[func(func(assert.TestingT) *func(...string)) *func() *func()](false,
		t, "func(func(github.com/stretchr/testify/assert.TestingT) *func(...string)) *func() *func()")
	check[func() *[]*func(...string) *func() (int, *testStruct)](false,
		t, "func() *[]*func(...string) *func() (int, *github.com/samber/go-type-to-string.testStruct)")
	check[func() *[]*func(...string) *func() (int, *func() *[]*func(...string) *func())](true,
		t, "func() *[]*func(...string) *func() (int, *func() *[]*func(...string) *func())")

	// anonymous types
	check[func()](true,
		t, "func()")
	check[struct{ foo int }](true,
		t, "struct { foo int }")
	check[struct{ foo testStruct }](false,
		t, "struct { foo github.com/samber/go-type-to-string.testStruct }")
	check[struct{ testStruct }](false,
		t, "struct { github.com/samber/go-type-to-string.testStruct }")
	check[func(struct{ foo testStruct })](false,
		t, "func(struct { foo github.com/samber/go-type-to-string.testStruct })")
	check[func(struct{ *testStruct })](false,
		t, "func(struct { *github.com/samber/go-type-to-string.testStruct })")
	check[chan struct{ foo testStruct }](false,
		t, "chan struct { foo github.com/samber/go-type-to-string.testStruct }")
	check[chan struct{ testStruct }](false,
		t, "chan struct { github.com/samber/go-type-to-string.testStruct }")
	check[struct {
		foo int
		testStruct
	}](false,
		t, "struct { foo int; github.com/samber/go-type-to-string.testStruct }")
	check[func(...struct{ bar string })](true,
		t, "func(...struct { bar string })")
	check[struct{ foo struct{ bar int } }](true,
		t, "struct { foo struct { bar int } }")
	check[interface{ Do() string }](true,
		t, "interface { Do() string }")
	check[*interface{ Do() string }](true,
		t, "*interface { Do() string }")
	check[[]interface{ Do() string }](true,
		t, "[]interface { Do() string }")
	check[func() interface{ Do() string }](true,
		t, "func() interface { Do() string }")
	check[interface {
		A() int
		B() string
	}](true,
		t, "interface { A() int; B() string }")
	check[interface {
		A() int
		interface{ B() string }
	}](true,
		t, "interface { A() int; B() string }")
	check[struct {
		foo int
		bar struct{ baz string }
	}](true,
		t, "struct { foo int; bar struct { baz string } }")

	// unsafe
	check[unsafe.Pointer](true,
		t, "unsafe.Pointer")
	check[func(unsafe.Pointer)](true,
		t, "func(unsafe.Pointer)")
	check[struct{ p unsafe.Pointer }](true,
		t, "struct { p unsafe.Pointer }")
	check[interface{ P() unsafe.Pointer }](true,
		t, "interface { P() unsafe.Pointer }")

	// any
	check[any](true,
		t, "interface {}")
	check[interface{}](true,
		t, "interface {}")
	check[*any](true,
		t, "*interface {}")
	check[**any](true,
		t, "**interface {}")

	// named types
	type ptr *any
	check[ptr](false, t, "github.com/samber/go-type-to-string.ptr")
	type slice []any
	check[slice](false, t, "github.com/samber/go-type-to-string.slice")
	type array [0]any
	check[array](false, t, "github.com/samber/go-type-to-string.array")
	type set map[any]struct{}
	check[set](false, t, "github.com/samber/go-type-to-string.set")
	type channel chan any
	check[channel](false, t, "github.com/samber/go-type-to-string.channel")
	type function func()
	check[function](false, t, "github.com/samber/go-type-to-string.function")
	type empty struct{}
	check[empty](false, t, "github.com/samber/go-type-to-string.empty")
	type aught interface{}
	check[aught](false, t, "github.com/samber/go-type-to-string.aught")

	check[*ptr](false, t, "*github.com/samber/go-type-to-string.ptr")
	check[[]ptr](false, t, "[]github.com/samber/go-type-to-string.ptr")
	check[chan<- ptr](false, t, "chan<- github.com/samber/go-type-to-string.ptr")

	// recursive types
	type recursive struct {
		r *recursive
	}
	check[recursive](false, t, "github.com/samber/go-type-to-string.recursive")

	// all mixed
	check[[]chan *[]*map[*testStruct][]map[chan int]*map[testInterface]func(int, string) bool](false,
		t, "[]chan *[]*map[*github.com/samber/go-type-to-string.testStruct][]map[chan int]*map[github.com/samber/go-type-to-string.testInterface]func(int, string) bool")
	check[[]chan *[]*map[*func()][]map[chan int]*map[struct{ int }]func(int, string) (bool, <-chan struct{})](true,
		t, "[]chan *[]*map[*func()][]map[chan int]*map[struct { int }]func(int, string) (bool, <-chan struct {})")
	check[[]chan *[10]*map[*func()][]map[chan int]*map[*func() <-chan func()]func(int, string) (bool, <-chan func(chan<- int))](true,
		t, "[]chan *[10]*map[*func()][]map[chan int]*map[*func() <-chan func()]func(int, string) (bool, <-chan func(chan<- int))")
}

func TestGetValueType(t *testing.T) {
	is := assert.New(t)

	var a any
	is.Equal("interface {}", GetValueType(a))
	is.Equal("*interface {}", GetValueType(&a))

	b := 123
	is.Equal("int", GetValueType(b))

	c := 1.23
	is.Equal("float64", GetValueType(c))

	d := true
	is.Equal("bool", GetValueType(d))

	var i interface{ f() }
	is.Equal("interface { typetostring.f() }", GetValueType(i))

	// @TODO: show "interface {}" or the underlying type ?
	for _, v := range []any{
		i, &a, 0, "", []any{}, [1]any{}, make(chan any), struct{}{},
	} {
		is.Equal("interface {}", GetValueType(v), reflect.TypeOf(v))
	}
}
