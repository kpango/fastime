# fastime [![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT) [![release](https://img.shields.io/github/release/kpango/fastime.svg)](https://github.com/kpango/fastime/releases/latest) [![CircleCI](https://circleci.com/gh/kpango/fastime.svg?style=shield)](https://circleci.com/gh/kpango/fastime) [![codecov](https://codecov.io/gh/kpango/fastime/branch/master/graph/badge.svg)](https://codecov.io/gh/kpango/fastime) [![Codacy Badge](https://api.codacy.com/project/badge/Grade/b9fa9b846ec343d3860b8f69e802c09b)](https://www.codacy.com/app/i.can.feel.gravity/fastime?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=kpango/fastime&amp;utm_campaign=Badge_Grade) [![Go Report Card](https://goreportcard.com/badge/github.com/kpango/fastime)](https://goreportcard.com/report/github.com/kpango/fastime) [![GoDoc](http://godoc.org/github.com/kpango/fastime?status.svg)](http://godoc.org/github.com/kpango/fastime) [![Join the chat at https://gitter.im/kpango/fastime](https://badges.gitter.im/kpango/fastime.svg)](https://gitter.im/kpango/fastime?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
Blazing Fast Time Library for golang with zero allocs

## Requirement
Go 1.8

## Installation
```shell
go get github.com/kpango/fastime
```

## Example
```go
    now := fastime.Now()
    defer fastime.Stop()

    // Create Instance
    ft := fastime.New()
    defer ft.Stop()
    ft.Now()
```

## Benchmark

```
go test -count=5 -run=NONE -bench . -benchmem
goos: darwin
goarch: amd64
pkg: github.com/kpango/fastime
BenchmarkFastime-8   	2000000000	         0.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastime-8   	2000000000	         0.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastime-8   	2000000000	         0.48 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastime-8   	2000000000	         0.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkFastime-8   	2000000000	         0.47 ns/op	       0 B/op	       0 allocs/op
BenchmarkTime-8      	100000000	        16.6 ns/op	       1 B/op	       0 allocs/op
BenchmarkTime-8      	100000000	        15.4 ns/op	       1 B/op	       0 allocs/op
BenchmarkTime-8      	100000000	        15.6 ns/op	       1 B/op	       0 allocs/op
BenchmarkTime-8      	100000000	        15.6 ns/op	       1 B/op	       0 allocs/op
BenchmarkTime-8      	100000000	        14.9 ns/op	       1 B/op	       0 allocs/op
PASS
ok  	github.com/kpango/fastime	12.936s
```
## Contribution
1. Fork it ( https://github.com/kpango/fastime/fork )
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

## Author
[kpango](https://github.com/kpango)

## LICENSE
fastime released under MIT license, refer [LICENSE](https://github.com/kpango/fastime/blob/master/LICENSE) file.
