VERSION ?= dev

.PHONY: binary
binary: clean version build

.PHONY: clean
clean:
	rm -rf *.deb
	rm -rf *.rpm
	rm -rf bla
	rm -rf .tmpBuildRoot

.PHONY: version
version:
	sed -i .bak -e 's/dev/${VERSION}/g' version.go

.PHONY: build
build:
	mkdir -p src/github.com/mengzhuo
	@-ln -s ${PWD} ${PWD}/src/github.com/mengzhuo/bla && ([ $$? -eq 0 ] )
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
