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

iter4:
	./shortenertestbeta -test.v -server-port=8080 -binary-path=cmd/shortener/shortener -test.run=^TestIteration4$

iter5:
	./shortenertestbeta -test.v -server-port=8080 -binary-path=cmd/shortener/shortener -test.run=^TestIteration5$
	
iter6:
	./shortenertestbeta -test.v -source-path=./ -test.run=^TestIteration6$
	
iter7:
	./shortenertestbeta -test.v -source-path=./ -binary-path=cmd/shortener/shortener -test.run=^TestIteration7$
	
iter8:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -test.run=^TestIteration8$

iter9:
	./shortenertestbeta -test.v -source-path=. -binary-path=cmd/shortener/shortener -file-storage-path="storage.json" -test.run=^TestIteration9$

iter10:
	./shortenertestbeta -test.v -source-path=. -binary-path=cmd/shortener/shortener -database-dsn='postgres://postgres:@localhost:5432/url_shortener?sslmode=disable&search_path=public' -test.run=^TestIteration10$

iter11:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -database-dsn='postgres://postgres:@localhost:5432/url_shortener?sslmode=disable&search_path=public' -test.run=^TestIteration11$

iter12:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -database-dsn='postgres://postgres:@localhost:5432/url_shortener?sslmode=disable&search_path=public' -test.run=^TestIteration12$

iter13:
	./shortenertestbeta -test.v -binary-path=cmd/shortener/shortener -database-dsn='postgres://postgres:@localhost:5432/url_shortener?sslmode=disable&search_path=public' -test.run=^TestIteration13$

localtests: build linter tests

autotests: iter1 iter2 iter3 iter4 iter5 iter6 iter7 iter8 iter9 iter10 iter12 iter13
