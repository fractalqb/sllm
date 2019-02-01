GOSRC:=$(shell find . -name '*.go')

README.html: README.md
	pandoc -f gfm -t html -s -M title="sllm – README" README.md > README.html

# → https://blog.golang.org/cover
cover: coverage.html

benchmark:
	go test -bench=.

cpuprof:
	go test -cpuprofile cpu.prof -bench BenchmarkExpandArgs
# Read with '$ go tool pprof cpu.prof' >>> e.g. '(pprof) web'

coverage.html: coverage.out
	go tool cover -html=$< -o $@

coverage.out: $(GOSRC)
	go test -coverprofile=$@ ./... || true
#	go test -covermode=count -coverprofile=$@ || true

cov=$(shell go tool cover -func=coverage.out \
            | egrep '^total:' \
            | awk '{print $$3}' \
            | tr "%" " ")
