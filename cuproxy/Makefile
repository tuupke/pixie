OS:=$(shell go env GOOS)
ARCH:=$(shell go env GOARCH)

local: cuproxy-${OS}-${ARCH}

linux: cuproxy-linux-amd64 cuproxy-linux-386 cuproxy-linux-arm64 cuproxy-linux-arm
macos: cuproxy-macos-amd64 cuproxy-macos-arm64

run: local
	./cuproxy-${OS}-${ARCH}

releases: release
release: vendor linux macos docker

docker: vendor
	docker build -t tuupke/cuproxy .

vendor:
	go mod vendor

cuproxy-linux-%:
	GOOS=linux GOARCH=$* CGO_ENABLED=0 go build -mod=vendor -o $@ .

cuproxy-macos-%:
	GOOS=darwin GOARCH=$* CGO_ENABLED=0 go build -mod=vendor -o $@ .

clean:
	rm -f cuproxy-{linux,macos}-{amd64,386,arm,arm64}
