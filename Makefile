build:
	go build -buildvcs=false -o cmd/shortener/shortener cmd/shortener/main.go

tests:
	go test -v ./...

linter:
	go vet -vettool=statictest ./...

iter1:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -test.run=^TestIteration1$

iter2:
	./shortenertestbeta -test.v -source-path=./ -test.run=^TestIteration2$

localtest: build linter tests

autotests: iter1 iter2
