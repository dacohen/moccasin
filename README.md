# moccasin

[![Build Status][ci-badge]][ci-runs] [![Go Reference][reference-badge]][reference]

The comfortable, un-opinionated mocking tool.
Write better, more flexible mocks in go, making writing tests less of a chore.

## Why moccasin?
We've got you covered: [Why Moccasin?](docs/why_moccasin.md)

## Installation
Install as you would any other go package:
```
go get github.com/rainforestpay/moccasin
```

## Using
Usage is straightforward, simply embed moccasin in the struct you want to mock:

```go
type MyStruct struct {
	moccasin.Embed
	
	Name string
	Value int64
	// ...
}
```

See the examples in the documentation for more information.

[ci-badge]:            https://github.com/rainforestpay/moccasin/actions/workflows/test.yaml/badge.svg
[ci-runs]:             https://github.com/rainforestpay/moccasin/actions
[reference-badge]:     https://pkg.go.dev/badge/github.com/rainforestpay/moccasin.svg
[reference]:           https://pkg.go.dev/github.com/rainforestpay/moccasin