OS   = linux
ARCH = arm

.PHONY: build

build:
	GOOS=$(OS) GOARCH=$(ARCH) go build -o build/app

.PHONY: clean
clean:
	go clean -i && rm -rf build
