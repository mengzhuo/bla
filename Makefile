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
	mkdir -p .tmpBuildRoot/etc/bla/webroot/
	mkdir -p .tmpBuildRoot/etc/systemd/system/
	mkdir -p .tmpBuildRoot/etc/logrotate.d/
	mkdir -p .tmpBuildRoot/usr/local/bin/
	cp bla .tmpBuildRoot/usr/local/bin/

deb: clean build pkg
	fpm -t deb -s dir -n bla --config-files /etc .tmpBuildRoot

rpm: clean build pkg
	fpm -t rpm -s dir -n bla --config-files /etc .tmpBuildRoot
