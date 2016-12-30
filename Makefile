VERSION := $(shell git describe --tags)
DESTDIR?=.tmpBuildRoot

.PHONY: binary
binary: clean build

.PHONY: clean
clean:
	rm -rf *.deb
	rm -rf *.rpm
	rm -rf bla
	rm -rf ${DESTDIR} 

.PHONY: build
build:
	mkdir -p src/github.com/mengzhuo
	@-ln -s ${PWD} ${PWD}/src/github.com/mengzhuo/bla && ([ $$? -eq 0 ] )
	go build -o bla -ldflags '-X main.Version=${VERSION}' cmd/bla/main.go

.PHONY: pkg
pkg:
	rm -rf ${DESTDIR}
	mkdir ${DESTDIR}
	cp -rf buildRoot/* ${DESTDIR}/
	mkdir -p ${DESTDIR}/usr/local/bin
	mkdir -p ${DESTDIR}/var/log/bla/
	cp bla ${DESTDIR}/usr/local/bin/

deb: clean build pkg
	fpm -t deb -s dir -n bla -v $(VERSION:v%=%) -C ${DESTDIR}

rpm: clean build pkg
	fpm -t rpm -s dir -n bla -v ${VERSION} -C ${DESTDIR} 
