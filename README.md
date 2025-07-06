
# Extract a string representation of Go type

[![tag](https://img.shields.io/github/tag/samber/go-type-to-string.svg)](https://github.com/samber/go-type-to-string/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18.0-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/go-type-to-string?status.svg)](https://pkg.go.dev/github.com/samber/go-type-to-string)
![Build Status](https://github.com/samber/go-type-to-string/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/go-type-to-string)](https://goreportcard.com/report/github.com/samber/go-type-to-string)
[![Coverage](https://img.shields.io/codecov/c/github/samber/go-type-to-string)](https://codecov.io/gh/samber/go-type-to-string)
[![Contributors](https://img.shields.io/github/contributors/samber/go-type-to-string)](https://github.com/samber/go-type-to-string/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/go-type-to-string)](./LICENSE)

## Motivations

For the [samber/do](https://github.com/samber/do) project, I needed to convert a Go type into a string. I used to convert it with `fmt.Sprintf("%T", t)` -> `mypkg.MyStruct`, but it does not insert package path into type representation, leading to collision when types from different pacakges match.

This package export type using the following representation:

```go
  *[]**<-chan   github.com/samber/example    .Example
  |             |                            ^
  |             |                            Type name
  |             ^
  |             The package path (including package name)
  ^
  Type indicators (map, slice, pointer, channel...)
```

This library supports:
- primitive types
- pointers
- structs
- functions with input and output
- vaargs
- interfaces
- maps
- arrays
- slices
- channels
- generics
- anonymous types
- named types
- `unsafe.Pointer`
- recursive types

Known limitations:
- structs in generic type
- `any("foobar")` is currently reported as `any` instead of `string` (see [#2](https://github.com/samber/go-type-to-string/issues/2))

## Examples

Using the following types:

```go
package example

type testStruct struct{}
type testGeneric[T any] struct{ t T }
type testNamedType testStruct
```

| Type                                         | Exported                                                                                         |
| -------------------------------------------- | ------------------------------------------------------------------------------------------------ |
| `int`                                        | `int`                                                                                            |
| `*int`                                       | `*int`                                                                                           |
| `**[]*int`                                   | `**[]*int`                                                                                       |
| `**[]*map[int]bool`                          | `**[]*map[int]bool`                                                                              |
| `func (a string, b bool) int`                | `func (string, bool) int`                                                                        |
| `func(int, ...string) (bool, error)`         | `func(int, ...string) (bool, error)`                                                             |
| `testStruct`                                 | `github.com/samber/example.testStruct`                                                           |
| `*testStruct`                                | `*github.com/samber/example.testStruct`                                                          |
| `***testStruct`                              | `***github.com/samber/example.testStruct`                                                        |
| `***testNamedType`                           | `***github.com/samber/example.testNamedType`                                                     |
| `[][3]***testStruct`                         | `[][3]***github.com/samber/example.testStruct`                                                   |
| `testGeneric[string]`                        | `github.com/samber/example.testGeneric[string]`                                                  |
| `*map[testStruct]chan<- testGeneric[string]` | `*map[github.com/samber/example.testStruct]chan<- github.com/samber/example.testGeneric[string]` |

See more examples [here](https://github.com/samber/go-type-to-string/blob/main/converter_test#L13)

## 🚀 Install

```sh
go get github.com/samber/go-type-to-string
```

This library is v1 and follows SemVer strictly. No breaking changes will be made to exported APIs before v2.0.0.

## 💡 How to

GoDoc: [https://pkg.go.dev/github.com/samber/go-type-to-string](https://pkg.go.dev/github.com/samber/go-type-to-string)

```go
package example

import converter "github.com/samber/go-type-to-string"

type Example struct{
    foo string
    bar int
}

func main() {
    name1 := converter.GetType[*Example]()
    // "*github.com/samber/example.Example"

    name2 := converter.GetValueType(Example{})
    // "github.com/samber/example.Example"
}
```

## 🤝 Contributing

- Ping me on Twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/go-type-to-string)
- Fix [open issues](https://github.com/samber/go-type-to-string/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/go-type-to-string)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
