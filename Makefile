build:
	go build -i

install:
	go install -v

fmt:
	go fmt ./...

test:
	go test ./... -v

release: VERSION := $(shell cat ./VERSION)
release:
	git tag -a $(VERSION) -m "Release" || true
	git push origin $(VERSION)
	goreleaser --rm-dist

.PHONY: build
