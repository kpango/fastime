GO_VERSION:=$(shell go version)

.PHONY: all clean bench bench-all profile lint test contributors update install

GOPATH := $(eval GOPATH := $(shell go env GOPATH))$(GOPATH)
GOLINES_MAX_WIDTH     ?= 200

all: clean install lint test bench

clean:
	go clean ./...
	rm -rf ./*.log
	rm -rf ./*.svg
	rm -rf ./go.mod
	rm -rf bench
	rm -rf pprof
	rm -rf vendor
	go mod init
	go mod tidy


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
	\
	go test -count=10 -run=NONE -bench=BenchmarkTime -benchmem -o pprof/default-test.bin -cpuprofile pprof/cpu-default.out -memprofile pprof/mem-default.out
	go tool pprof --svg pprof/default-test.bin pprof/mem-default.out > mem-default.svg
	go tool pprof --svg pprof/default-test.bin pprof/cpu-default.out > cpu-default.svg
	\
	mv ./*.svg bench/

profile-web-cpu:
	go tool pprof -http=":6061" \
		pprof/fastime-test.bin \
		pprof/cpu-fastime.out

profile-web-mem:
	go tool pprof -http=":6062" \
		pprof/fastime-test.bin \
		pprof/mem-fastime.out

profile-web-cpu-default:
	go tool pprof -http=":6063" \
		pprof/default-test.bin \
		pprof/cpu-default.out

profile-web-mem-default:
	go tool pprof -http=":6064" \
		pprof/default-test.bin \
		pprof/mem-default.out



lint:
	gometalinter --enable-all . | rg -v comment

test: clean
	GO111MODULE=on go test --race -v $(go list ./... | rg -v vendor)

contributors:
	git log --format='%aN <%aE>' | sort -fu > CONTRIBUTORS

run:
	go run example/main.go

format:
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/golines -w -m $(GOLINES_MAX_WIDTH)
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/gofumpt -w
	find ./ -type d -name .git -prune -o -type f -regex '.*[^\.pb]\.go' -print | xargs $(GOPATH)/bin/strictgoimports -w
	find ./ -type d -name .git -prune -o -type f -regex '.*\.go' -print | xargs $(GOPATH)/bin/goimports -w
