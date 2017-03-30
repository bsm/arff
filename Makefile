default: test

deps:
	go get -t ./...

test:
	go test ./... -v 1

bench:
	go test ./... -bench=. -v 1

README.md: README.md.tpl $(wildcard *.go)
	becca -package $(subst $(GOPATH)/src/,,$(PWD))
