package typetostring

import (
	"reflect"
	"testing"

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
	check[testGen[int]](false,
		t, "github.com/samber/go-type-to-string.testGen[int]")
	// @TODO: fix this
	// check[testGen[testGen[int]]](false,
	// 	t, "github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testGen[int]]")

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
	// @TODO: fix this
	// check[struct{ foo testStruct }](false,
	// 	t, "struct { foo github.com/samber/go-type-to-string.testStruct }")

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

	a = ""
	is.Equal("interface {}", GetValueType(a)) // not string ?

	is.Equal("*interface {}", GetValueType(&a))

	var i interface{ f() }
	a = i
	is.Equal("interface {}", GetValueType(a)) // not interface{ f() } ?

	for _, v := range []any{
		i, &a, 0, "", []any{}, [1]any{}, make(chan any), struct{}{},
	} {
		is.Equal("interface {}", GetValueType(v), reflect.TypeOf(v))
	}
}
