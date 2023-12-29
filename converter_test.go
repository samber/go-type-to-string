package typetostring

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testFunc1()                             {}              //nolint:unused
func testFunc2(string, assert.TestingT) bool { return true } //nolint:unused
func testFunc3(...string)                    {}              //nolint:unused

func TestGetType(t *testing.T) {
	is := assert.New(t)

	type testStruct struct{}          //nolint:unused
	type testInterface interface{}    //nolint:unused
	type testGen[T any] struct{ t T } //nolint:unused

	// simple types
	name := GetType[int]()
	is.Equal("int", name)
	name = GetType[string]()
	is.Equal("string", name)
	name = GetType[complex128]()
	is.Equal("complex128", name)
	name = GetType[uint32]()
	is.Equal("uint32", name)

	// stdlib types
	name = GetType[error]()
	is.Equal("*error", name)

	// simple types with pointer and slices
	name = GetType[[]int]()
	is.Equal("[]int", name)
	name = GetType[*int]()
	is.Equal("*int", name)
	name = GetType[*[]int]()
	is.Equal("*[]int", name)
	name = GetType[[]*int]()
	is.Equal("[]*int", name)
	name = GetType[*[]*int]()
	is.Equal("*[]*int", name)
	name = GetType[*[]*[]**int]()
	is.Equal("*[]*[]**int", name)

	// structs and interfaces
	name = GetType[testStruct]()
	is.Equal("github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[testInterface]()
	is.Equal("*github.com/samber/go-type-to-string.testInterface", name)

	// structs and interfaces with pointer and slices
	name = GetType[[]testStruct]()
	is.Equal("[]github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[*testStruct]()
	is.Equal("*github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[*[]testStruct]()
	is.Equal("*[]github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[[]*testStruct]()
	is.Equal("[]*github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[*[]*testStruct]()
	is.Equal("*[]*github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[*[]*[]**testStruct]()
	is.Equal("*[]*[]**github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[***testStruct]()
	is.Equal("***github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[*testInterface]()
	is.Equal("*github.com/samber/go-type-to-string.testInterface", name)
	name = GetType[***testInterface]()
	is.Equal("***github.com/samber/go-type-to-string.testInterface", name)

	// generic types
	name = GetType[testGen[int]]()
	is.Equal("github.com/samber/go-type-to-string.testGen[int]", name)
	// @TODO: fix this
	// name = GetType[testGen[testGen[int]]]()
	// is.Equal("github.com/samber/go-type-to-string.testGen[github.com/samber/go-type-to-string.testGen[int]]", name)

	// generic types with pointer and slices
	name = GetType[[]testGen[int]]()
	is.Equal("[]github.com/samber/go-type-to-string.testGen[int]", name)
	name = GetType[*testGen[int]]()
	is.Equal("*github.com/samber/go-type-to-string.testGen[int]", name)
	name = GetType[*[]testGen[int]]()
	is.Equal("*[]github.com/samber/go-type-to-string.testGen[int]", name)
	name = GetType[[]*testGen[int]]()
	is.Equal("[]*github.com/samber/go-type-to-string.testGen[int]", name)
	name = GetType[*[]*testGen[int]]()
	is.Equal("*[]*github.com/samber/go-type-to-string.testGen[int]", name)
	name = GetType[*[]*[]**testGen[int]]()
	is.Equal("*[]*[]**github.com/samber/go-type-to-string.testGen[int]", name)

	// maps
	name = GetType[map[string]int]()
	is.Equal("map[string]int", name)
	name = GetType[map[*string]int]()
	is.Equal("map[*string]int", name)
	name = GetType[*map[string]int]()
	is.Equal("*map[string]int", name)
	name = GetType[*[]*map[*testStruct]testInterface]()
	is.Equal("*[]*map[*github.com/samber/go-type-to-string.testStruct]github.com/samber/go-type-to-string.testInterface", name)
	name = GetType[*[]*map[*testStruct][]map[int]*testInterface]()
	is.Equal("*[]*map[*github.com/samber/go-type-to-string.testStruct][]map[int]*github.com/samber/go-type-to-string.testInterface", name)

	// channels
	name = GetType[chan int]()
	is.Equal("chan int", name)
	name = GetType[<-chan int]()
	is.Equal("<-chan int", name)
	name = GetType[chan<- int]()
	is.Equal("chan<- int", name)
	name = GetType[chan testStruct]()
	is.Equal("chan github.com/samber/go-type-to-string.testStruct", name)
	name = GetType[chan testInterface]()
	is.Equal("chan github.com/samber/go-type-to-string.testInterface", name)
	name = GetType[chan *[]*map[*testStruct][]map[chan int]*testInterface]()
	is.Equal("chan *[]*map[*github.com/samber/go-type-to-string.testStruct][]map[chan int]*github.com/samber/go-type-to-string.testInterface", name)

	// functions
	name = GetValueType(testFunc1)
	is.Equal("func()", name)
	name = GetValueType(testFunc2)
	is.Equal("func(string, assert.TestingT) bool", name)
	name = GetValueType(testFunc3)
	is.Equal("func(...string)", name)
	// @TODO: fix this
	// name = GetValueType(func() *testStruct { return nil })
	// is.Equal("func() *github.com/samber/go-type-to-string.testStruct", name)

	// anonymous types
	name = GetType[func()]()
	is.Equal("func()", name)
	name = GetType[struct{ foo int }]()
	is.Equal("struct { foo int }", name)
	// @TODO: fix this
	// name = GetType[struct{ foo testStruct }]()
	// is.Equal("struct { foo github.com/samber/go-type-to-string.testStruct }", name)

	// any
	name = GetType[any]()
	is.Equal("interface {}", name)
	name = GetType[interface{}]()
	is.Equal("interface {}", name)
	name = GetType[*any]()
	is.Equal("*interface {}", name)

	// all mixed
	name = GetType[[]chan *[]*map[*testStruct][]map[chan int]*map[testInterface]func(int, string) bool]()
	is.Equal("[]chan *[]*map[*github.com/samber/go-type-to-string.testStruct][]map[chan int]*map[github.com/samber/go-type-to-string.testInterface]func(int, string) bool", name)
}

func TestGetReflectValueType(t *testing.T) {
	is := assert.New(t)

	type testStruct struct{}          //nolint:unused
	type testInterface interface{}    //nolint:unused
	type testGen[T any] struct{ t T } //nolint:unused

	// random tests
	name := GetReflectValueType(reflect.ValueOf(42))
	is.Equal("int", name)
	name = GetType[[]int]()
	is.Equal("[]int", name)
	name = GetReflectValueType(reflect.ValueOf(testStruct{}))
	is.Equal("github.com/samber/go-type-to-string.testStruct", name)
	name = GetReflectValueType(reflect.ValueOf(testFunc2))
	is.Equal("func(string, assert.TestingT) bool", name)

	// @TODO: missing tests
}
