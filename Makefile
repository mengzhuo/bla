VERSION?=dev

.PHONY: clean
clean:
	rm -rf *.deb
	rm -rf *.rpm
	rm -rf bla
	rm -rf .tmpBuildRoot

.PHONY: build
build: 
	go build -o bla cmd/bla/main.go

.PHONY: pkg
pkg:
	rm -rf .tmpBuildRoot
	mkdir .tmpBuildRoot
	cp -rf buildRoot/* .tmpBuildRoot/
	cp bla .tmpBuildRoot/usr/local/bin/

deb: clean build pkg
	fpm -t deb -s dir -n bla  .tmpBuildRoot

rpm: clean build pkg
	fpm -t rpm -s dir -n bla  .tmpBuildRoot
