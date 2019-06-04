GO_VERSION:=$(shell go version)

.PHONY: all clean bench bench-all profile lint test contributors update install

all: clean install lint test bench

clean:
	go clean ./...
	rm -rf ./*.log
	rm -rf ./*.svg
	rm -rf ./go.mod
	rm -rf ./go.sum
	rm -rf bench
	rm -rf pprof
	rm -rf vendor


bench: clean
	go test -count=10 -run=NONE -bench . -benchmem

profile: clean
	rm -rf bench
	mkdir bench
	mkdir pprof
	\
	go test -count=10 -run=NONE -bench=BenchmarkFastime -benchmem -o pprof/fastime-test.bin -cpuprofile pprof/cpu-fastime.out -memprofile pprof/mem-fastime.out
	go tool pprof --svg pprof/fastime-test.bin pprof/cpu-fastime.out > cpu-fastime.svg
	go tool pprof --svg pprof/fastime-test.bin pprof/mem-fastime.out > mem-fastime.svg
	go-torch -f bench/cpu-fastime-graph.svg pprof/fastime-test.bin pprof/cpu-fastime.out
	go-torch --alloc_objects -f bench/mem-fastime-graph.svg pprof/fastime-test.bin pprof/mem-fastime.out
	\
	go test -count=10 -run=NONE -bench=BenchmarkTime -benchmem -o pprof/default-test.bin -cpuprofile pprof/cpu-default.out -memprofile pprof/mem-default.out
	go tool pprof --svg pprof/default-test.bin pprof/mem-default.out > mem-default.svg
	go tool pprof --svg pprof/default-test.bin pprof/cpu-default.out > cpu-default.svg
	go-torch -f bench/cpu-default-graph.svg pprof/default-test.bin pprof/cpu-default.out
	go-torch --alloc_objects -f bench/mem-default-graph.svg pprof/default-test.bin pprof/mem-default.out
	\
	mv ./*.svg bench/

cpu:
	go tool pprof pprof/fastime-test.bin pprof/cpu-fastime.out

mem:
	go tool pprof --alloc_space pprof/fastime-test.bin pprof/mem-fastime.out

lint:
	gometalinter --enable-all . | rg -v comment

test: clean
	GO111MODULE=on go test --race -v $(go list ./... | rg -v vendor)

contributors:
	git log --format='%aN <%aE>' | sort -fu > CONTRIBUTORS
