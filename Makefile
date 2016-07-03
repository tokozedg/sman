VERSION := $(shell awk -F= '/version =/ {print $$2}' lib/root.go | tr -d "\" ")

all:
	go build -v

release:
	@$(MAKE) linux_amd64
	@$(MAKE) linux_386
	@$(MAKE) linux_arm
	@$(MAKE) darwin_amd64

linux_amd64: GOOS=linux
linux_amd64: GOARCH=amd64
linux_amd64: build

linux_386: GOOS=linux
linux_386: GOARCH=386
linux_386: build

linux_arm: GOOS=linux
linux_arm: GOARCH=arm
linux_arm: build

darwin_amd64: GOOS=darwin
darwin_amd64: GOARCH=amd64
darwin_amd64: build

build:
	env GOOS=${GOOS} GOARCH=${GOARCH} go build -o bin/sman-${GOOS}-${GOARCH}-${VERSION}
	cd bin; tar -czf sman-${GOOS}-${GOARCH}-${VERSION}.tgz sman-${GOOS}-${GOARCH}-${VERSION}
	rm bin/sman-${GOOS}-${GOARCH}-${VERSION}

test:
		go -v test ./...

watch:
	CompileDaemon -command="go test -v ./..." -color=True --log-prefix=False --exclude-dir=.git

.PHONY: all release linux_amd64 linux_386 linux_arm \
	darwin_amd64 build test watch
