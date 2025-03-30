build:
	go build -buildvcs=false -o cmd/shortener/shortener.exe cmd/shortener/main.go

tests:
	go test -v ./...

linter:
	go vet -vettool=statictest ./...

iter1:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -test.run=^TestIteration1$

iter2:
	./shortenertestbeta -test.v -source-path=./ -test.run=^TestIteration2$
	
iter3:
	./shortenertestbeta -test.v -source-path=./ -test.run=^TestIteration3$

localtests: build linter tests

autotests: iter1 iter2 iter3
